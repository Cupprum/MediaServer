# Initial configuration of Raspberry PI

## Configure Raspberry Pi

Execute the `/init/setup-rpi.sh` shell script to configure the device.

**Note:** Pia automatically installs systemd service, which starts it up after reboot. Because Pia can now run in background, it can enforce the Kill Switch, which disables connectivity to the internet.

## Configure Kubernetes

After the Raspberry Pi reboots, execute the `/init/setup-k8s.sh` shell script to install tooling around Kubernetes and configure ArgoCD.