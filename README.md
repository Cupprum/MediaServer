# Setup Raspberry Pi to host torrent downloading docker based media center

Install the latest Raspbian OS using `Raspberry Pi Imager.app`.

Pin the desired IP Address of `raspberrypi` in admin console of router, specifically `LAN` setting tab. I want the Raspberry Pi be accessible under `192.168.0.100`.

Copy over public SSH key of my main device to Raspberry Pi: `ssh-copy-id -i ~/.ssh/id_rsa.pub x42@pi.local`

Create a fine grained PAT in GitHub and store it as a `MEDIASERVER_GITHUB_TOKEN` Env var on Raspberry Pi. Make sure its configured to allow read only to `MediaServer` repository.

SSH to Raspberry Pi and setup SSH key: `ssh-keygen` and press enter couple of times. Copy the public key and add it to GitHub: `cat ~/.ssh/id_rsa.pub`.

Execute the `setup.sh` shell script located in the root folder of the project on developer computer.
This script will configure the Raspberry Pi, deploy kubernetes cluster into it and spin up all services to run and monitor MediaServer.