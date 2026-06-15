# How the txpack build works: answer pack, lock, determinism, and the fork changes

This is the mechanism document. It is for reviewers and engineers who want to know *what changed
in each fork and why*, and how the pieces fit. If you just want to build the spec, read
[README.md](README.md); for a flag/file reference read [BUILD-REFERENCE.md](BUILD-REFERENCE.md);
for the nightly pipeline and CI gates read [REFRESH-AND-CI.md](REFRESH-AND-CI.md). For intent and
the trust model in prose see the repo-root [FUTURE.md](../../FUTURE.md). The legacy build path
described in the root [README.md](../../README.md) (gradle/ant/Azure) is bypassed entirely.

Three Git forks cooperate: **core** (`org.hl7.fhir.core`, branch `txpack-future`), **kindling**
(branch `txpack-future`), and **this repo** (`jmandel/fhir`, branch `txpack-future`, which wires
them together). Everything below is grounded in those branches.

---

## 1. The problem this solves

A stock cold spec build is slow (18+ minutes), depends on a live terminology server
(`tx.fhir.org`), and is not byte-reproducible:

- **Slow & network-bound.** Every unknown code, value-set expansion and capability probe is an
  HTTP round-trip to `tx.fhir.org`, serialized through a small connection budget.
- **Flaky.** `tx.fhir.org` sheds load with 404/429/503; a transient failure becomes a cached
  `"Error performing tx ..."` answer and silently corrupts the run.
- **Non-reproducible.** Two builds of the same commit differed (measured ~21 differing output
  files run-to-run) from a mix of accidental nondeterminism (shared mutable state, hash-set
  iteration order) and inherent volatility (clocks, UUIDs).

The txpack approach answers terminology from a pinned, content-addressed **answer pack** before
ever touching the network; runs the build **hermetic by default** (any network attempt is a hard
failure, proving offline-completeness); and makes the output **deterministic** by fixing the
accidental sources and normalizing the inherent ones. See FUTURE.md and
`docs/txpack-proposal.md` for the failure-mode motivation in full.

---

## 2. The answer pack & seed layer (CORE)

`TerminologyCache` (`org.hl7.fhir.r5/.../terminologies/utilities/TerminologyCache.java`) gains a
**read-only seed layer** loaded from a pinned pack:

```
-Dorg.hl7.fhir.tx.pack=<packDir|packZip>
```

The pack loads into immutable maps that are consulted **before** the mutable on-disk cache and are
**never written back**. `packLookup()` is the first thing `getValidation()` / `getExpansion()`
check; a hit increments `packHitCount` and returns immediately. The pack covers:

| Pack content | Loaded into | Purpose |
|---|---|---|
| validate-code / expand answers (`*.cache` pages) | `packCaches` | the bulk terminology answers |
| server `CapabilityStatement` / `TerminologyCapabilities` + `servers.ini` | `packCapabilityStatements`, `packTerminologyCapabilities` | client init never re-fetches them |
| findTxResource externals (`vs-externals.json` / `cs-externals.json` + `vs-<uuid>.json` / `cs-<uuid>.json`) | `packVsExternals`, `packCsExternals` | externally-resolved ValueSets/CodeSystems |
| tx-registry resolutions (`system-map.json`) | `packSystemMapSource` | `/resolve` results, so the registry is never re-queried |

Two properties matter for hermetic correctness:

- **NEGATIVE answers are first-class.** An externals entry whose value is JSON `null` is a
  recorded *"this canonical is not on the server"* answer. It loads as a map entry with a `null`
  value (`loadPackExternals`) and is what stops a pack-seeded run from re-asking the server for a
  resource the recording run already learned it does not have.
- **Pack-read failure is a HARD error.** `loadPack` throws `IOException` if the pack is missing, or
  if an externals index references a missing per-resource file. Hermetic runs depend on the pack
  being complete and self-describing, so a half-loaded pack must fail loudly, not degrade.

**Misses fall through.** If a request is not in the pack, lookup returns `null` and the request
proceeds exactly as it would without a pack — to the mutable cache, then (in non-hermetic mode)
the network. The pack is a seed, not a wall.

---

## 3. Cache-key canonicalization

The cache/pack key is a hash of the request JSON. Two synthetic identity labels in that JSON are
per-run random, so without intervention they make every key run-unique and **unpackable**.
`canonicalizeRequest()` rewrites only those labels **in the key** (the stored/wire text is never
touched):

- **Rule 1 — `profile-url`.** A `urn:uuid`-shaped `profile-url` expansion parameter is a synthetic
  label naming the parameter set; some harnesses mint a fresh UUID per run. It is replaced by a
  fixed placeholder (`PROFILE_URL_UUID_PLACEHOLDER =
  urn:uuid:00000000-0000-0000-0000-000000000000`). The rule is restricted to `urn:uuid` shapes so
  a real `http(s)` profile-url that names a server-resolvable expansion profile is left alone.
- **Rule 2 — `cache-id`.** A per-session random token the client mints so the server can correlate
  previously-sent resources; it is a transport optimization, never semantic. Replaced by the same
  placeholder. **Rule 2 alone fixed 24 hermetic violations per run** — all narrative-path
  validate-codes whose only run-varying content was the `cache-id`.

Both rules are collision-safe: every actual semantic parameter is serialized in full in the same
request text, so two genuinely different requests still differ after normalization.

Canonicalization is **gated by `canonicalKeysActive()`**, which is true only when a pack is
configured (`-Dorg.hl7.fhir.tx.pack` set) **or** a recording run is active
(`recordSemanticErrors`). Default user runs keep stock key derivation — otherwise a routine
upgrade would silently re-key (and invalidate) every existing on-disk cache. The same
`cacheKeyFor()` is used by token generation, mutable-cache load and pack load, so build-time and
lookup-time keys can never disagree.

---

## 4. Hermetic mode

```
-Dorg.hl7.fhir.tx.hermetic=true
```

In `ManagedFhirWebAccessor` (`org.hl7.fhir.utilities/.../http/ManagedFhirWebAccessor.java`), the
single choke point all FHIR/terminology HTTP traffic flows through, `httpCall()` throws on the
first network attempt:

```java
if (HERMETIC) {
  throw new TxHermeticViolationError(describeBlockedRequest(httpRequest));
}
```

`TxHermeticViolationError` **extends `Error`, not `Exception`** — and that is load-bearing.
Terminology clients catch `Exception` and convert failures into cached `"Error performing tx ..."`
results; an `Exception` here would be swallowed into the mutable cache and *hide* the violation. An
`Error` propagates and fails the run loudly at the first attempt, naming the blocked request (and
its body, truncated to 4 KB).

To keep init traffic from tripping the gate, client bootstrap paths skip their version-pings and
clears under pack/hermetic mode. `TerminologyClientManager` / `TerminologyCacheManager` gate on
`isTxPackMode()` (true when `tx.pack` is set or `tx.hermetic=true`): they keep whatever version
stamp the recording run left, and a server whose CapabilityStatement the pack carries is treated as
already-described (`isPackDescribedServer`) so no probe is sent.

---

## 5. The packager

`TerminologyCachePackager` (`org.hl7.fhir.r5/.../terminologies/utilities/`) is a build/merge/
verify/diff CLI. **The on-disk cache directory format *is* the pack format**, so packaging is
mostly filtering and copying.

```
TerminologyCachePackager build <sourceCacheDir> <outputParentDir>
TerminologyCachePackager merge <cacheDir1> <cacheDir2> [...] <outputParentDir>
TerminologyCachePackager verify <packDirOrZip> [sampleCount]
TerminologyCachePackager diff  <oldPackDirOrZip> <newPackDirOrZip>
```

(From kindling's front door these are reached as `SpecBuild diff-packs ...`, which prepends `diff`
to the packager's arguments.)

### Poison vs kept refusals — one choke point

The packager classifies every answer at a single `PackPageWriter.add()` point, which throws
`PoisonEntryException` rather than write a poison entry. So a pack page can only ever be assembled
from clean entries.

- **Poison (filtered):** transport/transient failures — `TRANSPORT_MARKERS`
  (`SocketTimeoutException`, `Read timed out`, `Connection reset/refused`, `UnknownHostException`,
  …) and the client-side wrappers `HTTP_WRAPPER_MARKERS` (`"Error from http"`, `"Error performing
  tx"`). These wrap whatever went wrong on the wire and must never be replayed.
- **Kept (replayable):** `DETERMINISTIC_REFUSAL_MARKERS` — server-authored refusals that are
  deterministic for a pinned server+content edition: `"has a grammar, and cannot be enumerated"`,
  `"too costly to expand"`, `"could not be found, so the value set cannot be expanded"`. These
  reach the cache wrapped in `"Error from http..."`, so a naive wrapper-only filter would drop them
  and leave permanent per-run `$expand` traffic. `poisonMarkerIn()` checks transport markers
  **first and unconditionally**, then treats a wrapper marker as poison **only if** it is not also a
  deterministic refusal.

### Canonical diff & reproduceFilter

`diffPacks` compares two packs **canonically**, after `canonicalize()` normalizes the fields that
legitimately vary between two recordings of the same server (`VOLATILE_PATTERNS`):

| Pattern | Replacement | What it normalizes |
|---|---|---|
| `(?<![0-9])\d+ms ` | `Nms ` | server step-timings in diagnostics |
| `urn:uuid:[0-9a-f-]{36}` | `urn:uuid:X` | expansion identifiers |
| ISO-8601 timestamp | `TS` | expansion timestamps |

`diff` exits **0** when canonically identical and **3** on a real change (markdown summary on
stdout). `reproduceFilter` keeps a delta only if a *second* recording reproduces it canonically —
an un-reproduced "change" is a server flake, not a real change.

The timing regex must be `(?<![0-9])\d+ms` — **a `\b` word boundary fails here**: the timings sit
after a JSON-encoded `\n`, so the character before the digit is the word-char `n`, and `\b` matches
only at the very first timing. The non-digit lookbehind is what catches every timing; **this one
fix removed 668 phantom "changes."**

`system-map.json` is written with Gson `serializeNulls()` on purpose: the negative resolutions are
stored as JSON `null`, and without `serializeNulls` Gson silently drops every negative entry, so a
hermetic run would re-ask the network.

---

## 6. Shadow recording for exhaustive packs

```
-Dorg.hl7.fhir.tx.recordSemanticErrors=true
```

Two gaps would otherwise make a recorded pack incomplete:

1. **Dropped semantic errors.** The default cache policy never persists a server's definitive
   `CODESYSTEM_UNSUPPORTED` ("I don't know this code system") answer for an *unversioned* request.
   With recording on, `store()` rescues these via `isRecordableSemanticError()` — but the class
   check is not trusted alone: the decision is made on the **exact bytes the entry would persist**,
   re-checked against `TerminologyCachePackager.isPoison`, so a recorded entry can never be one the
   packager would later have to filter.

2. **One-shape-per-system from memo suppression.** During a run the per-run unknown-system memo
   answers every probe after the first *locally*, so the recorded pack would carry only one request
   shape per unknown system; a replay whose thread timing differs hits a different shape first,
   misses, and goes to the server (~91 residual requests/run measured). At each suppression point,
   when recording is active, `shadowRecordSuppressedCodingValidation` /
   `shadowRecordSuppressedBatchValidation` synchronously send the **exact request the un-suppressed
   path would have sent** through the normal server machinery and store its answer — but never
   return it to the caller, so observable behavior is unchanged. The cache-token lookup precedes the
   shadow write, so each unique shape is shadow-sent at most once. Result: the pack captures every
   request shape per unknown system regardless of thread timing.

---

## 7. `fhir.lock`, TxLock, and the trust model

`fhir.lock` (repo root) is a **profile of npm's `package-lock` v3** — chosen so there is no new
format to learn. It pins the content the checkout builds against:

```json
{
  "lockfileVersion": 3,
  "packages": {
    "hl7.fhir.r6.txpack": {
      "version": "20260614",
      "resolved": "https://.../txpack-78abc477.zip",
      "integrity": "sha256-UF5bLyQYSzc2WLA484JTvuk374Xe1bP8QEhMzqcfMk0=",
      "contentSha256": "78abc4773bb45c4f526c5e4adcb701e147b00810e2aa2579d9af0c76c3931906",
      "recordedAgainst": "tx.fhir.org ..., 2026-06-14"
    }
  },
  "expectedOutput": { "errors": 0, "warnings": 3694, "information": 349 }
}
```

Two deliberate divergences from npm: keys are **plain package names** (no `node_modules/` vendor
prefix — entries resolve into a shared store, never a per-project folder), and a top-level
`expectedOutput` records the build-output signature this content produces (it changes in the same
single-writer commit as the pack it describes).

`TxLock` (`org.hl7.fhir.r5/.../terminologies/utilities/TxLock.java`) resolves the `*.txpack` entry:

- `resolvePackPath()` downloads `resolved` into the content-addressed store
  `~/.fhir/tx-packs/<sha256>.zip` and **verifies the sha256 against the SSRI `integrity` on every
  use**. A mismatch on an existing store entry can only mean local damage (the store is
  content-addressed), so the entry is deleted and refetched once; a freshly downloaded mismatch is
  a hard error.
- `expectedSignature()` exposes the lock's pinned `{errors, warnings, information}` for the build's
  signature gate.
- `-Dorg.hl7.fhir.tx.lock=ignore` disables lock resolution.

This is the **find / trust / keep** model in one line: registries **FIND** bytes (`resolved`), the
lock **TRUSTS** bytes (`integrity`, verified on every use), the content-addressed store **KEEPS**
bytes (immutable, verify-or-refetch). There are two non-colliding writers: the maintainer owns
`tools/build/kindling-wrapper.properties` (the toolchain pin — `toolUrl` + `toolSha256`, resolved
the same way into `~/.fhir/tools/<sha256>.jar`), and the refresh bot owns `fhir.lock` (the content
pin).

---

## 8. KINDLING build changes

### Terminology-fold

`BuildWorkerContext.lookupLoinc()` no longer self-constructs a `tx.fhir.org` client and issues a
direct `$lookup`. It routes through `super.validateCode(Coding)`:

```java
vr = super.validateCode(new ValidationOptions(),
        new Coding().setSystem("http://loinc.org").setCode(code), null);
```

so **every** terminology answer flows through `TerminologyCache` and is pack-replayable. A missing
master client is now an `Error` (a real build always configures one) rather than something papered
over by a default client; infrastructure-class failures are distinguished from authoritative
"code not found" so `serverOk` (which drives WARNING-vs-ERROR downstream) stays correct.

### Parallel prefetch + serial snapshot pre-gen

- `Publisher.prefetchExpansions()` (called first thing in `produceSpec()`) warms the expansion
  cache in parallel (`fhir.build.expansion.threads`, default `min(cores, 12)`) over exactly the set
  page production will expand, using the identical `expandVS(...)` entry point so the cache token
  matches. Failures are swallowed — page production will redo the identical call serially and
  surface any error the stock way.
- In `validateFiles()`, every validated profile's snapshot is **pre-generated serially** before the
  parallel example-validation pool starts, so the parallel phase only ever *reads* finished
  snapshots. This closes a data race: snapshot generation mutates the shared cached
  `StructureDefinition` and its base in place, which under the parallel pool produced an
  intermittent torn read (~1 build in 3 showing `Errors=1`, never at `threads=1`).

> **Superseded spike.** `prefetchExpansions` is the integrated alternative that *replaced* an
> earlier never-integrated two-pass approach. The system property `-Dfhir.build.tx.twopass` is
> **absent on kindling `txpack-future`** (it lives only on core's `spike/core-batch-tx`). Do not
> document it as a current flag.

---

## What the fhir repo fork adds

The core and kindling sections above cover the *engine* changes. This fork (`jmandel/fhir`,
`txpack-future`) is the integration layer that wires those engines into a runnable, pinned,
judgeable build. Its own contributions are:

- **`tools/build/launch.jar`** — the committed launcher (built from `tools/build/launcher-src/Launcher.java`); the single entry point that resolves the toolchain and content pins and invokes the build.
- **`tools/build/kindling-wrapper.properties`** — the toolchain pin (`toolUrl` + `toolSha256`), resolved into `~/.fhir/tools/<sha256>.jar`; the maintainer-owned writer of the find/trust/keep pair.
- **`tools/build/build.sh`** — the thin front-door alias that just invokes the launcher jar.
- **`fhir.lock`** (repo root) — the content lock (npm `package-lock` v3 profile) pinning the `*.txpack` and the `expectedOutput` signature; the refresh-bot-owned writer (see §7).
- **`tools/build/fhir-settings.json`** — the FHIR settings file the build runs against.
- **`tools/build/ref.manifest`** + **`tools/build/noise-files-v2.txt`** — the judge inputs: the reference output manifest to diff against and the static allowlist of excused-volatility files (see §9).
- **`.github/workflows/txpack-future.yml`**, **`.github/workflows/txpack-record.yml`**, **`.github/workflows/txpack-refresh.yml`** — the build/record/refresh CI pipeline (see REFRESH-AND-CI.md).
- **README / FUTURE doc edits** — the repo-root [README.md](../../README.md) and [FUTURE.md](../../FUTURE.md) front-door narrative.

---

## 9. The determinism contract

Determinism is **layered**:

**Accidental nondeterminism — FIXED at source** (mostly in core, two in kindling):

| Fix | Where | What it cured |
|---|---|---|
| defensive-copy designations before appending synthetic `preferredForLanguage` | `ValueSetExpander` | repeated expansions accumulated duplicates → flickering designation counts |
| `crossLink` idempotency guard | core | repeated linking side effects |
| `HashSet` → `LinkedHashSet` for `allResources` | `CanonicalResourceManager` | `getList()` order varied run-to-run |
| CRC32-derived patient-photo file names | `PatientRenderer` | non-deterministic photo names |
| `volatile` lazy caches + `lazyCacheLock`; serialized `generateSnapshot` under `SNAPSHOT_GEN_LOCK` | `ContextUtilities` | torn reads / half-built lists under the parallel pool |
| order-independent unknown-CodeSystem result shape | `BaseWorkerContext` | shape depended on memo-arming order |
| deterministic ShEx emission (sort ValueSets by `getVersionedUrl`, `TreeSet` iteration) | `ShExGenerator` | reference/known-resource/value-set blocks reordered |
| stable-order Jena `SortingGraph` proxy | `FHIRResourceFactory` (kindling) | `fhir.ttl` triple/bnode ordering varied each build |
| null-safe profile fetch | `ExampleInspector` (kindling) | NPE / inconsistent path on missing profile |
| deterministic NamingSystem registry sort (by `fullUrl`) | `Publisher` (kindling) | `namingsystem-terminologies` bundle order from a `HashMap` keySet |

**Inherent volatility — NORMALIZED:** clock, UUID, section-number and path variance are normalized
where they cannot be removed — in pack diffs via `VOLATILE_PATTERNS`, and in the output judge
(`OutputManifest`) via the order-insensitive `O:` hash and timestamp/uuid normalization.

**Server nondeterminism — DETECTED and DEGRADED:** the recorder's `diff` / `reproduceFilter`
detect changes a re-recording cannot reproduce, poison filtering drops transport flakes, and
carry-forward keeps a known-good answer rather than spuriously removing one.

> The judge is honest by construction: `OutputManifest.EXCLUDED_PATHS` is `\A(?!x)x`, a regex that
> matches **nothing** — a deliberate non-exclusion, because a blanket path exclusion once hid the
> `fhir.ttl` ordering bug. Real volatility is excused only with *evidence* (`-prev`, a same-commit
> twin build) or the static `-allowlist` (`noise-files-v2.txt`), never by a blanket pattern.

---

## 10. Caveats on hosting and stale numbers

- **Hosting is demo scaffolding.** The `resolved` URL in `fhir.lock` and the toolchain `toolUrl`
  currently point at `jmandel/fhir-perf` GitHub Releases. That is demo scaffolding; the intended
  channel is **npm-format FHIR packages** (the `package-lock` v3 profile is exactly so the lock
  needs no new format when that channel exists). Treat the GitHub URLs as a stand-in.
- **Stale references not to propagate.** `Launcher.java`'s javadoc and `docs/txpack-proposal.md`
  still say "tx.lock"; the real file is **`fhir.lock`**. `proposal.md` also carries pre-fix timing
  numbers. The authoritative output signature is `TxLock.expectedSignature()` from the lock
  (`{0, 3694, 349}`), not any comparison string hardcoded elsewhere; the refresh workflow's
  hardcoded `Warnings=3693/Info=345` comment is stale — the gate reads the lock.

---

### See also

- [README.md](README.md) — the 60-second quickstart
- [BUILD-REFERENCE.md](BUILD-REFERENCE.md) — every flag, file, and system property
- [REFRESH-AND-CI.md](REFRESH-AND-CI.md) — the nightly recorder, refresh PRs, and CI gates
- [FUTURE.md](../../FUTURE.md) — intent and the find/trust/keep trust model
- root [README.md](../../README.md) — the legacy build (bypassed by this branch)
