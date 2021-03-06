#!/bin/bash
#


ADB=$(which adb.exe) # for windows-linux

set -ex
ADB=${ADB:-"adb"}

DEST="/data/local/tmp/realdevice-agent"

echo "Build binary for arm ..."
#GOOS=linux GOARCH=arm go build

#go generate
GOOS=linux GOARCH=arm go build -tags vfs

$ADB push atx-agent $DEST
$ADB shell chmod 755 $DEST
$ADB shell $DEST server --stop
#$ADB shell $DEST server -d "$@"
$ADB shell $DEST server -d --nouia "$@"

$ADB forward tcp:7912 tcp:7912
curl localhost:7912/wlan/ip
