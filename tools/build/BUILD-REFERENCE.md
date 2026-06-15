# txpack build & verification: command, flag, and property reference

The complete, ground-truthed control surface for the lock-driven ("future world") FHIR spec
build. This is the maintainer/operator/CI-author reference: every command, flag, system
property, and on-disk path below is taken from the committed sources, not from prose.

The whole developer experience is one line:

```sh
java -jar tools/build/launch.jar build .
```

or, on a unix shell, the equivalent alias:

```sh
./tools/build/build.sh
```

The launcher is pure JDK — no bash, curl, or python — so the same command runs on Windows.
In PowerShell (the launcher reads the `HEAP` and `JAVA_OPTS` env vars for heap/JVM tuning):

```powershell
$env:HEAP="11g"; java -jar tools/build/launch.jar build .
```

or in `cmd.exe`:

```bat
set HEAP=11g && java -jar tools\build\launch.jar build .
```

For *how it fits together* and *why*, see [HOW-IT-WORKS.md](HOW-IT-WORKS.md). For the refresh
bot and CI wiring, see [REFRESH-AND-CI.md](REFRESH-AND-CI.md). For the quick-start, see
[README.md](README.md) in this folder. The repo-root [README.md](../../README.md) is the
top-level quick start: it now *leads* with this txpack flow (`java -jar tools/build/launch.jar
build .`, the `fhir.lock` signature, `SIGNATURE_GATE`, the `HEAP` override) and keeps the stock
Gradle/Ant build below a clearly-marked `> Legacy / stock toolchain` section. (Only the
repo-root `build.sh` script is genuinely legacy — see
[Stale references](#7-stale-references-do-not-propagate).)

---

## 1. The three committed artifacts and their single writers

Three files in this repo drive everything. Each has exactly one writer; the two human-owned
files and the bot-owned file never collide.

| Artifact | Writer | What it pins |
|---|---|---|
| `tools/build/launch.jar` | a maintainer, by recompiling `tools/build/launcher-src/Launcher.java` | the bootstrap logic itself (3892 bytes, stable, rarely changes) |
| `tools/build/kindling-wrapper.properties` | a maintainer, deliberately, on toolchain-release cadence | the **toolchain**: `toolUrl` + `toolSha256` of the build jar |
| `fhir.lock` | the refresh **bot** only | the **content**: the terminology answer pack + its `expectedOutput` signature |

The split is the point: the maintainer owns the toolchain (`kindling-wrapper.properties`), the
bot owns the content (`fhir.lock`), and the launcher (`launch.jar`) is the small, stable thing
that resolves one and hands off to a tool that consumes the other. The bot never touches
`kindling-wrapper.properties`; the maintainer never hand-edits `fhir.lock`.

Current pins (verify against the live files; do not trust this table to stay current):

- `kindling-wrapper.properties`: `toolUrl=…/txpack-future-v1/kindling-future-v11.jar`,
  `toolSha256=ea27d7c0b1cab47a81bb7c69420fe0140279d5fcb8d4b2f5b930ea0a0748c2bc`
- `fhir.lock`: package `hl7.fhir.r6.txpack`, `version 20260614`,
  `expectedOutput {errors:0, warnings:3694, information:349}`

---

## 2. The launcher bootstrap chain

`Launcher.java` is the gradle-wrapper pattern. It does no terminology work; it resolves the
toolchain pin, verifies it, and execs a child JVM. Exact sequence (from
`tools/build/launcher-src/Launcher.java`):

1. **CWD must be the spec root.** It reads `tools/build/kindling-wrapper.properties` as a
   relative path; if that file is absent it prints
   `no tools/build/kindling-wrapper.properties - run from the spec checkout root` and exits **2**.
2. **Read the pin.** `toolUrl` and `toolSha256` are loaded. `toolSha256` must match
   `[0-9a-f]{64}`; otherwise exit **2**.
3. **Download-once, content-addressed.** The cached jar is `~/.fhir/tools/<sha256>.jar`. If it
   is missing or its sha256 does not match the pin, the launcher downloads `toolUrl` to a
   `<sha256>.jar.download` temp file, **sha256-verifies** it, and only then atomically renames it
   into place. A verification failure deletes the partial and exits **1** with
   `downloaded tooling does not match tx.lock tooling.sha256 - refusing to run it`.
   (The message text says "tx.lock"; the verification is purely against
   `kindling-wrapper.properties`. See [Stale references](#7-stale-references-do-not-propagate).)
4. **Redirects use `Location` as-is.** Up to 5 hops are followed manually if the JDK stops
   early; the `Location` header is resolved against the request URL but never decoded, so signed
   URLs carrying encoded query parameters survive intact. A non-200 final response throws.
5. **Exec the child JVM.** The command built is:

   ```
   <java.home>/bin/java -Xmx<HEAP|12g> [JAVA_OPTS…] -cp <~/.fhir/tools/<sha>.jar> \
       org.hl7.fhir.tools.publisher.SpecBuild <your args…>
   ```

   - Heap defaults to `-Xmx12g`; override with the `HEAP` env var (e.g. `HEAP=11g`).
   - `JAVA_OPTS`, if set, is split on whitespace and inserted before `-cp`.
   - The child's exit code is propagated as the launcher's exit code.

The launcher itself is OS-independent (it picks `java.exe` on Windows). A JDK is the only
dependency — no bash, curl, or python.

---

## 3. Subcommand table (`SpecBuild`)

`main()` in `SpecBuild.java` dispatches on the first argument. Anything else prints usage and
exits **2**.

| Command | Delegates to | Purpose |
|---|---|---|
| `build [folder] [flags]` | in-class → `Publisher.main` | the lock-driven build (hermetic by default); see §4 |
| `manifest <publishDir>` | `OutputManifest.run` | fingerprint a publish directory |
| `compare <ref> <new> [-prev m] [-allowlist f]` | `OutputManifest.run` | judge two manifests; exit **1** on unexplained diffs |
| `impact <prev> <new>` | `OutputManifest.run` | content-level diff of two manifests |
| `diff-packs <old> <new>` | `TerminologyCachePackager.main` (prepends `diff`) | canonical answer-pack comparison; exit **3** on real change |
| `record [folder] [-out d] [flags]` | in-class | cold refresh recorder (live server); see [REFRESH-AND-CI.md](REFRESH-AND-CI.md) |
| `reproduce -fresh A -confirm B -pinned P [-out d]` | in-class | reproduce-before-propose gate; run in a **separate JVM** |
| `help` / `--help` / `-h` | — | usage |

`record` and `reproduce` are real, dispatchable commands even though the class-level javadoc
(and `FUTURE.md`) only document `build/manifest/compare/impact/diff-packs`. They are the two
halves of the recorder feeding the refresh pipeline.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | success / no change (also: `diff-packs` packs canonically identical; `compare` clean) |
| `1` | `build` signature mismatch (when gate enforces) or build failure; `compare` unexplained diffs |
| `2` | bad usage / unknown command / missing `fhir.lock` / unusable wrapper props |
| `3` | `diff-packs`: packs differ (a real, non-volatile change) |

`diff-packs` and `compare` are the two machine-readable verdicts a refresh or determinism job
keys on: `diff-packs` is `0`/`3` (identical/changed), `compare` is `0`/`1` (clean/unexplained).

---

## 4. `build` flags

```sh
java -jar tools/build/launch.jar build . [--online|--hermetic] [--judge] [--manifest] [--impact] \
     [-fhir-settings <file>] [-other -publisher -flags…]
```

`build [folder]` defaults `folder` to `.`. A bare token (not starting with `-`) is taken as the
folder; any `-flag` not recognized below is passed through to `Publisher.main`.

| Flag | Effect |
|---|---|
| (default) / `--hermetic` | hermetic build: sets `-Dorg.hl7.fhir.tx.hermetic=true`; **any** terminology network attempt is a hard failure |
| `--online` | pack misses fall through to the live server; enables miss counting (logs to a temp NDJSON via `org.hl7.fhir.tx.logMisses`, reports "terminology questions not answered by the pack: N") |
| `--manifest` | also fingerprint `publish/` into `build-future.manifest` (enables later `--judge`/`--impact`) |
| `--judge` | implies `--manifest`; after the build, runs `compare` of the manifest vs `tools/build/ref.manifest`, with `-allowlist tools/build/noise-files-v2.txt` and (if a previous build manifest exists) `-prev <that>` |
| `--impact` | implies `--manifest`; lists which published files changed vs your previous build manifest |
| `-fhir-settings <file>` | passed through to the Publisher (terminology server config) |
| any other `-flag` | passed through verbatim to `Publisher.main` |

Behavior baked into `build()` regardless of flags:

- **`fhir.lock` is required.** If `<folder>/fhir.lock` is absent, it prints
  `no fhir.lock … - this command drives lock-pinned builds; use the stock publisher otherwise`
  and returns **2**.
- **Auto fhir-settings.** If `tools/build/fhir-settings.json` exists and you did not pass
  `-fhir-settings`, it adds `-fhir-settings <abs path to that file>` (which pins
  `txFhirProduction = https://tx.fhir.org`).
- **Pinned environment, in-process.** `Locale.US`, timezone `America/Chicago`, `file.encoding=UTF-8`
  — all three leak into request keys and published bytes otherwise, so they are set
  programmatically (no JVM flags needed beyond heap).
- **Concurrency defaults.** If unset, `org.hl7.fhir.tx.maxConcurrency=12` and
  `org.hl7.fhir.tx.localFirst=true` are applied (both only when not already specified).
- **Heap warning.** If max heap `< 10GB` it warns (`the spec build wants ~12GB — add e.g.
  -Xmx12g`). This is advisory, not fatal.
- **Always passes `-nosound -nopartial`** to the Publisher, plus `-folder <abs root>`.
- **Signature gate.** Output is teed to `build-future.log`; the line
  `Summary: Errors=E, Warnings=W, Information messages=I` is scraped and compared to
  `TxLock.expectedSignature(fhir.lock)` (the lock's `expectedOutput`). On mismatch the build
  returns **1**, *unless* the gate is set to report-only (see below), in which case it prints
  `SIGNATURE CHANGED (reported, not enforced)` and returns 0. A missing signature line also
  returns **1**.

The signature gate is controlled by `org.hl7.fhir.spec.signatureGate` (system property) or the
`SIGNATURE_GATE` env var; the value `report` makes it non-enforcing. The default is `enforce`.

> Note: the gate compares against `TxLock.expectedSignature` (currently
> `Warnings=3694, Information=349`). Some CI comment strings elsewhere hardcode older numbers;
> the lock is authoritative. See [Stale references](#7-stale-references-do-not-propagate).

---

## 5. System properties (`org.hl7.fhir.tx.*`) and env vars

These are read by the core/kindling runtime. `build` sets several of them for you (§4); set them
yourself via `JAVA_OPTS` (e.g. `JAVA_OPTS="-Dorg.hl7.fhir.tx.maxConcurrency=8"`) when overriding.

| Property | Default | Effect |
|---|---|---|
| `org.hl7.fhir.tx.pack` | unset | path (dir or zip) to a terminology answer pack loaded read-only before the mutable cache. Setting it also activates cache-key canonicalization. |
| `org.hl7.fhir.tx.hermetic` | `false` | `true` makes every FHIR-server HTTP request throw `TxHermeticViolationError`. `build` sets this in the default/`--hermetic` mode. |
| `org.hl7.fhir.tx.lock` | (resolve) | `ignore` disables `fhir.lock` pack resolution entirely (`TxLock.LOCK_SYSTEM_PROPERTY`). |
| `org.hl7.fhir.tx.logMisses` | unset | path to an NDJSON file; pack misses are appended. `build --online` sets this to a temp file to count new questions. |
| `org.hl7.fhir.tx.localFirst` | `false` (runtime); `true` under `build` | answer UCUM/mimetype/BCP-47 grammar shapes locally when byte-identical to the server. |
| `org.hl7.fhir.tx.maxConcurrency` | `4` (runtime); `12` under `build` | permits in the JVM-wide fair request-throttle semaphore. Absent/unparseable/≤0 → 4. |
| `org.hl7.fhir.tx.recordSemanticErrors` | `false` | persist semantic `CODESYSTEM_UNSUPPORTED` answers (and shadow-send suppressed requests) for an exhaustive pack. Also activates key canonicalization. Used by `record`. |
| `org.hl7.fhir.tx.adaptiveConcurrency` | `false` | *(core-only)* `true` enables the opt-in AIMD `AdaptiveThrottle` (initial max 64, or `maxConcurrency` if set). Off by default — the 404-as-load-shedding heuristic is wrong for general traffic. |
| `org.hl7.fhir.tx.evalMemo` | `true` | *(core-only)* run-scoped expansion memoization; set `false` to disable. |
| `org.hl7.fhir.spec.signatureGate` | `enforce` | `report` makes the `build` signature gate non-enforcing (env equivalent: `SIGNATURE_GATE`). |

Env vars consumed by the launcher (§2): `HEAP` (heap size, default `12g`) and `JAVA_OPTS`
(extra JVM flags inserted before `-cp`).

---

## 6. The `~/.fhir` layout

Two shared, content-addressed stores under the user's home. Both are immutable and
verified-on-use; a corrupt entry is deleted and refetched, never silently trusted.

| Path | Holds | Verified by |
|---|---|---|
| `~/.fhir/tools/<sha256>.jar` | the build toolchain jar | `Launcher.java` — sha256 must equal `kindling-wrapper.properties` `toolSha256` |
| `~/.fhir/tx-packs/<sha256>.zip` | terminology answer packs | `TxLock.resolvePackPath` — sha256 must equal the SSRI `integrity` in `fhir.lock`, **on every use** |

Because both stores are keyed by content hash, a failed verification can only mean local
damage (never a legitimate update): `TxLock` deletes the bad `<sha>.zip` and refetches once from
the lock's `resolved` URL; the launcher does the same for the jar. Multiple checkouts share both
stores — the pack/tool is downloaded once per content hash, machine-wide.

---

## 7. Stale references (do not propagate)

A few strings in the tree still say `tx.lock`; the real committed file is **`fhir.lock`**.

- `Launcher.java`'s class javadoc and its mismatch message
  (`… does not match tx.lock tooling.sha256 …`) predate the rename. The launcher only ever
  reads `tools/build/kindling-wrapper.properties` — it never touches any `tx.lock`/`fhir.lock`.
- `TxLock.LOCK_SYSTEM_PROPERTY` javadoc says "when a tx.lock file is present"; the file is
  `fhir.lock`. The property itself (`org.hl7.fhir.tx.lock=ignore`) is correctly named.
- The refresh workflow hardcodes an older comparison string (`Warnings=3693 / Info=345`); the
  enforced gate uses `TxLock.expectedSignature` from the lock (currently `3694 / 349`). The lock
  is authoritative; the comment is cosmetic.
- The repo-**root** `build.sh` is the **legacy** entry point (bash + git-svn + `publish.sh` +
  ant, ~18+ min, Azure-oriented) and is bypassed entirely by the txpack flow. The thing you run
  is `tools/build/build.sh`, a 4-line unix alias that `cd`s to the spec root and execs
  `java -jar tools/build/launch.jar build . "$@"`. The repo-root **README**, by contrast, is
  *not* legacy: it now documents the txpack build up top (leading with
  `java -jar tools/build/launch.jar build .`) and demotes the stock Gradle/Ant toolchain to a
  section explicitly marked `> Legacy / stock toolchain`.

---

### See also

- [README.md](README.md) — quick start for this build
- [HOW-IT-WORKS.md](HOW-IT-WORKS.md) — the determinism contract and trust model end to end
- [REFRESH-AND-CI.md](REFRESH-AND-CI.md) — the recorder, the refresh bot, and the CI workflows
- [../../README.md](../../README.md) — the repo-root quick start; leads with the txpack build
  and keeps the stock Gradle/Ant toolchain below a marked Legacy section
