# Setup Raspberry PI to host torrent downloading docker based media center

* Setup SSH key and add it to GitHub.

* Clone the project: `git clone git@github.com:Cupprum/MediaServer.git`

* Execute the `setup.sh` shell script located in the root folder of the project, or setup each part manually from the raspberry pi by going through the following steps:

    * Start with setting up Raspberry PI, read more in `rpi` folder.

    * Install the Kubernetes cluster, read more in the `k8s` folder.

    * Install the whole Media server Jellyfin setup, read more in the `server` folder.

    * Install Grafana, read more in the `monitoring` folder.