# A peek into the txpack future

This branch is a working demo of the proposed terminology + performance future for the FHIR
spec build. **A JDK is the entire dependency footprint, on any OS.** Try it:

```
git clone -b txpack-future https://github.com/jmandel/fhir.git && cd fhir
java -jar eng/future/launch.jar build .
```

That is a fully cold build that finishes in ~3.5 minutes (12-core; ~4 on a small CI runner)
and **makes zero terminology network requests** — provably: any attempt is a hard failure
naming the request. The same build on the stock toolchain takes ~18+ minutes cold and varies
heavily with tx.fhir.org's load (we measured 415s–585s for the identical build at different
hours), failing outright when the server has a bad day.

## Three small artifacts, zero circularity

| Artifact | Written by | What it does |
|---|---|---|
| [`eng/future/launch.jar`](eng/future/launcher-src/Launcher.java) — **4KB, committed, source alongside** | almost never changes | The gradle-wrapper move: the only thing you run. Reads the wrapper pin, fetches the build tooling once (sha256-verified, cached in `~/.fhir/tools`), runs it with a sane heap (`HEAP` env to override). |
| [`eng/future/kindling-wrapper.properties`](eng/future/kindling-wrapper.properties) | **humans**, deliberately, on toolchain releases | The toolchain pin: `toolUrl` + `toolSha256` (exactly `gradle-wrapper.properties`' shape; with Maven-released tooling the URL is a pure function of the version). |
| [`fhir.lock`](fhir.lock) | **the refresh bot — its only writer** | The content lock: which terminology answers this commit builds against. |

Each stage bootstraps the next; nothing downloads itself. Two writers — maintainer and bot —
who never touch each other's files, so lockfile merge conflicts are impossible by construction
and review policies differ per file (toolchain bumps get changelog review; content bumps get
the evidence-laden PR described below).

### `fhir.lock` is not a new format

It is a profile of npm's `package-lock.json` v3 — a `packages` map of name →
`{version, resolved, integrity}` — because the *semantics* (single writer, integrity verified
on every use) are the contribution; the syntax should be one every developer already knows.
The terminology pack is simply the first entry, named `*.txpack`; ordinary FHIR package
dependencies (`hl7.terminology`, extensions — whose drift CI actually observed) can join the
same file with the same mechanism later. Two deliberate divergences, stated once: keys are
plain package names — **there is no vendor folder**; npm vendors per-project because it cannot
trust shared state, while here entries resolve into a shared content-addressed store
(`~/.fhir/tx-packs/<sha256>`) that is safe to share *because* every use re-verifies the hash —
and a top-level `expectedOutput` records the build signature this content produces, living in
the lock because it changes in the same bot commit as the pack it describes.

One line holds the whole trust model: **registries find bytes (`resolved`), the lock trusts
bytes (`integrity`), the store keeps bytes (immutable, verify-or-refetch).** Verification
happens once, at recording time — a recording is a fully-gated live build — and every
downstream consumer inherits that trust through the hash chain.

## How people work here

| Who | What they run | When | Takes | What happens |
|---|---|---|---|---|
| **Editor, any OS** | `java -jar eng/future/launch.jar build .` | every edit cycle | ~3.5 min | hermetic, zero terminology traffic, signature checked against the lock |
| Editor adding new codes | `… build . --online` | when the spec gains terminology | + ~1s per new code | only the new questions go to the server; the build reports "terminology questions not answered by the pack: N" |
| Editor checking blast radius | `… build . --impact` | before pushing | + ~30s | "your edit changed these 4 published files" — possible only because output is reproducible |
| **Content-PR author/reviewer** | nothing extra | — | — | PRs carry content only; **editors never touch `fhir.lock`**; the miss count surfaces in the CI log as a signal, not a gate |
| **CI, every push** | one `build` | per push | ~4 min | the entire everyday CI cost |
| **Refresh bot** (sole `fhir.lock` writer) | recording run + `diff-packs` | nightly | one build | most nights: canonical hash unchanged → silence (a free daily proof the server answers consistently); on change → lock-bump PR |
| **Lock-bump reviewer** (Grahame's seat) | reads the PR | on real change | ~1 min | layered evidence: machine answer-diff, published-output impact, AI-written explanation (prose only; merge policy keys off the machine facts) |

`build.sh` is a 3-line unix alias for the launcher. The CLI inside the tooling
(`SpecBuild build | manifest | compare | impact | diff-packs`) pins locale and timezone
in-process (both leak into request keys and published bytes otherwise) and detects an
undersized heap with the exact fix to apply.

### The Grahame story, end to end

His server fixes a display string overnight. The nightly recording sees it, passes its gate
(rc=0 + parity with its own previous output), and the canonical pack hash changes for the
first time in weeks. A PR opens itself: *"1 answer changed (SNOMED display for X); 3 published
pages change, listed; explanation attached."* It merges; every editor receives the fix at
their next `git pull` — visibly, atomically, identically. His server's total cost to propagate
the fix to the entire world: one build's worth of requests. Today the same fix arrives never,
instantly, or per-machine — depending on cache state nobody can inspect — while every cold
build anywhere sends him ~2,000 requests.

## The refresh pipeline (prototyped here)

Bumps are ordinary commits; **the recording run is the verification**, so no bump-time
re-verification ceremony exists. [`txpack-refresh.yml`](.github/workflows/txpack-refresh.yml)
exercises the downstream machinery: canonical pack comparison (`SpecBuild diff-packs` —
volatile fields like server step-timings and expansion timestamps normalized, so "nothing
really changed" is detectable and ends the job silently), then a PR layering the
machine-verified answer diff, the published-output impact, and an AI-written explanation
(GitHub Models, plain inference, no tools, no authority). *(In production the impact evidence
falls out of the recording run's own before/after at zero cost; this demo has no recorder in
CI, so the workflow synthesizes it with a same-machine A/B.)* CI-validated so far: the
no-change path (silent stop) and the harmful path — a candidate that removed answers was
blocked twice before any PR existed.

## The reproducibility track (optional, separate)

Nothing above needs output comparison. Separately, this branch pursues a stretch goal — *same
commit → same published bytes on any machine* — judged against a committed reference manifest
(`build . --judge`; CI runs it as the on-demand double-build job). Chasing it surfaced **seven
classes of environment leak in the stock toolchain** (locale → terminology request keys;
checkout path and OS username → page bytes; installed fonts → spreadsheet column widths;
filesystem enumeration order → archive member order; thread timing → element ordering and a
flickering designations table), each now pinned or normalized, several reported upstream. The
judge excuses *documented* nondeterminism with per-run evidence — a file must provably vary
between two same-commit builds to be excused, and order-only changes are detected structurally
— so a planted content change in a "noisy" file is still caught (tamper-verified). The stock
ordering bug still surfaces ~1–2 newly-sampled files per fresh environment; curation rule:
verify the diff is pure reordering/flicker, then allowlist.

## Deliberately honest footnotes

- Wallclock numbers are this hardware/runner; the claim that survives any hardware:
  **cold = warm, zero terminology requests, pinned output**.
- FHIR *package* downloads remain the un-pinned channel (their drift is real; we observed it).
  The fix is already designed: they join `fhir.lock` as ordinary integrity-pinned entries.
- The committed reference is pinned to this toolchain + pack configuration; against a live
  tx.fhir.org build, ~49 pages differ in server-dataset metadata — precisely the silent-drift
  class this design exists to make visible and pinnable.
- The stock-baseline CI job is manual-only, so this repo doesn't hammer tx.fhir.org.
- GitHub Packages was rejected for hosting (auth required even for public reads); the real
  distribution channel for packs is npm-format FHIR packages on the existing package servers
  (precedent: `hl7.fhir.rX.expansions`), with `fhir.lock` carrying the integrity the package
  infrastructure doesn't enforce. Tooling belongs on Maven, not npm.

## The full design

- [The settled vision](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md) — moving parts, user experiences, trust model
- [The design proposal](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-proposal.md) — failure-mode analysis of today's system
- [The 15-bug upstream report](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-bugs.md) this work surfaced
- [The upstreaming plan](https://github.com/jmandel/fhir-perf/blob/main/docs/upstream-plan.md)
