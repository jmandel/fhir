# The txpack FHIR spec build

A fast, hermetic, reproducible build of the FHIR specification, driven by a pinned terminology
**answer pack**. This is the single source of truth for understanding and using the whole system.

- **Fast.** A stock cold build is 18+ minutes and network-bound; this one is ~4 minutes.
- **Hermetic.** Every terminology question (`validate-code`, `expand`, server capabilities,
  `findTxResource`, tx-registry `/resolve`) is answered from a pinned, content-addressed pack
  *before* any network call. By default **any** attempt to reach a terminology server is a hard
  failure — proving the build is offline-complete. Adding new codes the pack hasn't seen is the
  one expected exception: the failure names the missing answer and tells you to re-run with
  `--online` (see [§1.2](#12-hermetic-is-the-default), which asks the server *only* those new
  questions and reports the count). CI follows the same split — content PRs that add terminology
  are green with a "N new terminology questions" signal; strict offline-completeness is enforced
  on pin changes, not on every content push.
- **Reproducible.** Accidental nondeterminism is fixed at source and inherent volatility is
  normalized, so the same commit produces the same bytes (run-to-run variance went from ~21
  differing output files to ~0). The output is checked against a signature pinned in `fhir.lock`.
- **One dependency: a JDK.** No bash, Python, ant, curl, or git-svn.

For the full design rationale and trust model, see
[docs/txpack-vision.md](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md).

---

## 1. Quickstart

### Prerequisites

- A **JDK** (Java 17+). Nothing else.
- ~**12 GB of heap** for the JVM (the build wants `-Xmx12g`; it warns under ~10 GB but does not fail).
- Run **from the spec checkout root** — the directory containing `fhir.lock` and `tools/build/`.
  The launcher refuses to run from anywhere else (exit 2).

### The one command

```sh
java -jar tools/build/launch.jar build .
```

On a unix shell, `tools/build/build.sh` is a thin alias that `cd`s to the repo root first:

```sh
./tools/build/build.sh
```

The launcher is pure JDK, so the **same command runs on Windows**. In PowerShell:

```powershell
java -jar tools\build\launch.jar build .
```

or in `cmd.exe`:

```bat
java -jar tools\build\launch.jar build .
```

Heap and JVM options come from the **environment**, not from `build` flags (see [§1.4](#14-heap--jvm-tuning)):

```powershell
# PowerShell
$env:HEAP="11g"; java -jar tools\build\launch.jar build .
```

```bat
:: cmd.exe
set HEAP=11g && java -jar tools\build\launch.jar build .
```

> **This replaces the legacy `build.sh`.** The repo-**root** `build.sh` is the old entry point
> (bash + git-svn + `publish.sh` + ant, ~18+ min, Azure-oriented) and is bypassed entirely. The
> thing you run is `tools/build/build.sh`, a 4-line alias that `cd`s to the spec root and execs
> `java -jar tools/build/launch.jar build . "$@"`.

### 1.1 What happens (launcher → wrapper → SpecBuild)

The launcher (`tools/build/launcher-src/Launcher.java`, committed as `launch.jar`) is the
gradle-wrapper pattern. It does no terminology work; it resolves and verifies the toolchain pin
and execs a child JVM:

1. **CWD must be the spec root.** It reads `tools/build/kindling-wrapper.properties` (the
   "wrapper") as a relative path; if absent, prints
   `no tools/build/kindling-wrapper.properties - run from the spec checkout root` and exits **2**.
2. **Read the pin.** `toolUrl` + `toolSha256` are loaded; `toolSha256` must match `[0-9a-f]{64}`
   or it exits **2**.
3. **Download once, content-addressed, verified.** The cached jar is `~/.fhir/tools/<sha256>.jar`.
   If missing or hash-mismatched, it downloads `toolUrl` to a temp file, **sha256-verifies** it,
   and only then atomically installs it. A verification failure deletes the partial and exits **1**.
   (Up to 5 redirects are followed manually with the `Location` header used as-is, so signed URLs
   survive intact.)
4. **Exec the child JVM:**

   ```
   <java.home>/bin/java -Xmx<HEAP|12g> [JAVA_OPTS…] -cp ~/.fhir/tools/<sha>.jar \
       org.hl7.fhir.tools.publisher.SpecBuild <your args…>
   ```

   The child's exit code is the launcher's exit code. The launcher picks `java.exe` on Windows.

`SpecBuild build` then pins the environment in-process (locale `en-US`, timezone
`America/Chicago`, `file.encoding=UTF-8` — all three leak into request keys and published bytes),
sets `maxConcurrency=12` and `localFirst=true` if unset, runs the hermetic cold build, and checks
the output signature against `fhir.lock`.

Output lands in:

| Path | What it is |
|---|---|
| `publish/` | the generated specification site |
| `build-future.log` | the Publisher build log (tee'd from stdout) |
| `build-future.manifest` | output fingerprint (only with `--manifest`/`--judge`/`--impact`) |

> `build-future.log` captures the Publisher's own stdout. SpecBuild's own summary lines
> (`hermetic:`, `BUILD ok…`, `signature: …`, `terminology questions…`) print to the console
> *after* the log closes, so they appear on screen but not in the file.

### 1.2 Hermetic is the default

The plain `build` is hermetic: every terminology answer is served from the pack, and any network
attempt is a hard failure. On start you'll see:

```
hermetic: any terminology network attempt is a hard failure (use --online when adding new codes)
```

When your edit introduces genuinely new codes the pack hasn't seen, run with `--online` so pack
misses fall through to the live server. `--online` does not silence anything — it counts the
unanswered questions and reports the blast radius:

```sh
./tools/build/build.sh --online
```
```
terminology questions not answered by the pack: N (these will fold into the pack at the next refresh)
```

Those new answers are folded into the pack by the nightly recorder; you never edit the pack or
the lock yourself (see [§3](#3-the-answer-pack--fhirlock) and [§4](#4-the-ci-workflows--verification-layers)).

### 1.3 `build` flags and subcommands

`SpecBuild` dispatches on its first argument; an unknown command prints usage and exits **2**.

| Command | Purpose |
|---|---|
| `build [folder] [flags]` | the lock-driven build (hermetic by default); `folder` defaults to `.` |
| `manifest <publishDir>` | fingerprint a publish directory into a manifest |
| `compare <ref> <new> [-prev m] [-allowlist f]` | judge two manifests; exit **1** on unexplained diffs |
| `impact <prev> <new>` | content-level diff of two manifests (which files changed) |
| `diff-packs <old> <new>` | canonical answer-pack comparison; exit **3** on a real change |
| `record [folder] [-out d] [-bootstrap] [-fhir-settings f]` | cold recorder vs a live server ([§4](#4-the-ci-workflows--verification-layers)) |
| `reproduce -fresh A -confirm B -pinned P [-out d]` | reproduce-before-propose gate (separate JVM) |
| `help` / `--help` / `-h` | usage |

**`build` flags** (passed after the command, e.g. `./tools/build/build.sh --impact`):

| Flag | Effect |
|---|---|
| (default) / `--hermetic` | hermetic build: sets `-Dorg.hl7.fhir.tx.hermetic=true`; any tx network attempt is a hard failure |
| `--online` | pack misses fall through to the live server; counts and reports the miss count |
| `--manifest` | also fingerprint `publish/` into `build-future.manifest` (enables later `--judge`/`--impact`) |
| `--judge` | implies `--manifest`; after the build, `compare` the manifest vs `tools/build/ref.manifest` (with `-allowlist tools/build/noise-files-v2.txt`, and `-prev` if a prior build manifest exists) |
| `--impact` | implies `--manifest`; list which published files changed vs your previous build manifest |
| `-fhir-settings <file>` | passed through to the Publisher (terminology server config) |
| any other `-flag` | passed verbatim to `Publisher.main` |

Baked into every `build` regardless of flags: `fhir.lock` is required (else exit **2**); if
`tools/build/fhir-settings.json` exists and you didn't pass `-fhir-settings`, it is auto-added;
`-nosound -nopartial -folder <abs root>` are always passed; and the **signature gate** runs
(below). `--impact` needs a prior `--manifest` build — the first time prints
`impact: no previous build manifest found (run 'build --manifest' first)`.

**Exit codes:** `0` success / no change · `1` signature mismatch or build failure, or `compare`
unexplained diffs · `2` bad usage / missing `fhir.lock` / unusable wrapper props · `3` `diff-packs`
packs differ.

**The signature gate.** After the build, the line `Summary: Errors=E, Warnings=W, Information
messages=I` is scraped from the log and compared to `TxLock.expectedSignature(fhir.lock)` (the
lock's `expectedOutput`, currently **`0 / 3694 / 349`**). A mismatch returns **1**, unless the
gate is set to report-only via `-Dorg.hl7.fhir.spec.signatureGate=report` or `SIGNATURE_GATE=report`,
in which case it prints `SIGNATURE CHANGED (reported, not enforced)` and returns 0.

### 1.4 Heap / JVM tuning

The launcher reads these from the **environment** (not as `build` flags):

| Env var | Effect | Example |
|---|---|---|
| `HEAP` | sets `-Xmx` on the child JVM (default `12g`) | `HEAP=11g ./tools/build/build.sh` |
| `JAVA_OPTS` | extra JVM options, whitespace-split, inserted before `-cp` | `JAVA_OPTS="-XX:+UseParallelGC" ./tools/build/build.sh` |

Use `JAVA_OPTS` to override the `org.hl7.fhir.tx.*` system properties (e.g.
`JAVA_OPTS="-Dorg.hl7.fhir.tx.maxConcurrency=8"`).

---

## 2. Concepts

A handful of terms recur. Defined plainly:

- **A RECORDING** — a *cold, online* build (`SpecBuild record`) run against a **live** terminology
  server with recording on. It deletes the local tx cache and ignores the pinned pack, so it
  re-asks the build's **full** terminology question set and captures what the server says **today**.
  Running cold (not seeded from the pack) is what makes a server *fix* to an existing answer
  visible. The output is a candidate answer pack.

- **A REFRESH** — the upkeep loop that keeps the pinned pack current: a nightly recording, a
  carry-forward merge against the current pinned pack, a reproduce-before-propose confirmation,
  and — only on a confirmed change — a one-line `fhir.lock` bump PR. A frozen pack drifts
  (the spec grows new codes; the server fixes old answers); the refresh closes that gap.

- **reproduce-before-propose** — the server-flakiness defense. Before any change is proposed, a
  **second independent recording** runs in a **separate JVM** (the Publisher is not re-entrant
  in-process). `SpecBuild reproduce` keeps a delta-vs-pinned **only if both recordings confirm it**;
  a delta that appears in one but not the other is dropped as a load-induced flap. So a proposal
  is, by construction, a real server change — never a single-run hiccup.

- **The JUDGE** (`OutputManifest`, via `SpecBuild compare` / `build --judge`) — decides whether the
  **published spec output** changed in a way that matters. It fingerprints every file in `publish/`
  (one `HASH<TAB>relpath` line, sorted) and compares against a reference. Timestamps, UUIDs, and
  section numbers are normalized; pure line **reordering** is excused (an order-insensitive `O:`
  hash); a real content change flips the hash and flags. It excludes **nothing** by path on purpose
  (`EXCLUDED_PATHS` is a regex that matches nothing) — a blanket exclusion once hid a real bug.
  Volatility is excused only with *evidence* (`-prev`, a same-commit twin build) or the static
  `-allowlist` (`tools/build/noise-files-v2.txt`), never by pattern.

- **PIN BOOTSTRAP / from-scratch pinning** — creating the *very first* pinned pack when there is no
  prior pin. `SpecBuild record -bootstrap` records cold against the live server, skips
  carry-forward and diff (there is nothing to carry forward from), and emits the pure fresh
  recording as a first pin. The `full_pin` CI job then **verifies** it with a hermetic build (zero
  misses = complete) and proposes the new signature.

  > **Do not confuse this with the launcher's _tool bootstrap_.** "Tool bootstrap" = the launcher
  > downloading/verifying/caching the tool jar ([§1.1](#11-what-happens-launcher--wrapper--specbuild)).
  > "Pin bootstrap" = the recorder creating a first content pin from scratch. Different mechanisms.

---

## 3. The answer pack & `fhir.lock`

### Find / trust / keep

`fhir.lock` (repo root) is a **profile of npm's `package-lock` v3** — so there is no new format to
learn. It pins the content this checkout builds against:

```json
{
  "lockfileVersion": 3,
  "packages": {
    "hl7.fhir.r6.txpack": {
      "version": "20260615",
      "resolved": "https://github.com/jmandel/fhir/releases/download/txpack-store/txpack-fa8828e5.zip",
      "integrity": "sha256-+ogo5S8ClXr73ma1upOEAVjq418hz9VCTTTyw/SgEfY=",
      "contentSha256": "fa8828e52f02957afbde66b5ba93840158eae35f21cfd5424d34f2c3f4a011f6",
      "recordedAgainst": "tx.fhir.org …, 2026-06-15"
    }
  },
  "expectedOutput": { "errors": 0, "warnings": 3694, "information": 349 }
}
```

This is the **find / trust / keep** model in one line:

- **FIND** — registries find bytes (`resolved` is where to fetch the pack).
- **TRUST** — the lock trusts bytes (`integrity`, an SSRI sha256, verified on **every** use).
- **KEEP** — the content-addressed store keeps bytes (`~/.fhir/tx-packs/<sha256>.zip`, immutable;
  a hash mismatch can only mean local damage, so it is deleted and refetched once).

Two deliberate divergences from npm: keys are **plain package names** (no `node_modules/` vendor
prefix — entries resolve into a shared store, never a per-project folder), and a top-level
`expectedOutput` records the build-output signature this content produces (it changes in the same
single-writer commit as the pack it describes). `-Dorg.hl7.fhir.tx.lock=ignore` disables lock
resolution.

### Where artifacts live now (and it's all no-token)

| Artifact | Lives on | Pinned by | Writer |
|---|---|---|---|
| **Answer packs** (`txpack-*.zip`) | **this repo's own** `txpack-store` release (`jmandel/fhir`) | `fhir.lock` (`resolved` + `integrity`) | the refresh **bot** |
| **Tool jar** (`kindling-future-v13.jar`) | the `jmandel/fhir-perf` release | `tools/build/kindling-wrapper.properties` (`toolUrl` + `toolSha256`) | a **maintainer** |

Both are public read (plain `curl`), content-addressed, and verified-on-use. **No cross-repo PAT
is needed:** packs publish to this repo's *own* release using the default `GITHUB_TOKEN`. The
recorder→refresh handoff does **not** push a request file (a `GITHUB_TOKEN` push won't trigger a
workflow); instead the candidate rides a `refresh-request` **artifact** and the refresh workflow
chains off the record run via `workflow_run` ([§4](#4-the-ci-workflows--verification-layers)).

### The three committed artifacts and their single writers

| Artifact | Writer | Pins |
|---|---|---|
| `tools/build/launch.jar` | maintainer (recompiles `launcher-src/Launcher.java`) | the bootstrap logic |
| `tools/build/kindling-wrapper.properties` | maintainer, on toolchain cadence | the **toolchain** (`toolUrl` + `toolSha256`) |
| `fhir.lock` | the refresh **bot** only | the **content** (the pack + its `expectedOutput`) |

The split is the point: the maintainer owns the toolchain, the bot owns the content, and the two
never collide. **As an editor you change neither** — a content PR carries content only; the lock
moves in a separate, evidence-backed bot PR.

> The toolchain `toolUrl` and the lock `resolved` currently point at GitHub Releases. That is
> intentional demo scaffolding; the `package-lock` v3 profile is chosen precisely so the lock
> needs no new format when an npm-format package channel exists.

---

## 4. The CI workflows & verification layers

Three workflows live on `txpack-future`. `txpack-future.yml` gates builds (a pack-seeded build on
every push; strict hermetic on pin changes); `txpack-record.yml` + `txpack-refresh.yml` run the
refresh loop.

### The verification layers, at a glance

| Layer | What it ensures | When |
|---|---|---|
| **Pack-seeded build + new-terminology signal** | the spec builds; intentional new codes reach the server and are *reported*, not blocked | every push |
| **Hermetic 0-miss** | offline completeness — the pinned pack answers every question, zero tx network | on pin / tooling change |
| **Drift diff** (record vs pinned) | server drift — the pinned pack still matches reality | nightly |
| **Reproduce-before-propose** | flake filter — a proposed change reproduces independently | nightly, on a delta |
| **Judge vs `ref.manifest`** | output reproducibility — published bytes match the reference | on `[parity]` / dispatch |
| **Determinism gate** | local nondeterminism is empty — two same-commit builds are byte-identical (no excusal) | on `[determinism]` / dispatch |
| **Stock baseline** | sanity vs the live server / legacy timing | manual dispatch only |

Determinism is what makes every *output* comparison meaningful.

### `txpack-future.yml` — the build gates

| Job | Trigger | What it does |
|---|---|---|
| **`future-build`** | **every push** (required) | one pack-seeded build, `SIGNATURE_GATE=report HEAP=11g ./tools/build/build.sh --online` — a no-new-code push still touches the network zero times; intentional new codes reach the server and surface as a "N new terminology questions" signal; the signature is reported, not gated. The everyday CI cost. |
| **`pinned-world-hermetic`** | **pin / tooling change** (`fhir.lock` or wrapper) | strict hermetic cold build `HEAP=11g ./tools/build/build.sh` — zero tx network + signature, **both enforced**: proves the pin is offline-complete. Skipped on content-only pushes. |
| **`reproducibility`** | `[parity]` commit msg, or dispatch `parity=true` (`continue-on-error`) | two builds: `--manifest` (convergence + evidence) then `--judge` (compare vs `ref.manifest`, with allowlist + `-prev`). Non-blocking signal; *allows* expected differences. |
| **`determinism-gate`** | `[determinism]` commit msg, or dispatch `parity=true` | three builds (first discarded for convergence), then `SpecBuild compare mA mB` with **no allowlist, no `-prev`** — any non-ordering content diff between two same-commit builds **fails**. *Forbids* differences. |
| **`stock-baseline`** | dispatch `stock_baseline=true` | the legacy `./publish.sh` cold build vs `tx.fhir.org`, for timing comparison. Slow; run rarely. |

### `txpack-record.yml` — the nightly recorder (upstream half)

Cron `11 9 * * *` (09:11 UTC, off-peak for `tx.fhir.org`), plus dispatch with a `server` input.
The default `record` job:

1. Fetch the pinned tool jar (sha-verified against `kindling-wrapper.properties`).
2. **Record (build A)** — `SpecBuild record . -fhir-settings … -out candidate-A` (cold,
   carry-forward `merge([pinned, fresh])` with pinned first, so fresh supersedes on overlap and a
   pinned answer survives only where a transient failure dropped it — a flake never becomes a
   "removed"). A `candidate-A.fresh.zip` sidecar appears **only** on a real delta.
3. **Reproduce** — only if the sidecar exists: a second cold recording (`-out candidate-B`,
   separate JVM), then `SpecBuild reproduce -fresh candidate-A.fresh.zip -confirm
   candidate-B.fresh.zip -pinned fhir.lock -out candidate-pack`.
4. **Propose** — if `candidate-pack.zip` survived, name it `txpack-<sha[:8]>.zip`, upload it to
   **this repo's own `txpack-store` release** (default `GITHUB_TOKEN`, created on demand), and emit
   a `refresh-request` **artifact** (`candidate_url` + `candidate_sha256`).

Most nights: one cold build, no delta, silent. The second build runs only when build A saw a delta.

**The `full_pin` dispatch job** ([pin bootstrap](#2-concepts)): cold-record a **complete** pack with
`record . -bootstrap` against the public server, **select the pure-fresh pack**, **verify** it
builds the spec hermetically with `-Dorg.hl7.fhir.tx.pack=<pack> --hermetic` (zero misses = a hard
failure proves completeness, with `SIGNATURE_GATE=report` since a from-scratch pin sets a *new*
signature), and upload it plus a `pin-proposal.json` as an artifact. Because a GitHub runner can
only reach the public server, the pack is provably recorded against `tx.fhir.org`.

### `txpack-refresh.yml` — the lock-bump PR (downstream half)

Chained off the recorder via `workflow_run` (so the default `GITHUB_TOKEN` suffices — no PAT, no
triggering push); also manually runnable with explicit `candidate_url` / `candidate_sha256` inputs.
It mutates **only `fhir.lock`** (branch `txpack-bump-<date>`); the pack never lands in git.

1. **Resolve candidate + fetch packs** — read the `refresh-request` artifact (or dispatch inputs);
   pull `current.zip` (from `fhir.lock`, sha-verified against `integrity`) and `candidate.zip`
   (sha-verified against the request). No artifact = no-delta night = exit quietly.
2. **Canonical compare** — `SpecBuild diff-packs current.zip candidate.zip`: exit **0** = nothing
   changed (stop); exit **3** = a real change (continue). The rc gate.
3. **Output A/B impact** — build hermetically with `-Dorg.hl7.fhir.tx.pack=current.zip --manifest`,
   then again with the candidate, then `SpecBuild impact` the two manifests. Hermetic both times, so
   a candidate that **drops** a needed answer fails loudly. Reports exactly which files change — or
   "output-inert".
4. **Explanation (GitHub Models)** — an LLM writes a short prose summary **from the diffs only**.
   Narrative, never authority; it may fall back to "service unavailable" without blocking.
5. **Open the PR** — rewrite `fhir.lock`'s `resolved` + `integrity`, push, `gh pr create` with the
   machine diff + impact + AI narrative.

**Merge policy** keys on the *machine facts*, never the prose: **additions-only** candidates are
safe to auto-merge; **changed answers** get human review.

---

## 5. What changed in each fork, and why

Three Git forks cooperate, all on branch `txpack-future`: **core** (`org.hl7.fhir.core`),
**kindling**, and **this repo** (`jmandel/fhir`) which wires them together.

### core (`org.hl7.fhir.core`) — the terminology engine

- **Answer-pack seed layer.** `TerminologyCache` gains a read-only seed (from
  `-Dorg.hl7.fhir.tx.pack`) consulted *before* the mutable cache and never written back. Covers
  validate/expand answers, server capabilities, `findTxResource` externals, and tx-registry
  `/resolve` results. **Negative answers are first-class** (a recorded "not on the server" stops a
  re-ask); a half-loaded pack is a **hard error**; misses fall through (it's a seed, not a wall).
- **Hermetic mode.** `-Dorg.hl7.fhir.tx.hermetic=true` makes the single HTTP choke point throw
  `TxHermeticViolationError` on the first network attempt. It **extends `Error`, not `Exception`** —
  load-bearing, because tx clients catch `Exception` and would otherwise swallow the violation into
  a cached error.
- **Cache-key canonicalization.** Per-run synthetic labels (`profile-url` urn:uuid, `cache-id`) are
  normalized *in the key only*, so packs are replayable; gated on pack/recording mode so default
  runs keep stock keys.
- **`TerminologyCachePackager` + poison filtering.** Build/merge/verify/diff CLI (the cache dir
  format *is* the pack format). Transport flakes are filtered; deterministic server refusals are
  kept; canonical diff normalizes volatile fields; `reproduceFilter` keeps only confirmed deltas.
- **TxLock** — resolves `fhir.lock`, verifies `integrity` on every use, exposes the signature.
- **Local-first / memo / throttle** — synthesize byte-identical grammar/unknown-system answers
  locally; run-scoped expansion memoization; a fair request-throttle semaphore.
- **Determinism fixes** — no shared-designation mutation, thread-safe `ContextUtilities`/snapshot,
  order-independent unknown-CodeSystem shape, deterministic ShEx emission, and more.

### kindling — the build front end

- **`SpecBuild` CLI** — the OS-independent command surface (`build`/`manifest`/`compare`/`impact`/
  `diff-packs`/`record`/`reproduce`) that replaces shell scripting.
- **The recorder** (`record`, with `-bootstrap`) — the cold, online, carry-forward recording loop
  ([§4](#4-the-ci-workflows--verification-layers)).
- **Terminology-fold** — `lookupLoinc()` routes through `validateCode` so *every* terminology
  answer flows through the cache and is pack-replayable.
- **Determinism** — deterministic `fhir.ttl` (stable-order Jena graph), null-safe profile fetch,
  **serial snapshot pre-generation** before the parallel validation pool (closes a torn-read race),
  deterministic NamingSystem sort, and an **honest judge** that excuses nothing by path.

### this repo (`jmandel/fhir`) — the integration layer

- **`tools/build/launch.jar`** (+ `launcher-src/Launcher.java`) — the entry point.
- **`tools/build/kindling-wrapper.properties`** — the maintainer-owned toolchain pin.
- **`fhir.lock`** — the bot-owned content lock.
- **`tools/build/build.sh`** — the thin unix alias.
- **`tools/build/fhir-settings.json`** — the FHIR settings the build runs against.
- **`tools/build/ref.manifest`** + **`noise-files-v2.txt`** — the judge's reference and noise allowlist.
- **The three `.github/workflows/txpack-*.yml`** — the build/record/refresh CI pipeline.

---

## 6. The determinism contract (in brief)

Determinism is what makes every output comparison sound. It is layered:

- **Accidental nondeterminism — FIXED at source.** Shared-mutable-state, hash-set iteration order,
  CRC-derived file names, torn reads under the parallel pool, etc. (mostly in core, two in kindling).
- **Inherent volatility — NORMALIZED.** Clocks, UUIDs, section numbers, step timings are normalized
  where they can't be removed — in pack diffs (`VOLATILE_PATTERNS`) and in the output judge
  (the order-insensitive `O:` hash, timestamp/uuid normalization).
- **Server nondeterminism — DETECTED and DEGRADED.** Poison filtering drops transport flakes;
  carry-forward keeps a known-good answer rather than spuriously removing one; reproduce-before-propose
  drops any change a re-recording can't confirm.

The result: the same commit produces the same bytes, so the judge and signature gate compare real
content rather than noise. For the full treatment, see
[docs/txpack-vision.md](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md).
