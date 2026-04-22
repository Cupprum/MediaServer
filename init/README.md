# Initial configuration of Raspberry PI

## Configure Raspberry Pi

Install the latest Raspbian OS using `Raspberry Pi Imager.app`.

Pin the desired IP Address of `raspberrypi` in admin console of router, specifically `LAN` setting tab. I want the Raspberry Pi be accessible under `192.168.0.100`.

Execute the `/init/setup-rpi.sh` shell script to configure the device.

**Note:** Pia automatically installs systemd service, which starts it up after reboot. Because Pia can now run in background, it can enforce the Kill Switch, which disables connectivity to the internet.

## Configure Kubernetes

After the Raspberry Pi reboots, execute the `/init/setup-k8s.sh` shell script to install tooling around Kubernetes and configure ArgoCD.