# The refresh pipeline and CI: recording packs, the judge, and the lock-bump PR

This is the operator/reviewer manual for the *upkeep* side of the txpack build: how the
pinned terminology answer-pack stays current, how the determinism judge decides whether a
candidate is real, and how a pack change becomes a one-line `fhir.lock` bump that a human
reviews. If you build the spec but never touch the pack, you want
[tools/build/README.md](README.md) and [tools/build/HOW-IT-WORKS.md](HOW-IT-WORKS.md)
instead. This document is for two seats:

- the **refresh-bot operator** — owns the nightly recorder and the candidate it proposes;
- the **lock-bump reviewer** (Grahame's seat) — the sole human gate, who reads the
  machine-verified diff and merges (or not).

Companion docs: [README.md](../../README.md) ·
[BUILD-REFERENCE.md](BUILD-REFERENCE.md) · [HOW-IT-WORKS.md](HOW-IT-WORKS.md).

---

## 1. Why a recorder at all

The hermetic build replays a **pinned** answer pack — every terminology question
(`validate-code`, `expand`, server capabilities, `findTxResource`, tx-registry `/resolve`)
is answered from immutable, content-addressed bytes, and any network attempt is a hard
failure. That is great for speed and reproducibility, but a frozen pack drifts from reality:
the spec grows new codes the pack never recorded, and the canonical server *fixes* answers
the pack still holds at their old value.

The recorder closes that gap. It re-asks the build's full terminology question set against a
**live** server, captures what the server says **today**, and — only if the answers actually
changed and a second run confirms it — proposes a new pack. Two design rules frame everything
below:

- **The bot is the sole `fhir.lock` writer.** The maintainer owns the toolchain pin
  (`kindling-wrapper.properties`); the bot owns the content pin (`fhir.lock`). They never
  collide. No human edits `fhir.lock` by hand.
- **The server is assumed flaky.** A transient failure must never look like a deliberate
  removal, and a load-induced flap must never reach a PR. Carry-forward and
  reproduce-before-propose (below) are the two mechanisms that enforce this.

---

## 2. `record`: a COLD recording against the live server

`SpecBuild record` runs the build **cold** — *not* seeded from the pinned pack — online, with
recording on. Seeding would serve every already-known answer locally and only send genuinely
*new* questions to the server, so a server **fix** to an existing answer (the whole point of a
refresh) would be invisible. Running cold re-asks everything, so the fresh recording reflects
what the server says now; diffing it against the pinned pack catches changes, additions, and
(carried-forward, never auto-proposed) removals.

The recorder pins this configuration in-process before delegating to `Publisher.main`
(verified in `SpecBuild.record()`):

| Setting | Value | Why |
|---|---|---|
| `~/.fhir/tx-cache` | **deleted** | cold start: the build re-asks every question |
| `org.hl7.fhir.tx.lock` | `ignore` | never seed a pack during a refresh recording |
| `org.hl7.fhir.tx.pack` | **cleared** | same — nothing pre-answers the server |
| `org.hl7.fhir.tx.recordSemanticErrors` | `true` | persist semantic `CODESYSTEM_UNSUPPORTED` answers the default policy drops, and shadow-send suppressed unknown-system requests so the pack is exhaustive regardless of thread timing |
| `org.hl7.fhir.tx.localFirst` | `true` | **match the pinned build/pack**: grammar/unknown-system answers are synthesized in the same shape the pack stores, so a no-drift night diffs to exactly zero (mismatched shapes report phantom drift) |
| `fhir.build.validation.threads` | `1` | serial validation dodges the known parallel search-param flake — fidelity over speed here |
| locale / timezone | `en-US` / `America/Chicago` | both leak into request keys; must reproduce the pack's recording environment |

It needs a **live terminology server** — point it at one with `-fhir-settings`:

```bash
java -Xmx11g -cp tool.jar org.hl7.fhir.tools.publisher.SpecBuild record . \
  -fhir-settings /tmp/record-settings.json -out candidate-A
```

(The demo pack was recorded against a local FHIRsmith; the committed
[`tools/build/fhir-settings.json`](fhir-settings.json) and the nightly both point at
`tx.fhir.org` — see [Honest notes](#6-honest-notes). The loop is identical either way.)

**Carry-forward merge `[pinned-first, fresh]`.** The candidate is
`merge([pinned, fresh])` with the **pinned pack listed first**, so the fresh recording
*supersedes* it on every key the build re-asked — a server fix is therefore never masked
(that is the distinction from seeding, which would have skipped the ask). The pinned answer
survives **only** where the fresh recording is absent for it — i.e. a transient server
failure this run — so **a flake can never become a spurious "removed."** A hard "every request
must succeed" gate would never complete against a flaky server; carry-forward degrades
gracefully instead.

The candidate is then diffed against the pinned pack:

- **identical** → exit `0` silently (the common nightly outcome, one build, no proposal);
- **changed** → emit two artifacts:
  - **`candidate-A.zip`** — the carry-forward candidate (what would be proposed), and
  - **`candidate-A.fresh.zip`** — a *sidecar* of the **raw** fresh recording (no carry-forward),
    which the reproduce step consumes.

> Removals are carried forward, never auto-proposed. A genuine code retirement is surfaced by
> a separate staleness signal, not invented from one night's missing answer.

---

## 3. `reproduce-before-propose`: a second, independent recording

A single recording can still capture a load-induced flap. So before anything is proposed, a
**second independent recording runs in a SEPARATE JVM** (the Publisher is not re-entrant
in-process — you cannot run two builds in one JVM) and produces its own
`candidate-B.fresh.zip`. Then:

```bash
java -cp tool.jar org.hl7.fhir.tools.publisher.SpecBuild reproduce \
  -fresh candidate-A.fresh.zip -confirm candidate-B.fresh.zip \
  -pinned fhir.lock -out candidate-pack
```

`reproduce` runs `reproduceFilter(A, B, pinned)`: it keeps a delta-vs-pinned **only if both
recordings confirm it**, carries the pinned answer forward for everything else, and proposes
`candidate-pack.zip` only if a confirmed change survives. A delta that appears in A but not B
is dropped as server flakiness. As in `record`, **removals are never auto-proposed** — only
confirmed *additions and changes* survive. Output:

- no surviving change → `reproduce: no change survived confirmation (... dropped ...)`, exit `0`;
- a surviving change → `candidate-pack.zip` + its sha256, ready for the downstream PR workflow.

So a proposal is, by construction, a real change in the canonical server's answers — never our
own noise, never a single-run hiccup.

---

## 4. The judge: `OutputManifest`

The recorder decides whether the **pack** changed. The judge (`OutputManifest`, driven by
`SpecBuild compare` / `--judge`) decides whether the **published spec output** changed in a way
that matters. It fingerprints every file in `publish/` into a manifest (one
`HASH<TAB>relpath` line per file, sorted) and compares.

**Hash kinds** (from the manifest line prefix):

| Kind | Applies to | What it hashes |
|---|---|---|
| `N:` | text | normalized-exact: bytes after `TS_PATTERNS` neutralize the build clock, UUIDs, section numbers, usernames |
| `O:` | text | order-insensitive "second chance": the sorted multiset of trimmed non-empty lines |
| `X:` | `.xlsx` | recursive member hash with POI font-metric column-width variance (`width="…"`) removed |
| `A:` | archives (`.zip .jar .pack .tgz`) | recursive normalized member signature; build-date normalized per member, **entry mtimes ignored** |
| `B:` | binaries (`.png .gif .pdf .woff …`) | raw content hash |

**The `O:` second chance.** Text files carry both an `N:` and an `O:` hash.
`contentDifferent()` flags a file only when the `N:` parts differ **and** the `O:` parts also
differ — i.e. a file whose exact bytes moved but whose set of lines is unchanged (pure element
**reordering**, the documented stock nondeterminism class) compares **equal** and is excused.
Any real content change still flips the `O:` hash and flags.

**`EXCLUDED_PATHS` matches NOTHING by design.** It is literally
`Pattern.compile("\\A(?!x)x")` — a regex that can never match. The judge content-checks
**every** published file and excludes nothing by path. This is deliberate: a blanket path
exclusion is an excuse that can hide a real regression — one previously hid `fhir.ttl`'s
ordering bug. Everything that was once path-excused (`fhir.ttl`, all `.shex`/`.shex.html`,
`definitions.*.zip`, `validator.pack`) is now genuinely deterministic at the source and stays
in the gate.

**Two excusal channels** (both opt-in, neither hides a content change):

- **`-prev <evidence>`** — a manifest from a **second same-commit build**. A file whose hash
  provably differed between the two twin builds is excused *as observed nondeterminism this
  run* — evidence, not assertion.
- **`-allowlist <file>`** — static noise excusal. The committed list is
  [`tools/build/noise-files-v2.txt`](noise-files-v2.txt) (246 plain paths).

**The committed reference is [`tools/build/ref.manifest`](ref.manifest)** — the manifest the
judge compares a fresh build against. `SpecBuild build --judge` runs:

```bash
compare tools/build/ref.manifest build-future.manifest \
  -allowlist tools/build/noise-files-v2.txt [-prev <previous build's manifest>]
```

and returns nonzero on any unexplained content difference.

**The signature gate is separate from the judge.** After every `SpecBuild build`, the tool
scrapes the build log's `Summary: Errors=E, Warnings=W, Information messages=I` line and
compares it to `TxLock.expectedSignature()` — the `{errors, warnings, information}` pinned in
`fhir.lock` under `expectedOutput` (currently `0 / 3694 / 349`). A mismatch returns `1` unless
`SIGNATURE_GATE=report` (env) or `-Dorg.hl7.fhir.spec.signatureGate=report` downgrades it to a
report. This is a cheap whole-build aggregate check; the manifest judge is the byte-level one.

---

## 5. CI workflow walkthrough

Three workflows live on `txpack-future`. One gates every push; two run the refresh loop.

### `txpack-future.yml` — the build gates

| Job | Trigger | What it does |
|---|---|---|
| **`future-hermetic`** | **every push** (required) | one hermetic cold build, `HEAP=11g ./tools/build/build.sh` — zero terminology network (enforced) + signature check. The everyday cost. |
| **`reproducibility`** | `[parity]` commit msg or dispatch (`continue-on-error`) | two builds: `build.sh --manifest` (convergence + evidence pass) then `build.sh --judge` (compare vs `ref.manifest`, with `-allowlist` and `-prev`). Non-blocking signal. |
| **`determinism-gate`** | `[determinism]` commit msg, or workflow_dispatch with the `parity` input (it shares reproducibility's dispatch checkbox — there is no separate determinism input) | three builds (first discarded for convergence) then `SpecBuild compare mA.manifest mB.manifest` with **no `-allowlist` and no `-prev`** — any non-ordering content diff between two same-commit builds **fails**. This is the "accidental nondeterminism must be empty" gate; only `O:`-equal reordering is excused. |
| **`stock-baseline`** | dispatch only | the legacy `./publish.sh` cold build against `tx.fhir.org`, for timing comparison. Slow; run rarely. |

The `reproducibility` job *allows* expected differences (it judges against the committed
reference with the noise allowlist and evidence excusal); the `determinism-gate` *forbids*
them (no allowlist, no evidence — the build must be a pure function of its inputs). They are
deliberately different strictnesses.

### `txpack-record.yml` — the nightly recorder (upstream half)

Cron `11 9 * * *` (09:11 UTC, off-peak for `tx.fhir.org`), plus dispatch with a `server`
input. Steps:

1. fetch the pinned tool jar from `kindling-wrapper.properties` (sha-verified);
2. **Record (build A)** — `SpecBuild record . -fhir-settings … -out candidate-A` (cold,
   carry-forward). A `candidate-A.fresh.zip` sidecar appears **only** on a real delta;
3. **Reproduce** — if (and only if) the sidecar exists, a second cold recording
   (`-out candidate-B`, separate JVM) then `SpecBuild reproduce -fresh candidate-A.fresh.zip
   -confirm candidate-B.fresh.zip -pinned fhir.lock -out candidate-pack`;
4. **Propose** — if `candidate-pack.zip` survived, name it `txpack-<sha[:8]>.zip`. In
   production this uploads the pack to the immutable, hash-named pack store and writes
   `.txpack/refresh-request.json` (`candidate_url`, `candidate_sha256`), commits, and pushes —
   which triggers the refresh workflow.

Most nights: one cold build, no delta, silent. The second build only runs when build A already
saw a real change.

### `txpack-refresh.yml` — the lock-bump PR (downstream half)

Triggered by a push touching `.txpack/refresh-request.json`. It produces evidence, then opens
a PR that mutates **only `fhir.lock`** (branch `txpack-bump-<date>`). The pack itself never
lands in git — it lives, immutable and hash-named, in the pack store; the review gate is the
one-line lock change. Steps:

1. **Read request + fetch packs** — pull `current.zip` (resolved from `fhir.lock`,
   sha-verified against its `integrity`) and `candidate.zip` (sha-verified against the request);
2. **Canonical compare** — `SpecBuild diff-packs current.zip candidate.zip`. **Exit `0` =
   nothing changed** (stop); **exit `3` = a real change** (continue). This is the *rc gate*;
3. **Output A/B impact** — build hermetically with `-Dorg.hl7.fhir.tx.pack=current.zip`
   `--manifest`, then again with the candidate pack, then `SpecBuild impact` the two manifests.
   A candidate that **drops** a needed answer fails loudly here (hermetic both times). This
   reports exactly which published files change — or "output-inert";
4. **Explanation (GitHub Models)** — an LLM writes a short prose summary **from the diffs only**.
   It is narrative, never authority: the deterministic Java diff/impact is the source of truth,
   and the prose is allowed to fall back to "service unavailable" without blocking;
5. **Open the PR** — rewrite `fhir.lock`'s `resolved` + `integrity` to the candidate, commit,
   push, and `gh pr create` with the machine diff + impact + AI narrative in the body.

**Merge policy** is a dial on the *machine facts*, never the prose:

- **additions-only** candidates → safe to **auto-merge** (new answers, nothing changed);
- **changed answers** → **human review** (an existing answer moved — the reviewer reads the
  diff and the published-output impact and decides).

---

## 6. Honest notes

- **The "3693 / 345" string in `txpack-refresh.yml` is stale.** The Output A/B step prints a
  hardcoded comparison line — `current pack: Errors=0, Warnings=3693, Information messages=345`
  — that is **cosmetic and out of date**. The **real** signature gate uses
  `TxLock.expectedSignature()`, which reads `fhir.lock`'s `expectedOutput`
  (currently **`0 / 3694 / 349`**). Trust the lock and the gate, not that printed string.
- **The Publisher is not re-entrant.** Two builds cannot share one JVM, which is why
  reproduce-before-propose and the Output A/B impact both run as *separate processes*. This is
  a structural constraint, not a tuning choice — do not "optimize" it into one JVM.
- **The demo pack was recorded against a local FHIRsmith,** which the public nightly cannot
  reach. A real rollout would re-base the pack on `tx.fhir.org` for consistent provenance
  (`recordedAgainst` in `fhir.lock`); the recording/reproduce/refresh loop is identical
  either way.
- **`txpack-refresh.yml` is a prototype.** It takes the candidate as a request-file input so
  the downstream machinery runs in CI without a live server. The `push`-path trigger is a
  demo-branch workaround (the dispatch API requires the workflow to exist on the default
  branch); on a real deployment it lives on `main` and uses plain `workflow_dispatch`.
