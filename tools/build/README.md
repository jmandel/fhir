# The txpack FHIR spec build

**Building the FHIR spec used to mean asking a terminology server thousands of
questions over the network — slow, flaky, and never quite the same twice. So we
ask those questions *once*, write down all the answers in a file we pin like a
lockfile, and let every build replay from that file instead of the network.**
That pinned file of answers is the **pack**. The result: ~4 minutes instead of
18+, fully offline, byte-identical on every machine. The server gets touched
only when you genuinely add new codes.

If you read nothing else, read the [Quickstart](#quickstart). Everything after
it is "go deeper."

---

## The problem

A normal FHIR spec build doesn't just transform files. As it runs, it constantly
turns to a **terminology server** (like `tx.fhir.org`) and asks: *Is this code
valid? What's in this value set? Does this system exist?* — thousands of such
questions, each a network round-trip. That has three costs:

- **Slow.** A cold build is 18+ minutes, most of it waiting on the network.
- **Flaky.** A server hiccup, a timeout, a changed answer mid-build — any of these
  can derail or subtly alter the result.
- **Non-reproducible.** Two people, or two days, can get different output because
  the server's answers drifted underneath them.

## The idea: record the answers, then replay them

Here's the trick, and it's the same one your package manager already uses.

When you run `npm install`, npm doesn't re-resolve the whole dependency graph from
scratch — it reads `package-lock.json`, which pins *exactly* which bytes you get.
When you run a Gradle build, you don't install Gradle yourself — a tiny committed
`gradlew` wrapper fetches the exact pinned Gradle for you and verifies it. We do
both of those things, for terminology:

1. **Record once.** Run the build online against the live server with recording
   turned on. Capture every question and its answer into a content-addressed zip:
   the **pack**.
2. **Pin it like a lockfile.** A file called `fhir.lock` names that pack by URL and
   hash — exactly the shape of an npm lockfile entry.
3. **Replay offline.** Every subsequent build reads answers from the pack and never
   touches the network. This is called a **hermetic** build (hermetic = sealed: no
   outside influence, so the result depends only on your inputs).

### The three pinned files

Just like Gradle has *a wrapper script* + *a pinned distribution*, and npm has
*a lockfile*, this system has three committed files — and each has exactly one
writer, so they never step on each other:

| File | Analogy | Pins | Who edits it |
|---|---|---|---|
| `tools/build/launch.jar` | `gradlew` | the bootstrapper you run | a maintainer (recompiles `launcher-src/Launcher.java`) |
| `tools/build/kindling-wrapper.properties` | `gradle-wrapper.properties` | the build **tool** (by URL + sha256) | a maintainer, on tool-release cadence |
| `fhir.lock` | `package-lock.json` | the terminology **answers** (the pack) | the refresh **bot**, never a human |

The split is the whole point: maintainers own the *toolchain*, the bot owns the
*content*, and as an everyday editor **you change none of them**.

---

## Quickstart

**Prerequisites:** a JDK (Java 17+) and nothing else — no bash, Python, ant, curl,
or git. Give the JVM ~12 GB of heap. Run from the **spec checkout root** (the folder
that contains `fhir.lock` and `tools/`).

### The one command

```sh
java -jar tools/build/launch.jar build .
```

That's it, on any OS. On a unix shell there's a 4-line convenience alias that
`cd`s to the root for you:

```sh
./tools/build/build.sh
```

On Windows, run the exact same launcher (it's pure JDK):

```powershell
java -jar tools\build\launch.jar build .
```

By **default the build is hermetic**: it serves every answer from the pack and
touches the network zero times. At startup it prints:

```
hermetic: any terminology network attempt is a hard failure (use --online when adding new codes)
```

### When you add new codes

If your edit introduces terminology the pack has never seen, the hermetic build
**fails loudly and names what it needs** — it does not silently guess or skip.
That's your cue to re-run online so the new questions (and only those) go to the
server:

```sh
./tools/build/build.sh --online
```

`--online` doesn't hide anything; it answers everything it can from the pack, asks
the server only the genuinely-new questions, and tells you how many:

```
terminology questions not answered by the pack: N (these will fold into the pack at the next refresh)
```

You never edit the pack or `fhir.lock` to absorb those answers — the nightly
[refresh loop](#keeping-the-pack-correct-the-refresh-loop) does that for you.

### Common variations

| You want to… | Command |
|---|---|
| Build offline from the pin (default) | `./tools/build/build.sh` |
| Build after adding new codes | `./tools/build/build.sh --online` |
| Build and check output vs the reference | `./tools/build/build.sh --judge` |
| See which published files your build changed | `./tools/build/build.sh --impact` |
| Use less heap | `HEAP=11g ./tools/build/build.sh` |
| Pass extra JVM options | `JAVA_OPTS="-XX:+UseParallelGC" ./tools/build/build.sh` |

Heap and JVM tuning come from the **environment** (`HEAP`, `JAVA_OPTS`), not from
`build` flags, because the launcher applies them when it spawns the child JVM.

**Output** lands in `publish/` (the generated site) and `build-future.log` (the
Publisher log). SpecBuild's own summary lines (`hermetic:`, `BUILD ok…`,
`signature: …`) print to the console after the log closes, so they're on screen
but not in the file.

> The repo-**root** `build.sh` is a different, legacy thing (the old
> bash + git-svn + ant build). The front door is `tools/build/build.sh` / `launch.jar`.

---

## Concepts

A few terms recur. Each is defined here the first time you'll need it.

- **Pack** — the content-addressed zip of recorded terminology answers (validate,
  expand, server capabilities, resource lookups, registry resolutions). It *is* a
  terminology-cache directory, zipped. Negative answers ("not on the server") count
  too, so a recorded "no" stops a pointless re-ask.

- **Recording** — a *cold, online* build (`SpecBuild record`) that deletes the local
  cache, ignores the pinned pack, and re-asks the build's **full** question set against
  the live server, capturing today's answers. Going cold (not seeded from the old pack)
  is what makes a *changed* server answer visible. The output is a candidate pack.

- **Reproduce** — the flake defense. Before any proposed change is trusted, a **second
  independent recording** runs in a separate process, and a difference is kept **only
  if both recordings agree**. A difference that shows up once and not the other time is
  discarded as server noise. So a proposal is, by construction, a real change.

- **Refresh** — the nightly upkeep loop that keeps the pinned pack honest: record →
  reproduce → (only on a confirmed change) open a one-line `fhir.lock` bump PR. A frozen
  pack drifts as the spec grows and the server corrects answers; refresh closes the gap.

- **The judge** — *How do we know a pack-built spec still produces the right output?*
  We fingerprint every published file and compare it to a known-good reference. The hard
  part is telling a *real* change from meaningless noise — timestamps, generated UUIDs,
  section numbers, reordered lines. So we **normalize the known noise** (and treat pure
  reordering as equal), and we **deliberately never skip a file just because of its name**
  — a blanket by-name skip once hid a real bug. A file is excused as "expected to vary"
  only with *evidence* (a same-commit twin build via `-prev`) or via the explicit
  allowlist `tools/build/noise-files-v2.txt` — never by guessing from a pattern.
  Mechanically: `SpecBuild compare`, run for you by `build --judge` against
  `tools/build/ref.manifest`.

- **Pin bootstrap** — creating the *very first* pack when there's no prior pin to carry
  forward from (`SpecBuild record -bootstrap`). It records cold, skips the carry-forward
  merge, and emits a fresh pin. Don't confuse this with the launcher's **tool**
  bootstrap (downloading and verifying the tool jar) — different mechanism, similar word.

---

## Keeping the pack correct: the refresh loop

A pinned pack is a snapshot, and snapshots go stale: the spec gains new codes, and the
server occasionally corrects an old answer. The **refresh loop** keeps the pin matching
reality, automatically, so humans never hand-edit it.

Every night a job:

1. **Records** a fresh pack cold against the live server, carrying forward the pinned
   answers where the fresh run merely dropped one to a transient failure (a flake never
   becomes a phantom "removed").
2. **Reproduces** any difference with a second independent recording — keeping only what
   both runs confirm.
3. **Proposes** — only if a real, confirmed change survives — by publishing the candidate
   pack and opening a PR that rewrites *only* `fhir.lock`.

Most nights, nothing changed: one build, no diff, silent.

**Two personas, and that's the safety:**

- **You, the editor.** You change content. You never touch `fhir.lock`. If your edit needs
  new codes, you build `--online`, see the count, and move on — the bot folds them in later.
- **The bot.** It is the *only* writer of `fhir.lock`, and every bump it proposes is backed
  by a reproduced diff and an output-impact report. Maintainers, separately, are the only
  writers of the tool pin.

---

## How you know it's right

Four independent checks, each guarding a different failure mode:

| Check | The question it answers | When it runs |
|---|---|---|
| **Hermetic 0-miss** | Is the pack *complete* — can the build run fully offline? | on pin / tool changes (and every push, as a signal) |
| **The judge** | Does the published output still match the known-good reference? | on demand (`--judge`, `[parity]` CI) |
| **Determinism** | Do two builds of the same commit produce identical bytes? | on demand (`[determinism]` CI) |
| **Drift diff** | Does the pinned pack still match what the live server says today? | nightly |

A few notes that tie them together:

- **Hermetic 0-miss** is strong because a missing answer is a *hard failure that names the
  request*, not a silent fallback. Zero misses therefore proves completeness.
- **Determinism is what makes the judge meaningful.** If the same commit produced different
  bytes run-to-run, comparing output to a reference would be comparing noise. Run-to-run
  variance was driven from ~21 differing files to ~0 (see the determinism contract below).
- **The signature gate** is a lightweight always-on cousin of the judge: after every build,
  the error/warning/info summary is compared to the count pinned in `fhir.lock`. A mismatch
  fails the build, unless you set it to report-only with `SIGNATURE_GATE=report` (or
  `-Dorg.hl7.fhir.spec.signatureGate=report`) — useful while a change is intentionally
  moving the signature.

---

## CI workflows (the gate split)

Three workflows live on `txpack-future`. The key design choice is **what's a signal vs. what's
a hard gate**: everyday content pushes report new terminology rather than blocking on it;
strict offline-completeness is enforced only when the *pin or tooling* changes.

- **`txpack-future.yml` — build gates.**
  - *future-build* (every push): one pack-seeded `--online` build. A no-new-code push touches
    the network zero times; intentional new codes surface as an "N new terminology questions"
    signal; the signature is reported, not gated.
  - *pinned-world-hermetic* (only when `fhir.lock` or the wrapper changes): strict hermetic
    cold build with zero-network **and** signature **both enforced** — proves the pin is
    offline-complete.
  - *reproducibility* (`[parity]`): runs the judge vs `ref.manifest`. *determinism-gate*
    (`[determinism]`): two same-commit builds must be byte-identical, no excuses. *stock-baseline*
    (manual): the legacy build, for timing comparison.

- **`txpack-record.yml` — the nightly recorder.** Cron job that records → reproduces → publishes
  a candidate pack to this repo's own `txpack-store` release, then emits a `refresh-request`
  artifact. A manual `full_pin` job does a from-scratch [pin bootstrap](#concepts), hermetic-verified.

- **`txpack-refresh.yml` — the lock-bump PR.** Chained off the recorder via `workflow_run` (so no
  cross-repo token is needed): canonically diffs the candidate vs the pinned pack, builds an A/B
  output-impact report, has an LLM write a plain-English summary *from the diffs only*, and opens a
  PR that mutates **only `fhir.lock`**.

**Where artifacts live (all no-token, public `curl`, verified-on-use):** the answer **pack**
publishes to *this* repo's own `txpack-store` release (default `GITHUB_TOKEN`); the build **tool
jar** lives on the `jmandel/fhir-perf` release (human-bumped). Consumers need no token to read either.

---

## What changed in each fork (for the curious)

Three Git forks cooperate, all on `txpack-future`:

- **core (`org.hl7.fhir.core`) — the terminology engine.** Adds a read-only **pack seed layer**
  consulted before the live cache (misses fall through — it's a seed, not a wall); a **hermetic
  mode** whose violation is an `Error` (not an `Exception`, so tx clients can't swallow it);
  cache-key canonicalization so per-run labels don't poison replay; the pack build/merge/diff
  tooling with flake filtering; `TxLock` (reads `fhir.lock`, verifies integrity on every use); and
  a batch of determinism fixes.
- **kindling — the build front end.** Adds the **`SpecBuild` CLI** (`build`/`manifest`/`compare`/
  `impact`/`diff-packs`/`record`/`reproduce`) that replaces shell scripting, the cold **recorder**,
  routing all terminology through the cache so it's replayable, and an **honest judge** plus its
  own determinism fixes.
- **this repo (`jmandel/fhir`) — integration.** Holds the three pinned files, the unix alias, the
  FHIR settings, the judge's reference (`ref.manifest`) + allowlist, and the three CI workflows.

For the build-CLI flag-by-flag reference, run `./tools/build/build.sh help` (or any of
`build`/`manifest`/`compare`/`impact`/`diff-packs`/`record`/`reproduce`).

---

## The determinism contract (in brief)

Determinism is what makes every output comparison sound, and it's enforced in three layers:
*accidental* nondeterminism (shared state, hash-set ordering, torn parallel reads) is **fixed at
source**; *inherent* volatility (clocks, UUIDs, section numbers) is **normalized** in diffs and in
the judge; and *server* nondeterminism is **detected and degraded** via flake filtering, carry-forward,
and reproduce-before-propose. For the full treatment, see
[docs/txpack-vision.md](https://github.com/jmandel/fhir-perf/blob/main/docs/txpack-vision.md).
