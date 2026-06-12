#!/usr/bin/env bash
# build.sh [--hermetic] [--judge] [--online] — the txpack future-world spec build.
#
# Resolves the build tooling and the terminology answer pack from tx.lock (content-addressed,
# sha256-verified, cached under ~/.fhir/tx-packs and ~/.fhir/tools), then runs the publisher.
#   --hermetic  prove the build needs zero terminology network (any request = hard failure)
#   --online    allow misses to fall through to the configured tx server (editor default
#               when adding genuinely new codes; without it, --hermetic is implied)
#   --judge     after the build, compare published output against the committed reference
#               manifest (eng/future/ref.manifest), ignoring known build nondeterminism
set -euo pipefail
cd "$(dirname "$0")/../.."

HERMETIC=true; JUDGE=false; MANIFEST=false
for a in "$@"; do
  case "$a" in
    --hermetic) HERMETIC=true;;
    --online)   HERMETIC=false;;
    --judge)    JUDGE=true; MANIFEST=true;;
    --manifest) MANIFEST=true;;  # write the output manifest without judging (CI convergence pass)
    *) echo "unknown arg: $a" >&2; exit 1;;
  esac
done

lock() { python3 -c "import json,sys;d=json.load(open('tx.lock'));print(eval('d'+sys.argv[1]))" "$1"; }

PACK_URL=$(lock "['pack']['url']");         PACK_SHA=$(lock "['pack']['zipSha256']")
TOOL_URL=$(lock "['tooling']['url']");      TOOL_SHA=$(lock "['tooling']['sha256']")
EXP_E=$(lock "['expectedSignature']['errors']"); EXP_W=$(lock "['expectedSignature']['warnings']"); EXP_I=$(lock "['expectedSignature']['information']")

fetch() { # fetch URL SHA DEST — content-addressed download: verify-or-refetch, never mutate
  local url="$1" sha="$2" dest="$3"
  if [[ -f "$dest" ]] && echo "$sha  $dest" | sha256sum -c --quiet - 2>/dev/null; then return 0; fi
  mkdir -p "$(dirname "$dest")"
  echo "fetching $(basename "$dest") ..."
  curl -fsSL --retry 3 -o "$dest.part" "$url"
  echo "$sha  $dest.part" | sha256sum -c --quiet - || { echo "sha256 mismatch for $url" >&2; rm -f "$dest.part"; exit 1; }
  mv "$dest.part" "$dest"
}

PACK="$HOME/.fhir/tx-packs/$PACK_SHA.zip"
TOOL="$HOME/.fhir/tools/kindling-future-v2-$TOOL_SHA.jar"
fetch "$PACK_URL" "$PACK_SHA" "$PACK"
fetch "$TOOL_URL" "$TOOL_SHA" "$TOOL"

# locale and timezone are part of the pinned configuration: the pack's request keys embed the
# display language (recorded as en-US) and the reference output embeds recording-zone dates
FLAGS=(-Xmx"${HEAP:-12g}" -XX:+UseParallelGC -Dfile.encoding=UTF-8
  -Duser.language=en -Duser.country=US -Duser.timezone=America/Chicago
  -Dorg.hl7.fhir.tx.maxConcurrency=12
  -Dorg.hl7.fhir.tx.localFirst=true
  -Dorg.hl7.fhir.tx.pack="$PACK")
$HERMETIC && FLAGS+=(-Dorg.hl7.fhir.tx.hermetic=true)

start=$(date +%s)
java "${FLAGS[@]}" -cp "$TOOL" org.hl7.fhir.tools.publisher.Publisher \
  -nosound -nopartial -fhir-settings eng/future/fhir-settings.json | tee build-future.log
rc=${PIPESTATUS[0]}
echo "BUILD rc=$rc duration=$(( $(date +%s) - start ))s"
[[ $rc -eq 0 ]] || exit $rc

SIG=$(grep -E "Summary: Errors=" build-future.log | tail -1)
echo "signature: $SIG (expected Errors=$EXP_E, Warnings=$EXP_W, Information messages=$EXP_I)"
echo "$SIG" | grep -q "Errors=$EXP_E, Warnings=$EXP_W, Information messages=$EXP_I" \
  || { echo "OUTPUT SIGNATURE MISMATCH" >&2; exit 1; }

# a manifest from a prior same-commit build (CI's convergence pass, --manifest) upgrades the
# judge from a static noise allowlist to EVIDENCE-BASED excusal: a file is only excused if it
# provably varied between two same-commit builds in this environment. A real content change is
# stable across builds and differs from the reference -> flagged, even in historically-noisy
# files. Manifest writing is opt-in (~14s) so default editor builds pay nothing.
PREV_MANIFEST=""
if $MANIFEST; then
  [[ -f build-future.manifest ]] && { PREV_MANIFEST=/tmp/prev.manifest; cp build-future.manifest "$PREV_MANIFEST"; }
  python3 eng/future/manifest.py publish > build-future.manifest
fi

if $JUDGE; then
  echo "== judging published output against the committed reference manifest"
  join -t$'\t' -j2 <(sort -t$'\t' -k2 eng/future/ref.manifest) <(sort -t$'\t' -k2 build-future.manifest) \
    | awk -F'\t' '$2!=$3{print $1}' > /tmp/future.diff-files
  # irreducible exclusions (content genuinely appears/disappears or is binary-timestamped;
  # both reported upstream): .shex, .xls, all-valuesets.zip
  candidates=$(grep -vE '\.shex(\.html)?$|\.xls$' /tmp/future.diff-files | grep -vx 'all-valuesets.zip' || true)
  if [[ -n "$PREV_MANIFEST" ]]; then
    join -t$'\t' -j2 <(sort -t$'\t' -k2 "$PREV_MANIFEST") <(sort -t$'\t' -k2 build-future.manifest) \
      | awk -F'\t' '$2!=$3{print $1}' | sort > /tmp/noisy-now.txt
    # evidence first; the historically-observed ordering allowlist additionally covers files
    # whose ordering nondeterminism sampled differently across machines but happened to be
    # stable within this run's two builds (the residual blind spot is static-listed files that
    # are stable this run - the documented stock-ordering class, see FUTURE.md)
    unexplained=$(echo "$candidates" | sort | comm -23 - /tmp/noisy-now.txt | grep -vxFf eng/future/noise-files-v2.txt | sed '/^$/d' || true)
    excused=$(echo "$candidates" | sort | comm -12 - /tmp/noisy-now.txt | sed '/^$/d' | wc -l)
    echo "(evidence-based excusal: $excused files varied between this run's two builds)"
  else
    # single local build: fall back to the historically-observed noise allowlist
    unexplained=$(echo "$candidates" | grep -vxFf eng/future/noise-files-v2.txt || true)
  fi
  if [[ -n "$unexplained" ]]; then
    echo "UNEXPLAINED OUTPUT DIFFS (beyond known build nondeterminism):"; echo "$unexplained"
    # bundle the flagged files for offline diagnosis (CI uploads this as an artifact)
    echo "$unexplained" | head -40 | (cd publish && tar -czf ../parity-debug.tgz -T - 2>/dev/null) || true
    exit 1
  fi
  echo "byte parity: clean ($(wc -l < /tmp/future.diff-files) files differ, all evidenced or known-nondeterministic)"
fi
