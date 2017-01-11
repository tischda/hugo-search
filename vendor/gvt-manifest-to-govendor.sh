#!/bin/sh
# ----------------------------------------------------------------------
# convert GVT manifest to GOVENDOR fetch commands
# ----------------------------------------------------------------------

grep -e importpath -e revision bleve-*-manifest | sed -e 's/\s*"importpath": ./govendor fetch /g' -e 's/..$//g' -e 's/"revision": ./@/g' | sed 'N;s/\n\s*//'
