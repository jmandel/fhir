#!/usr/bin/env bash
# Thin bootstrap for the Java CLI - the only thing this script does is fetch the pinned tooling
# jar (sha256-verified) and exec it. EVERYTHING else (tx.lock pack resolution, hermetic mode,
# locale/timezone pinning, signature check, judge/impact) lives in the Java CLI, so Windows
# users skip this file entirely:
#   1. download the jar from tx.lock's tooling.url, check its sha256
#   2. java -Xmx12g -cp kindling-future-v5.jar org.hl7.fhir.tools.publisher.SpecBuild build . [--judge|--impact|--online|--manifest]
set -euo pipefail
cd "$(dirname "$0")/../.."
lockval() {
  sed -n "/\"$1\"[[:space:]]*:/,/}/p" tx.lock | grep -m1 "\"$2\"" \
    | sed 's/^[^:]*:[[:space:]]*//; s/^"//; s/"\{0,1\},\{0,1\}[[:space:]]*$//'
}
TOOL_URL=$(lockval tooling url); TOOL_SHA=$(lockval tooling sha256)
TOOL="$HOME/.fhir/tools/kindling-future-$TOOL_SHA.jar"
if ! { [[ -f "$TOOL" ]] && echo "$TOOL_SHA  $TOOL" | sha256sum -c --quiet - 2>/dev/null; }; then
  mkdir -p "$(dirname "$TOOL")"
  echo "fetching $(basename "$TOOL") ..."
  curl -fsSL --retry 3 -o "$TOOL.part" "$TOOL_URL"
  echo "$TOOL_SHA  $TOOL.part" | sha256sum -c --quiet - || { echo "sha256 mismatch for $TOOL_URL" >&2; rm -f "$TOOL.part"; exit 1; }
  mv "$TOOL.part" "$TOOL"
fi
exec java -Xmx"${HEAP:-12g}" ${JAVA_OPTS:-} -cp "$TOOL" org.hl7.fhir.tools.publisher.SpecBuild build . "$@"
