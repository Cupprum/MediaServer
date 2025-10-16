# Setup Raspberry Pi to host torrent downloading docker based media center

* Copy over public SSH key of my main device to Raspberry Pi: `ssh-copy-id -i ~/.ssh/id_rsa.pub x42@pi.local`

* Create a fine grained PAT in GitHub and store it as a `GH_TOKEN` Env var on Raspberry Pi. Make sure its configured to allow read only to `MediaServer` repository.

* SSH to Raspberry Pi and setup SSH key: `ssh-keygen` and press enter couple of times. Copy the public key and add it to GitHub: `cat ~/.ssh/id_rsa.pub`.

* Execute the `setup.sh` shell script located in the root folder of the project on developer computer, or setup each part manually from the raspberry pi by going through the following steps:

    * Start with setting up Raspberry Pi: `rpi` folder.

    * Install the Kubernetes cluster: `k8s` folder.

    * Install the whole Media server Jellyfin setup: `server` folder.

    * Install Grafana: `monitoring` folder.