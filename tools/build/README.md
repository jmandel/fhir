# The txpack FHIR spec build

Building the FHIR spec means answering thousands of terminology questions — *is this code
valid? what is in this value set?* — that the build would otherwise put to a server over the
network, one round-trip at a time.[\*](#faq) txpack answers them once, records the answers in a
hash-pinned file, and replays that file on every build. Cold builds drop from 18+ minutes to
~4, run fully offline, and produce identical bytes on every machine. The build touches the
server only when you add a genuinely new code.

New here? Read the [Quickstart](#quickstart); the rest is reference.

---

## The problem

As it runs, the build keeps turning to a **terminology server** (e.g. `tx.fhir.org`): *Is this
code valid? What does this value set contain? Does this system exist?* — thousands of times,
each a network round-trip. Three costs follow:

- **Slow.** A cold build spends most of 18+ minutes waiting on the network.
- **Flaky.** A timeout, a shed request, or an answer that shifts mid-build can derail the run or
  quietly change its output.
- **Unreproducible.** The server's answers drift, so the same source builds differently on
  different days and machines.

A per-machine cache softens this for repeat builds on one machine — but it is mutable, unshared,
and unreviewed, with failure modes of its own (see the [FAQ](#faq)).

## The idea: record the answers, replay them

Your package manager already does this. `npm install` reads `package-lock.json` to pin exactly
which bytes you get instead of re-resolving the graph; `gradlew` fetches and verifies a pinned
Gradle instead of trusting whatever is installed. txpack applies both moves to terminology:

1. **Record once.** Run online against the live server with recording on; capture every question
   and its answer into a content-addressed zip — the **pack**.
2. **Pin it.** `fhir.lock` names that pack by URL and hash, exactly like an npm lockfile entry.
3. **Replay offline.** Every later build reads from the pack and never touches the network — a
   **hermetic** build (sealed: the output depends only on your inputs).

### The three pinned files

Each has exactly one writer, so they never collide:

| File | Analogy | Pins | Writer |
|---|---|---|---|
| `tools/build/launch.jar` | `gradlew` | the bootstrapper you run | maintainer (recompiles `launcher-src/Launcher.java`) |
| `tools/build/kindling-wrapper.properties` | `gradle-wrapper.properties` | the build **tool** (URL + sha256) | maintainer, on tool releases |
| `fhir.lock` | `package-lock.json` | the terminology **answers** (the pack) | the refresh **bot**, never a human |

Maintainers own the toolchain, the bot owns the content, and as an everyday editor you touch none
of them.

---

## Quickstart

You need a JDK (17+) and nothing else — no bash, Python, ant, curl, or git. Run from the spec
root (the folder holding `fhir.lock` and `tools/`), and give the JVM ~12 GB of heap.

```sh
java -jar tools/build/launch.jar build .
```

On a unix shell, a four-line alias `cd`s to the root for you:

```sh
./tools/build/build.sh
```

Windows runs the same launcher:

```powershell
java -jar tools\build\launch.jar build .
```

The build is hermetic by default: every answer comes from the pack, and the network is never
touched. It says so at startup:

```
hermetic: any terminology network attempt is a hard failure (use --online when adding new codes)
```

### Adding new codes

When your edit needs terminology the pack has never seen, the hermetic build stops and names what
it needs — it never guesses or skips. Re-run online to send those questions, and only those, to
the server:

```sh
./tools/build/build.sh --online
```

`--online` hides nothing. It answers what it can from the pack, asks the server the rest, and
reports the count:

```
terminology questions not answered by the pack: N (these will fold into the pack at the next refresh)
```

You never edit the pack or `fhir.lock` to absorb those answers; the nightly
[refresh loop](#keeping-the-pack-correct-the-refresh-loop) does that.

### Common variations

| Goal | Command |
|---|---|
| Build offline from the pin (default) | `./tools/build/build.sh` |
| Build after adding new codes | `./tools/build/build.sh --online` |
| Check output against the reference | `./tools/build/build.sh --judge` |
| List which published files your build changed | `./tools/build/build.sh --impact` |
| Use less heap | `HEAP=11g ./tools/build/build.sh` |
| Pass extra JVM options | `JAVA_OPTS="-XX:+UseParallelGC" ./tools/build/build.sh` |

Heap and JVM options come from the environment (`HEAP`, `JAVA_OPTS`), because the launcher applies
them when it spawns the build's JVM.

The build writes the site to `publish/` and its log to `build-future.log`. SpecBuild's own summary
lines (`hermetic:`, `BUILD ok…`, `signature:…`) print to the console after the log closes.

> The repo-**root** `build.sh` is the old bash + git-svn + ant build, not this one. The front door
> is `tools/build/build.sh` / `launch.jar`.

---

## Concepts

- **Pack** — the content-addressed zip of recorded answers: validate-code, expand, server
  capabilities, resource lookups, registry resolutions. It is a terminology-cache directory,
  zipped. It records negative answers ("not on the server") too, so a recorded "no" heads off a
  pointless re-ask.

- **Recording** — a cold, online build (`SpecBuild record`) that clears the local cache, ignores
  the pinned pack, and re-asks the build's full question set against the live server. Going cold
  is what exposes a *changed* server answer; the result is a candidate pack.

- **Reproduce** — the flake defense. A second independent recording runs in a fresh process, and a
  difference survives only if both recordings agree. A one-time blip is discarded as server noise,
  so every proposed change is real by construction.

- **Refresh** — the nightly loop that keeps the pin honest: record → reproduce → open a one-line
  `fhir.lock` PR, but only on a confirmed change. A frozen pack drifts as the spec grows and the
  server corrects answers; refresh closes that gap.

- **The judge** — *does a pack-built spec still produce the right output?* It fingerprints every
  published file and compares against a known-good reference. The work is separating a real change
  from noise: timestamps, generated UUIDs, section numbers, reordered lines. So it normalizes that
  noise and treats pure reordering as equal. It never excuses a file by name — a blanket by-name
  skip once hid a real bug — so a file counts as "expected to vary" only on evidence: a same-commit
  twin build (`-prev`) or the explicit allowlist `tools/build/noise-files-v2.txt`. `build --judge`
  runs it against `tools/build/ref.manifest`.

- **Pin bootstrap** — recording the first pack when no prior pin exists to carry forward from
  (`SpecBuild record -bootstrap`): record cold, skip the carry-forward merge, emit a fresh pin.
  Distinct from the launcher's *tool* bootstrap, which fetches and verifies the tool jar.

---

## Keeping the pack correct: the refresh loop

A pinned pack is a snapshot, and snapshots go stale: the spec gains codes, and the server
occasionally corrects an answer. The refresh loop keeps the pin matching reality without anyone
hand-editing it. Each night a job:

1. **Records** a fresh pack cold against the live server, carrying a pinned answer forward wherever
   the fresh run only dropped it to a transient failure — so a flake never becomes a phantom removal.
2. **Reproduces** every difference with a second recording, keeping only what both runs confirm.
3. **Proposes** a confirmed change by publishing the candidate pack and opening a PR that rewrites
   *only* `fhir.lock`.

Most nights nothing changed: one build, no diff, silence.

Two writers keep this safe:

- **You, the editor,** change content and never touch `fhir.lock`. If an edit needs new codes, you
  build `--online`, note the count, and move on; the bot folds them in.
- **The bot** is the only writer of `fhir.lock`, and backs every bump with a reproduced diff and an
  output-impact report. Maintainers, separately, are the only writers of the tool pin.

---

## How you know it's right

Four checks, each guarding a different failure:

| Check | Question it answers | When |
|---|---|---|
| **Hermetic 0-miss** | Is the pack complete — can the build run fully offline? | on pin/tool changes (and every push, as a signal) |
| **The judge** | Does the published output still match the reference? | on demand (`--judge`, `[parity]` CI) |
| **Determinism** | Do two builds of one commit produce identical bytes? | on demand (`[determinism]` CI) |
| **Drift diff** | Does the pin still match what the server says today? | nightly |

Hermetic 0-miss is strong because a missing answer is a hard failure that names the request, not a
silent fallback — so zero misses proves completeness. Determinism is what makes the judge
meaningful: if one commit built to different bytes each run, comparing against a reference would
compare noise (run-to-run variance fell from ~21 differing files to ~0). The **signature gate** is
the judge's always-on cousin: after every build it compares the error/warning/info summary against
the count pinned in `fhir.lock`, and fails on a mismatch unless you set `SIGNATURE_GATE=report`
(or `-Dorg.hl7.fhir.spec.signatureGate=report`) while a change intentionally moves it.

---

## CI workflows

Three workflows run on `txpack-future`. The recurring design choice is **signal vs. gate**: a
content push reports new terminology rather than blocking on it, and strict offline-completeness is
enforced only when the pin or tooling changes.

- **`txpack-future.yml` — build gates.** *future-build* (every push) runs one pack-seeded
  `--online` build: a no-new-code push touches the network zero times, new codes surface as an
  "N new terminology questions" signal, and the signature is reported, not gated.
  *pinned-world-hermetic* (only when `fhir.lock` or the wrapper changes) runs a strict hermetic
  build with zero-network and the signature both enforced. *reproducibility* (`[parity]`) runs the
  judge; *determinism-gate* (`[determinism]`) demands two byte-identical builds of one commit;
  *stock-baseline* (manual) times the legacy build.

- **`txpack-record.yml` — the nightly recorder.** A cron job records → reproduces → publishes a
  candidate pack to this repo's own `txpack-store` release, then emits a `refresh-request` artifact.
  A manual `full_pin` job does a from-scratch [pin bootstrap](#concepts), hermetic-verified.

- **`txpack-refresh.yml` — the lock-bump PR.** Chained off the recorder by `workflow_run`, so it
  needs no cross-repo token: it diffs the candidate against the pin, builds an A/B output-impact
  report, has an LLM summarize those diffs in plain English, and opens a PR that touches **only
  `fhir.lock`**.

Reads need no token — they are public `curl`, verified on use. The **pack** publishes to this
repo's `txpack-store` release (default `GITHUB_TOKEN`); the **tool jar** lives on the
`jmandel/fhir-perf` release (human-bumped).

---

## What changed in each fork (for the curious)

Three forks cooperate, all on `txpack-future`:

- **core (`org.hl7.fhir.core`) — the terminology engine.** A read-only pack seed layer sits in
  front of the live cache (a miss falls through — it seeds, it does not wall); a hermetic mode
  whose violation is an `Error`, not an `Exception`, so tx clients cannot swallow it; cache-key
  canonicalization so per-run labels do not poison replay; the pack build/merge/diff tooling with
  flake filtering; `TxLock`, which reads `fhir.lock` and verifies integrity on every use; and a
  batch of determinism fixes.
- **kindling — the build front end.** The `SpecBuild` CLI
  (`build`/`manifest`/`compare`/`impact`/`diff-packs`/`record`/`reproduce`) that replaces shell
  scripting; the cold recorder; terminology routed through the cache so it replays; the honest
  judge; and its own determinism fixes.
- **this repo (`jmandel/fhir`) — integration.** The three pinned files, the unix alias, the FHIR
  settings, the judge's reference and allowlist, and the three CI workflows.

Run `./tools/build/build.sh help` for the flag-by-flag CLI reference.

---

## The determinism contract (in brief)

Determinism makes every output comparison sound. It is enforced in three layers: *accidental*
nondeterminism (shared state, hash-set ordering, torn parallel reads) is fixed at the source;
*inherent* volatility (clocks, UUIDs, section numbers) is normalized in diffs and in the judge;
*server* nondeterminism is detected and degraded by flake filtering, carry-forward, and
reproduce-before-propose. For the full treatment, see
[docs/txpack-vision.md](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md).

---

## FAQ

### \* Didn't the build already cache terminology answers — and download a zip from tx.fhir.org?

It did — and this is the core spec build, not an IG-Publisher-only path; both drive the same
`TerminologyCacheManager` in core. The stock toolchain keeps a per-machine cache at
`~/.fhir/tx-cache/{org}/{repo}/{branch}/` (the `{org}/{repo}/{branch}` comes from the checkout's
git branch), read at startup and written each run, so a *warm* build already finishes in ~197s
while a *cold* one takes ~18 min. It also moves that cache over the network: at startup it downloads
`https://tx.fhir.org/tx-cache/.../{branch}.zip` when its version stamp is stale, and a build that
holds a tx.fhir.org API key (in practice HL7's CI or the tx maintainer) zips the cache back up at
the end and PUTs it over the same URL in place. (That cache is separate from the IG *expansions
package*, which ships pre-expanded value sets.)

So txpack does not remove network calls the cache had already removed for warm builds. It changes
what the cache *is*. The stock cache is mutable, machine-local, ungated, overwritten in place, and
has no committed reference at all — a build just trusts whatever happens to be on disk or in the
shared zip. The pack is immutable and named by a hash that a reviewed `fhir.lock` pins. As with
`package-lock.json`, what's committed is the *pin*, not the data: the pack bytes live in a release,
fetched and integrity-checked on use (like `node_modules`, never vendored into git). A
machine-local convenience becomes a shared, reviewed contract.

### Then what does the pack fix?

The stock cache's documented failure modes — each removed by the immutable, pinned model:

| Stock cache failure | Pinned pack |
|---|---|
| A server flake writes an *error* into a `.cache` file, and every later warm build replays it as a failure (the folk remedy is `rm -rf ~/.fhir/tx-cache`). | A pack cannot hold a transport error: a recording stores a clean answer or nothing. |
| It is meant to be per-branch, but it picks the branch wrong — the *alphabetically-first* local branch, not the one checked out — so builds of different branches collide in one directory. | One content-addressed pin, the same on every branch and machine. |
| Any successful build that holds a tx.fhir.org API key overwrites the one shared zip — no hash, diff, review, or provenance. | A refresh proposes an immutable, hash-named pack through a reviewed PR with a machine diff. |
| Forks and CI start cold (their zips 404), and a cache-version bump re-colds everyone at once. | The answers arrive with the checkout, fetched by hash: cold equals warm for forks, CI, and after a bump. |
| A failed build discards its fetches, so the fail→fix→rerun loop re-pays the full network bill. | The build is hermetic and complete: every answer is present or it stops, and reruns are free and offline. |

The payoff is not "faster than warm." It is **cold equals warm, identical everywhere,
un-poisonable, and drift-free** — measured at ~195s cold-with-pack against ~197s warm, and ~231s
fully hermetic with the exact reference signature.

### Isn't this the expansions package?

Same instinct — ship terminology as data instead of asking for it — but the expansions package
ships only pre-expanded value sets. The pack covers the build's whole question surface:
validate-code, expand, server capabilities, resource lookups, registry resolutions, and negative
answers, all pinned by hash.
