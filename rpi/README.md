# Configure Raspberry Pi

Note: Raspberry Pi should be accessible on `pi.local` address from the browser.

**Steps:**
* Install the latest Raspbian OS using `Raspberry Pi Imager.app`.

* Pin the desired IP Address of `raspberrpi` in admin console of router, specifically `LAN` setting tab. I want the Raspberry Pi be accessible under `192.168.0.100`.

* Execute the `setup.sh` shell script to update the system, install Docker, minikube and disable wifi.

* Configure Kill Switch in `Pia` settings.

Note: Pia automatically installs systemd service, which starts it up after reboot. Because Pia can now run in background, it can enforce the Kill Switch, which disables connectivity to the internet.