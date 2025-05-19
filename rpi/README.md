# Configure Raspberry PI

NOTE: Raspberry Pi should be accessible on `pi.local` address from the browser.

**Steps:**
* Install the latest Raspbian OS using `Raspberry Pi Imager.app`.

* Execute the `startup.sh` shell script to disable wifi.

* Pin the desired IP Address of `raspberrpi` in admin console of router, `LAN` setting tab.

* Execute the `bootstrap.sh` shell script to update the system, install Docker and minikube.

* Use `sudo raspi-config` to open settings and enable `VNC`. Its under `Interfacer` option.

* Download [Pia VPN](https://www.privateinternetaccess.com/installer/x/download_installer_linux_arm/arm64). And install it by executing the downloaded script.

* Configure Kill Switch in `Pia` settings.

* Enable running Pia in the background by running: `piactl background enable`. Pia automatically installs systemd service, which starts it up after reboot. Because Pia can now run in background, it can enforce the Kill Switch, which disables connectivity to the internet.