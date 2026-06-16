#!/bin/bash
# ----------------------------------------------------------------------------------------------
# LEGACY stock build (svn + ant + publish.sh; ~20 min, requires a live terminology server).
# On the txpack-future branch this is NOT the front door. The fast, offline, JDK-only build is:
#
#     java -jar tools/build/launch.jar build .      # any OS; ~4 min; zero terminology network
#     ./tools/build/build.sh                        # unix alias for the same command
#
# Full guide: tools/build/README.md  (concepts, CLI, --online for new codes, CI, the refresh loop).
# This legacy script is kept for the existing Azure pipelines and is otherwise superseded.
# ----------------------------------------------------------------------------------------------
set -ev

NAME="Continuous Integration Build"
SVNREV=$(git log -1 | grep svn.fhir. | sed -r 's/^.*?@([0-9]+).*$/\1/')

antBuild (){
  ./publish.sh -svn $SVNREV -name \'$NAME\' -url http://build.fhir.org/ -nocheck
  checkStatus
}

checkStatus (){
  nf=`find publish -maxdepth 1 -type f | wc -l`
  if [ "$nf" -lt "100" ] ; then
     echo "< 100 files produced: bailing!"
     exit 1
  fi
  if [ $? -eq 0 -a ! -f fhir-error-dump.txt ]
  then
    echo "Build status OK"
  else
    echo "error dump:"
    cat fhir-error-dump.txt
    exit 1
  fi
}

antBuild
