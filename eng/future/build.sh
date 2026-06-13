#!/bin/sh
# Convenience alias for unix shells. The real front door is the committed launcher:
#   java -jar eng/future/launch.jar build .         (any OS; JDK is the only dependency)
cd "$(dirname "$0")/../.." && exec java -jar eng/future/launch.jar build . "$@"
