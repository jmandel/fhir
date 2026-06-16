# A peek into the txpack future

This branch (`txpack-future`) is a working demo of the proposed terminology + performance future
for the FHIR spec build. **A JDK is the entire dependency footprint, on any OS.** Try it:

```sh
git clone -b txpack-future https://github.com/jmandel/fhir.git && cd fhir
java -jar tools/build/launch.jar build .
```

That is a fully **cold** build that finishes in **~3.5 minutes** (12-core; ~4 on a small CI
runner) and makes **zero terminology network requests** — provably: any attempt is a hard
failure naming the request. The same build on the stock toolchain takes ~18+ minutes cold and
varies heavily with `tx.fhir.org`'s load (we measured 415–585s for the identical build at
different hours), failing outright when the server has a bad day.

## The idea, in one breath

The spec build's terminology answers become **an ordinary pinned dependency** — named in the
repo, immutable, fetched by hash, recorded once against the canonical server and replayed
deterministically forever. Three tiny artifacts, each bootstrapping the next, nothing
downloading itself:

| Artifact | Written by | What it does |
|---|---|---|
| [`tools/build/launch.jar`](tools/build/launcher-src/Launcher.java) | almost never | The gradle-wrapper move: the only thing you run. Fetches the pinned tooling (sha256-verified, cached), runs it. |
| [`tools/build/kindling-wrapper.properties`](tools/build/kindling-wrapper.properties) | **maintainers**, on toolchain releases | The toolchain pin (`toolUrl` + `toolSha256`). |
| [`fhir.lock`](fhir.lock) | **the refresh bot — its only writer** | The content lock: which terminology answers this commit builds against. |

Two writers — maintainer and bot — who never touch each other's files, so lockfile merge
conflicts are impossible by construction.

**Trust model in one line.** Verification happens once, at recording time (the recording *is* a
fully-gated live build); the pack replays those bytes deterministically forever; bumps are
ordinary single-writer commits that *inherit* that trust rather than re-earning it — and a
nightly recorder re-asks the live server and opens an evidence-laden PR when its answers
actually change.

## Everything else

The complete, current explanation — the build entry point and CLI, what a *recording* / *refresh*
/ *reproduce* / *the judge* / *from-scratch pin bootstrap* are, the answer pack and `fhir.lock`,
the CI workflows and their verification layers, what changed in each fork, and the determinism
contract — lives in one place:

➡️ **[`tools/build/README.md`](tools/build/README.md) — the complete txpack build guide.**

The full design rationale and failure-mode analysis live in the companion
[`docs/txpack-vision.md`](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md)
and [`docs/txpack-proposal.md`](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-proposal.md)
in the `jmandel/fhir-perf` workspace.
