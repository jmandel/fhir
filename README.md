## The FHIR Specification Publisher

This library builds and publishes the FHIR specification, based on the contained spreadsheet data in the project.

| CI Status ([master][Link-BuildFhirOrgMaster]) | CI Status ([R4B][Link-BuildFhirOrgR4B]) |
| :---: | :---: |
| [![Build Status][Badge-AzureMasterPipeline]][Link-AzureMasterPipeline] | [![Build Status][Badge-AzureR4BPipeline]][Link-AzureR4BPipeline] |

## Building the spec (txpack future)

This branch (`txpack-future`) ships a new, self-contained build front door. **A JDK is the
entire dependency footprint, on any OS** — no Gradle, no Ant, no `git-svn`, no local
terminology server. From the spec root:

```sh
java -jar tools/build/launch.jar build .
```

That is a fully **cold** build that finishes in **~3.5 minutes** (12-core; ~4 min on a small
CI runner) and makes **zero terminology network requests** — provably, because any attempt is a
hard failure naming the offending request. The build's output signature is checked against
[`fhir.lock`](fhir.lock) (`errors=0, warnings=3694, information=349`); by default a mismatch
fails the build. (To investigate a deliberate change, run with `SIGNATURE_GATE=report` — or
`-Dorg.hl7.fhir.spec.signatureGate=report` — to print `SIGNATURE CHANGED (reported, not
enforced)` and exit 0 instead of failing.)

The launcher needs nothing but a JDK on `PATH`. It reads the toolchain pin in
`tools/build/kindling-wrapper.properties`, downloads the pinned build tooling once
(sha256-verified, cached in `~/.fhir/tools`), and runs it. Override the heap with the `HEAP`
environment variable (default `12g`):

```sh
HEAP=11g java -jar tools/build/launch.jar build .
```

`tools/build/build.sh` is a 4-line unix alias for the same command (it `cd`s to the spec root
and `exec`s the launcher), so `./tools/build/build.sh` is equivalent.

### Common variations

| Goal | Command |
|---|---|
| Hermetic build (default; zero terminology traffic) | `java -jar tools/build/launch.jar build .` |
| Allow pack misses to reach the server (e.g. new codes) | `java -jar tools/build/launch.jar build . --online` |
| Report which published files your edit changed | `java -jar tools/build/launch.jar build . --impact` |
| Reproducibility judge (vs the committed reference manifest) | `java -jar tools/build/launch.jar build . --judge` |

### Personas

Editors run the one command above and nothing else — content PRs carry content only.
**The only writer of [`fhir.lock`](fhir.lock) is the nightly refresh bot**; editors never touch
it. The miss count under `--online` surfaces in the CI log as a signal, not a gate.

### Learn more

- [`tools/build/README.md`](tools/build/README.md) — quickstart for the new build.
- [`tools/build/BUILD-REFERENCE.md`](tools/build/BUILD-REFERENCE.md) — the full command surface
  (every subcommand, flag, and system property).
- [`tools/build/HOW-IT-WORKS.md`](tools/build/HOW-IT-WORKS.md) — what changed in each fork to make
  the build hermetic, fast, and deterministic.
- [`tools/build/REFRESH-AND-CI.md`](tools/build/REFRESH-AND-CI.md) — the nightly refresh/CI
  pipeline that regenerates packs and `fhir.lock`.
- [`FUTURE.md`](FUTURE.md) — the why: the txpack design, trust model, and refresh pipeline.

---

> **Legacy / stock toolchain.** Everything below documents the original Gradle/Ant publisher.
> It still works, but it is slow (~20 min) and depends on live network access to a terminology
> server. New work should use the `build` front door above.

### Important Links

This is the source for the FHIR specification itself. Only the editors of
the specification (a small group) need to build this. If that's not you,
one of these links should get you going:

* [Jira - Propose a change](https://jira.hl7.org/projects/FHIR/issues) - use this rather than making a PR directly, since all changes must be approved using the Jira workflow
* [FHIR chat](https://chat.fhir.org)
* [Stack Overflow questions](https://stackoverflow.com/tags/hl7-fhir)
* [Published FHIR Specification](http://hl7.org/fhir) or [Current Build of the specficiation](http://build.fhir.org)

### Publishing Locally (legacy)

1. Be sure you have at least 16 GB of RAM
2. Run `./gradlew publish` from the command line
3. Wait for it to finish (~20 minutes)

See also: [Getting Started][Link-Wiki] and [FHIR Build Process][Link-Confluence]

##### If running commands on the terminal is a frightening prospect for you...

We provide executable script files for Windows (publish.bat) and for a Bash shell for mac/linux/windows (publish.sh).

### Command line parameters (legacy)

There are multiple options available for publishing:

 * `--offline`: use this arg if you are offline and cannot fetch dependencies (doesn't work at the moment, and may never)

 * `-nogen`: don't generate the spec, just run the validation. (to use this,
   manually fix things in the publication directory, and then migrate the
changes back to source when done. this is a hack)

 * `-noarchive`: don't generate the archive. Don't use this if you're a core
   editor

 * `-web`: produce the HL7 ready publication form for final upload (only core
   editors)

 * `-diff`: the executable program to use if platform round-tripping doesn't
   produce identical content (default: c:\program files
(x86)\WinMerge\WinMergeU.exe)

 * `-name`: the "name" to go in the title bar of each of the specification

To add any of these options to the publish task, run the command as `./gradlew publish --args"<YOUR ARGS HERE>"`

For example, if you wanted to publish without generating the spec, just running the validation, you would run the command `./gradlew publish --args="-nogen"`

### Publishing Globally

Each time a pull request is opened, the [pull request pipeline][Link-AzurePRPipeline] runs. If the pipeline successfully publishes, it uploads the build as a
separate branch on [build.fhir.org/branches][Link-BuildFhirOrgBranches], where it can be reviewed to ensure accuracy.

Once merged to master, the [master branch pipeline][Link-AzureMasterPipeline] runs. If successful, the published specification is uploaded to the main
[build.fhir.org][Link-BuildFhirOrgMaster] webpage.

The only exception to the above is the build for `R4B`. The [R4B pipline][Link-AzureR4BPipeline] detects changes to the [R4B branch][Link-R4BGithub] in github, and
publishes any changes from that branch to [build.fhir.org/R4B][Link-BuildFhirOrgR4B].

### Maintenance
This project is maintained by [Grahame Grieve][Link-grahameGithub] and [Mark Iantorno][Link-markGithub] on behalf of the FHIR community.

---

[Link-AzureMasterPipeline]: https://dev.azure.com/fhir-pipelines/fhir-publisher/_build/latest?definitionId=44&branchName=refs%2Fpull%2F1084%2Fmerge
[Link-AzureR4BPipeline]: https://dev.azure.com/fhir-pipelines/fhir-publisher/_build/latest?definitionId=46&branchName=R4B
[Link-AzurePRPipeline]: https://dev.azure.com/fhir-pipelines/fhir-publisher/_build/latest?definitionId=42&branchName=refs%2Fpull%2F1084%2Fmerge
[Link-BuildFhirOrgMaster]: https://build.fhir.org
[Link-BuildFhirOrgBranches]: https://build.fhir.org/branches/
[Link-BuildFhirOrgR4B]: https://build.fhir.org/branches/R4B/
[Link-Wiki]: https://github.com/hl7/fhir/wiki/Get-Started-with-FHIR-on-GitHub
[Link-Confluence]: https://confluence.hl7.org/display/FHIR/FHIR+Build+Process
[Link-R4BGithub]: https://github.com/HL7/fhir/tree/R4B
[Link-grahameGithub]: https://github.com/grahamegrieve
[Link-markGithub]: https://github.com/markiantorno
[Badge-AzureMasterPipeline]: https://dev.azure.com/fhir-pipelines/fhir-publisher/_apis/build/status/Master%20Branch%20Pipeline?branchName=refs%2Fpull%2F1084%2Fmerge
[Badge-AzureR4BPipeline]: https://dev.azure.com/fhir-pipelines/fhir-publisher/_apis/build/status/R4B%20Pipeline?branchName=R4B
