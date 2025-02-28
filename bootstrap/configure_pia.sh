#!/usr/bin/env bash

set -e
set -u
set -o pipefail

echo "- Starting $0"

function main() {
    WorkDir="$(cd "$(dirname "$0")" || exit 1; pwd)"
    cd "$WorkDir" || exit 1
    echo "WorkDir is $WorkDir"

    echo '- Configuring PIA region'
    piactl set region dk-copenhagen

    echo '- Enable PIA to run in background'
    piactl background enable

    echo '- Copy LaunchAgent definition of PIA'
    cp com.privateinternetaccess.vpn.plist ~/Library/LaunchAgents/com.privateinternetaccess.vpn.plist

    echo '- Load PIA LaunchAgent definition'
    launchctl load ~/Library/LaunchAgents/com.privateinternetaccess.vpn.plist
}

main "$@"

echo "- Finished $0"
