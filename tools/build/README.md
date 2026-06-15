# Building the FHIR spec (txpack): the 60-second quickstart

This is the front door for building the FHIR specification on this branch. The entire
dependency footprint is **a JDK** — no bash, Python, ant, curl, or git-svn. The build is
hermetic (zero terminology network traffic), fast (~3.5 min cold on 12 cores), and its output
is signature-checked against a committed lock.

## 1. Prerequisites

- A **JDK** (Java 17+). Nothing else.
- ~**12 GB of heap** available to the JVM (the build wants `-Xmx12g`).
- Run **from the spec checkout root** (the directory that contains `fhir.lock` and
  `tools/build/`). The launcher refuses to run from anywhere else.

## 2. The one command

```
java -jar tools/build/launch.jar build .
```

On a unix shell there is a thin alias that `cd`s to the repo root for you:

```
tools/build/build.sh
```

(`build.sh` is literally `cd "$(dirname "$0")/../.." && exec java -jar tools/build/launch.jar build . "$@"`.)
Both forms accept the same trailing flags, e.g. `tools/build/build.sh --impact`.

## 3. What happens

1. The launcher reads `tools/build/kindling-wrapper.properties` (`toolUrl` + `toolSha256`).
2. It downloads the pinned build tool **once** into `~/.fhir/tools/<sha256>.jar`, verifies the
   SHA-256, and refuses to run anything that doesn't match (exit 1). Subsequent builds reuse
   the cached jar.
3. It execs a child JVM (`<java.home>/bin/java -Xmx12g -cp <jar>
   org.hl7.fhir.tools.publisher.SpecBuild build .`).
4. `SpecBuild` runs a **hermetic cold build** (~3.5 min), pinning locale (en-US), timezone
   (America/Chicago), and UTF-8 so the output is reproducible.
5. It scrapes the build's `Summary: Errors=…, Warnings=…, Information messages=…` line and
   compares it to the signature pinned in `fhir.lock`. By default (gate = `enforce`) a mismatch
   fails the build (exit 1). A non-default report-only mode
   (`-Dorg.hl7.fhir.spec.signatureGate=report`, or `SIGNATURE_GATE=report` in the environment)
   instead prints `SIGNATURE CHANGED (reported, not enforced)` and returns 0 — see
   [BUILD-REFERENCE.md](BUILD-REFERENCE.md).

Output lands in:

| Path | What it is |
|---|---|
| `publish/` | the generated specification site |
| `build-future.log` | the Publisher build log (tee'd from the Publisher's stdout) |
| `build-future.manifest` | output fingerprint (only written with `--manifest`/`--judge`/`--impact`) |

`build-future.log` captures the Publisher's own stdout. SpecBuild's wrapper summary lines — the
`hermetic:` banner, `BUILD ok…`, `signature: … (expected …)`, and `terminology questions…` —
are printed to the console *after* the log is closed, so they appear on screen but are not in
the log file.

## 4. Hermetic is the DEFAULT

The plain `build` is **hermetic**: every terminology answer is served from the pinned pack, and
**any** attempt to reach a terminology server is a hard failure. On start you'll see:

```
hermetic: any terminology network attempt is a hard failure (use --online when adding new codes)
```

When your edit introduces genuinely new codes the pack hasn't seen, run with `--online` so pack
misses fall through to the live server:

```
tools/build/build.sh --online
```

`--online` does not silence anything — it counts the questions the pack couldn't answer and
reports them so you (and CI) can see the blast radius:

```
terminology questions not answered by the pack: N (these will fold into the pack at the next refresh)
```

Those new answers are folded into the pack by the nightly refresh bot; you don't edit the pack
or the lock yourself (see §7).

## 5. Everyday flags

All of these are arguments to `build` (passed after the command, e.g.
`tools/build/build.sh --impact`):

| Flag | What it does |
|---|---|
| `--online` | let pack misses fall through to the live server; report the miss count (use when adding new codes) |
| `--hermetic` | the explicit form of the default (any tx network attempt is a hard failure) |
| `--impact` | list exactly which published files changed vs your previous build (implies `--manifest`) |
| `--judge` | verify the published output against the committed `tools/build/ref.manifest` (implies `--manifest`) |
| `--manifest` | fingerprint `publish/` into `build-future.manifest` so a later `--impact`/`--judge` has a baseline |

`--impact` needs a prior `--manifest` build to diff against; the first time you'll see
`impact: no previous build manifest found (run 'build --manifest' first)`.

Unrecognized `-flags` (and `-fhir-settings <file>`) are passed straight through to the
underlying Publisher.

## 6. Heap / JVM tuning (env vars, not flags)

The launcher reads these from the **environment** — they are not `build` flags:

| Env var | Effect | Example |
|---|---|---|
| `HEAP` | sets `-Xmx` on the child JVM (default `12g`) | `HEAP=11g tools/build/build.sh` |
| `JAVA_OPTS` | extra JVM options, whitespace-split, prepended before `-cp` | `JAVA_OPTS="-XX:+UseParallelGC" tools/build/build.sh` |

If the heap is under ~10 GB, the build warns:
`WARNING: max heap is …MB; the spec build wants ~12GB`.

## 7. Editors never touch `fhir.lock`

`fhir.lock` is the **content lock** — it pins which terminology pack (answers) this commit
builds against, plus the expected output signature (`{errors, warnings, information}`). Its
**only** writer is the refresh bot. A content PR carries content only; the lock changes in a
separate, evidence-backed bot PR. Likewise, the toolchain pin
(`tools/build/kindling-wrapper.properties`) is maintainer-owned and bumped deliberately. As an
editor you change neither.

## 8. Where to go next

- **[BUILD-REFERENCE.md](BUILD-REFERENCE.md)** — the full command/flag surface
  (`build`, `manifest`, `compare`, `impact`, `diff-packs`, `record`, `reproduce`), every system
  property (including the `SIGNATURE_GATE` escape hatch from §3/§4), and exit codes.
- **[REFRESH-AND-CI.md](REFRESH-AND-CI.md)** — the nightly recorder, lock-bump PRs, and the CI
  workflows.
- **[../../FUTURE.md](../../FUTURE.md)** — the narrative front door: why this exists, the three
  artifacts, and the trust model.
- **[txpack-vision.md](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md)**
  — the settled design intent (moving parts, user experiences, trust model).
