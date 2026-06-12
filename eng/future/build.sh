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

HERMETIC=true; JUDGE=false; MANIFEST=false; IMPACT=false
for a in "$@"; do
  case "$a" in
    --hermetic) HERMETIC=true;;
    --online)   HERMETIC=false;;
    --judge)    JUDGE=true; MANIFEST=true;;
    --manifest) MANIFEST=true;;  # write the output manifest without judging (CI convergence pass)
    --impact)   IMPACT=true; MANIFEST=true;;  # diff this build against my previous local build
    *) echo "unknown arg: $a" >&2; exit 1;;
  esac
done

# tx.lock parsing in plain shell - the default build path requires only bash, curl, sha256sum
# and java (in the non-demo future the publisher reads tx.lock natively and no script exists).
# python3 is needed only by the verification tooling (--judge / --manifest / CI / refresh).
lockval() { # lockval SECTION KEY -> value (string or number)
  sed -n "/\"$1\"[[:space:]]*:/,/}/p" tx.lock | grep -m1 "\"$2\"" \
    | sed 's/^[^:]*:[[:space:]]*//; s/^"//; s/"\{0,1\},\{0,1\}[[:space:]]*$//'
}
PACK_URL=$(lockval pack url);        PACK_SHA=$(lockval pack zipSha256)
TOOL_URL=$(lockval tooling url);     TOOL_SHA=$(lockval tooling sha256)
EXP_E=$(lockval expectedSignature errors); EXP_W=$(lockval expectedSignature warnings); EXP_I=$(lockval expectedSignature information)

fetch() { # fetch URL SHA DEST — content-addressed download: verify-or-refetch, never mutate
  local url="$1" sha="$2" dest="$3"
  if [[ -f "$dest" ]] && echo "$sha  $dest" | sha256sum -c --quiet - 2>/dev/null; then return 0; fi
  mkdir -p "$(dirname "$dest")"
  echo "fetching $(basename "$dest") ..."
  curl -fsSL --retry 3 -o "$dest.part" "$url"
  echo "$sha  $dest.part" | sha256sum -c --quiet - || { echo "sha256 mismatch for $url" >&2; rm -f "$dest.part"; exit 1; }
  mv "$dest.part" "$dest"
}

# TXPACK_ZIP overrides the lock's pack with a local zip (used by the refresh workflow's
# old-vs-new output A/B; everything else still comes from the lock)
if [[ -n "${TXPACK_ZIP:-}" ]]; then
  PACK="$TXPACK_ZIP"
else
  PACK="$HOME/.fhir/tx-packs/$PACK_SHA.zip"
  fetch "$PACK_URL" "$PACK_SHA" "$PACK"
fi
TOOL="$HOME/.fhir/tools/kindling-future-v2-$TOOL_SHA.jar"
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

if $IMPACT; then
  if [[ -z "$PREV_MANIFEST" ]]; then
    echo "impact: no previous build manifest found (run a build first); nothing to compare"
  else
    echo "== impact: files changed vs your previous build (content-level; ordering-only churn excluded)"
    join -t$'\t' -j2 <(sort -t$'\t' -k2 "$PREV_MANIFEST") <(sort -t$'\t' -k2 build-future.manifest) \
      | awk -F'\t' '{split($2,a,"/"); split($3,b,"/"); if ($2!=$3 && (length(a)<2 || a[2]!=b[2])) print "  " $1}'
    comm -13 <(cut -f2- "$PREV_MANIFEST" | sort) <(cut -f2- build-future.manifest | sort) | sed 's/^/  + /'
    comm -23 <(cut -f2- "$PREV_MANIFEST" | sort) <(cut -f2- build-future.manifest | sort) | sed 's/^/  - /'
  fi
fi

if $JUDGE; then
  echo "== judging published output against the committed reference manifest"
  # ordering second chance: a file whose exact (N:) hash differs but whose order-insensitive
  # (O:) hash matches changed only in element order - the documented stock nondeterminism
  # class - and is excused with that evidence; any content change differs under O too
  join -t$'\t' -j2 <(sort -t$'\t' -k2 eng/future/ref.manifest) <(sort -t$'\t' -k2 build-future.manifest) \
    | awk -F'\t' '{split($2,a,"/"); split($3,b,"/"); if ($2!=$3 && (length(a)<2 || a[2]!=b[2])) print $1}' > /tmp/future.content-diff
  candidates=$(grep -vE '\.shex(\.html)?$|\.xls$' /tmp/future.content-diff | grep -vx 'all-valuesets.zip' || true)
  if [[ -n "$PREV_MANIFEST" ]]; then
    join -t$'\t' -j2 <(sort -t$'\t' -k2 "$PREV_MANIFEST") <(sort -t$'\t' -k2 build-future.manifest) \
      | awk -F'\t' '$2!=$3{print $1}' | sort > /tmp/noisy-now.txt
    # evidence first; the historically-observed allowlist additionally covers files whose
    # nondeterminism sampled differently across machines but happened to be stable within this
    # run's two builds (residual blind spot: static-listed files stable this run - see FUTURE.md)
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
  echo "byte parity: clean (content-level diffs: $(wc -l < /tmp/future.content-diff), all evidenced or known-nondeterministic)"
fi
