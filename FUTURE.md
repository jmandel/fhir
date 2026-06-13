# A peek into the txpack future

This branch is a working, self-contained demo of the proposed terminology + performance future
for the core FHIR spec build. Everything it claims is runnable here, now:

```bash
git clone -b txpack-future https://github.com/jmandel/fhir.git && cd fhir
./eng/future/build.sh --judge        # or, on any OS, fetch the jar named in tx.lock and:
# java -Xmx12g -cp kindling-future-v5.jar org.hl7.fhir.tools.publisher.SpecBuild build . --judge
```

Everything is Java: the publisher reads `tx.lock` natively (content-addressed, sha256-verified
pack fetch), the `SpecBuild` CLI pins locale/timezone in-process, runs hermetic by default,
checks the output signature against the lock, and carries the verification subcommands
(`manifest` / `compare` / `impact` / `diff-packs`). The bash file is a 20-line bootstrap that
fetches the jar and execs it; Windows users skip it. **Dependency footprint: a JDK.**

That is a **fully cold build that makes zero terminology network requests** — provably (any
attempted request is a hard failure) — and ends by checking the published output byte-for-byte
against the committed reference manifest. On a 12-core machine it finishes in ~3.5 minutes; the
same build on the stock toolchain takes ~18 minutes cold against tx.fhir.org (±40% by hour),
and fails outright when the server is having a bad day.

## What's different on this branch

Three files and one directory — nothing else in the spec changed:

| What | Why |
|---|---|
| [`tx.lock`](tx.lock) | Pins the immutable **terminology answer pack** (content sha256 + URL) and, for this demo, the tooling jar. The repo now *names* exactly which terminology answers this commit builds against. |
| [`eng/future/build.sh`](eng/future/build.sh) | Resolves both by hash into content-addressed caches (`~/.fhir/tx-packs`, `~/.fhir/tools`), verifies, builds, checks the output signature, optionally judges byte parity. |
| [`eng/future/`](eng/future/) | The parity judge: a normalizing manifest hasher, the reference manifest, and the known-nondeterminism allowlist (the stock build is not byte-deterministic; that's [a reported bug](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-bugs.md), not something this branch hides). |
| [`.github/workflows/txpack-future.yml`](.github/workflows/txpack-future.yml) | CI: every push runs ONE hermetic build + signature check (~4 min, no terminology server anywhere). A separate reproducibility job (double build + judged byte parity) runs only when the pinned world changes or on demand. |

The tooling jar is one artifact: kindling (with the perf work: no forced GC per example,
parallel validation, indexed lookups, terminology fold) embedding the fhir-core txpack branch
(answer-pack seed layer, packager, hermetic mode, thread-safety fixes). Both live as reviewed
branches on [jmandel/kindling](https://github.com/jmandel/kindling) and
[jmandel/org.hl7.fhir.core](https://github.com/jmandel/org.hl7.fhir.core); design + evidence in
[jmandel/fhir-perf](https://github.com/jmandel/fhir-perf).

## A day in the life on this branch

**You edit a page, an example, a StructureDefinition — anything using terminology the spec
already uses.** `./eng/future/build.sh` (hermetic by default): every terminology question is
answered locally from the pack, byte-identically to what the server would say. Cold clone,
airplane, CI — all identical, all warm-speed.

**You add an example with three genuinely new SNOMED codes.** Build with `--online`: the three
misses go to the server (everything else still comes from the pack), and CI reports "this PR
introduces 3 new terminology questions" as a signal, not a gate. **Your PR contains only the
content change — you never touch `tx.lock`**, so nothing conflicts with anyone's long-lived
branch. After merge, the nightly refresh recording folds the three answers into a new pack and
bumps the lock in its own single-writer commit. (Hermetic mode stays available to *prove* the
no-new-codes case: it fails loudly naming any request that escapes.)

**The terminology server fixes a bug.** Nothing on your machine changes silently. A refresh
recording-run produces a new pack; the lock-bump PR renders exactly what changed ("12 SNOMED
display warnings resolved"); you get the fix when you pull, like every other change to the
spec. Today the same fix arrives never, or instantly, or per-machine — depending on cache
state nobody can inspect.

**Your build fails only on your machine at 3am.** It can't be the terminology cache anymore:
there is no mutable per-machine terminology state left to rot. Packs are immutable,
content-addressed, and shared.

## Commands, times, and what to expect

| Who | What they run | When | Takes | Expect |
|---|---|---|---|---|
| **Editor (any OS)** | `SpecBuild build .` (or `./eng/future/build.sh`) | every edit cycle | **~3.5 min** (12-core; ~4 min on a 4-core CI runner) | hermetic, zero terminology network, signature checked; first-ever run adds a one-time pack+jar fetch (~200MB) |
| Editor, adding new codes | `SpecBuild build . --online` | when the spec gains terminology | same + ~1s per new code | misses answered live and reported; PR carries content only — **never `tx.lock`** |
| Editor, checking blast radius | `SpecBuild build . --impact` | before pushing | build + ~30s | "your edit changed these 4 published files" — not a 226-file noise dump |
| **Content-PR reviewer** | nothing | — | — | content diff only; CI annotates "introduces N new terminology questions" |
| **CI, every push** | one `build` | per push | ~4 min | the entire everyday CI cost |
| **Refresh bot (the only `tx.lock` writer)** | recording run + `diff-packs` | nightly | one build | most nights: canonical hash unchanged → silence (a free daily server-consistency proof) |
| **Lock-bump reviewer (Grahame's seat)** | reads the PR | on real change | ~1 min | layered evidence: answer diff (machine), published-output impact, AI explanation (prose only); auto-merge for additions-only |
| **Anyone proving reproducibility** | `build --judge` (CI: double build) | lock/tooling bumps, optional nightly | 2× build + ~1 min | byte parity vs the committed reference, noise excused with per-run evidence |

Full CLI: `SpecBuild <build | manifest | compare | impact | diff-packs>` — one jar, one
command surface, identical on Windows/macOS/Linux (`java -Xmx12g -cp kindling-future-v5.jar
org.hl7.fhir.tools.publisher.SpecBuild help`).

### The Grahame story, end to end

His server fixes a display string overnight. The nightly recording sees it; the gate passes;
the canonical diff is non-empty for the first time in weeks; a lock-bump PR opens itself:
*"1 answer changed (SNOMED display for X); 3 published pages change, listed; explanation
attached."* He (or policy) merges it; every editor receives the fix with their next `git pull`
— visibly, atomically, identically. His server's total load for propagating the fix to the
entire world: **one build's worth of requests**. Today the same fix propagates never, or
instantly, or per-machine, depending on cache states nobody can inspect — while every cold
build anywhere hammers him with ~2,000 requests.

## What this demo deliberately keeps honest

- **Numbers**: warm 683s → ~197s and cold 1117s → ~211s were measured on a 12-core/62GB
  machine; GitHub's 4-core runners are slower in absolute terms for both stock and future.
  The headline that survives any hardware: **cold = warm, zero terminology requests, exact
  output parity**.
- **The reference manifest is pinned to this toolchain + pack configuration** (that's the
  txpack contract: pinned answers → pinned output). Against a build that asks tx.fhir.org live,
  ~49 pages differ in server-dataset metadata (e.g. one server annotates IANA-timezone bindings,
  the other doesn't) — exactly the silent drift class txpack exists to make visible and pinnable.
- **Scope**: the pack covers terminology traffic. FHIR *package* downloads (`~/.fhir/packages`)
  are a separate, already-content-versioned mechanism — CI caches them; a future `pkg.lock`
  could pin them the same way.
- **Byte-parity is the stretch goal, not the requirement.** The txpack goals (fast, offline,
  poison-free, pinned answers) need only hash-verified inputs + hermetic mode + the signature
  check — one build, no comparison. The reproducibility track (same commit → same bytes on any
  machine) is a separate ambition this demo also pursues; chasing it surfaced six classes of
  environment leak in the stock toolchain (locale → request keys, checkout path and OS username
  → page content, installed fonts → spreadsheet column widths, filesystem enumeration order →
  archive member order, thread timing → element ordering and a flickering designations table),
  each handled by pinning or by normalization in the judge.
- **How the judge stays honest about noise without going blind**: exact hashes are compared
  first; a file differing only under exact-but-not-order-insensitive hashing changed purely in
  element order (the documented stock ordering bug) and is excused with that evidence; the
  reproducibility job builds twice, and files that differed between its own two builds are
  excused as proven-nondeterministic-here; the small historical allowlist covers cross-machine
  ordering samples. A planted content change in an allowlisted file is still caught (verified).
- **The stock baseline** is a manual CI job (`workflow_dispatch`) so this repo doesn't hammer
  tx.fhir.org on every push.
- **Maven hosting**: GitHub Packages requires auth even for public reads, so the demo ships
  artifacts as content-addressed GitHub Release assets instead — which is also the more honest
  model for immutable, hash-named artifacts.

## Local modes and dependencies

`./eng/future/build.sh` (default): bash + curl + sha256sum + java, nothing else — and the
script itself is demo scaffolding for what would be ~30 lines of publisher code reading
`tx.lock` natively. `--impact`: after an edit, lists exactly which published files your change
touched (content-level; ordering churn excluded) — a capability the stock build cannot offer at
any price, since its output drowns intent in nondeterminism. `--judge`: verifies your
environment reproduces the pinned output (python3; used by CI and when bumping the lock).

## The refresh flow (prototyped here too)

Bumps are ordinary commits with a single writer (the refresh bot), and **the recording run is
the verification**: it is a fully-gated live build, so a bump inherits its trust rather than
re-earning it. [`txpack-refresh.yml`](.github/workflows/txpack-refresh.yml) prototypes the
pipeline: canonical pack comparison (`SpecBuild diff-packs` — volatile fields like server
step-timings and expansion timestamps are normalized, so "nothing really changed" is detectable
and ends the job silently), then a lock-bump PR whose body layers the machine-verified answer
diff (authoritative), the published-output impact, and an **AI-generated explanation section**
(GitHub Models, plain inference with the built-in token — no tools, no write access, prose
only). Merge policy keys off the machine facts, never the prose. *(In production the impact
evidence falls out of the recording run's own before/after at zero extra cost; this demo has no
recorder in CI, so the workflow synthesizes it with a same-machine old-vs-new A/B. The E2E test
of a regressing candidate — the previous pack offered as a "refresh" — was correctly blocked.)*

## Where the real proposal lives

- Design: [txpack proposal](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-proposal.md)
- The 14-bug report this work surfaced: [upstream bugs](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-bugs.md)
- PR sequencing for upstreaming all of it: [plan](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-plan.md)
