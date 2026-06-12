# A peek into the txpack future

This branch is a working, self-contained demo of the proposed terminology + performance future
for the core FHIR spec build. Everything it claims is runnable here, now:

```bash
git clone -b txpack-future https://github.com/jmandel/fhir.git && cd fhir
./eng/future/build.sh --hermetic --judge
```

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

**You add an example with three genuinely new SNOMED codes.** Hermetic mode fails loudly,
naming the exact three requests (that's its job — proving completeness). Run
`./eng/future/build.sh --online`: the three misses go to the configured server (everything else
still comes from the pack), the build completes. A maintainer top-up run then folds the three
answers into a new pack and bumps `tx.lock` in a reviewed commit whose diff *is* the
terminology change.

**The terminology server fixes a bug.** Nothing on your machine changes silently. A refresh
recording-run produces a new pack; the lock-bump PR renders exactly what changed ("12 SNOMED
display warnings resolved"); you get the fix when you pull, like every other change to the
spec. Today the same fix arrives never, or instantly, or per-machine — depending on cache
state nobody can inspect.

**Your build fails only on your machine at 3am.** It can't be the terminology cache anymore:
there is no mutable per-machine terminology state left to rot. Packs are immutable,
content-addressed, and shared.

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

[`txpack-refresh.yml`](.github/workflows/txpack-refresh.yml) prototypes the routine update
path: canonical pack comparison ([`pack-diff.py`](eng/future/pack-diff.py) — volatile fields
like server step-timings and expansion timestamps are normalized, so "nothing really changed"
is detectable and ends the job silently), then a lock-bump PR whose body is layered:
the machine-verified diff table (authoritative), plus an **AI-generated explanation section**
(GitHub Models, plain inference with the built-in token — no tools, no write access, prose
only). Merge policy keys off the machine facts, never the prose.

## Where the real proposal lives

- Design: [txpack proposal](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-proposal.md)
- The 14-bug report this work surfaced: [upstream bugs](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-bugs.md)
- PR sequencing for upstreaming all of it: [plan](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-plan.md)
