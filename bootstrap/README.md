# Configure Mac

**Steps:**

* Install [brew](https://brew.sh)

* Install PIA, colima, docker and tmux: `brew install private-internet-access colima docker tmux`

* Create a new tmux session: `tmux new -s keepawake` and stay awake: `sudo caffeinate -i`. Detach from the tmux session by pressing Ctrl+B followed by D. NOTE: This is required to do after restart. To attach to the session use: `tmux attach -t keepawake`.

TODO: Is this even required? doesnt it stay available either way?

* In System Preferences under `Lock Screen` configure display to never turn off.

* Enable Kill Switch and Advanced Kill Switch opitions in PIA Application GUI in the `Privacy` part of settings. This will allow network traffic to flow only through the VPN, if the connection to VPN is lost, no network traffic will work. I can still however access the server through SSH on local network, even if the VPN is disconnected.

* Configure PIA to run automatically on start up using the `configure_pia.sh` script.