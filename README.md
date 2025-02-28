# Setup MacOS to host torrent downloading docker based center

* Start with setting up mac, learn more in `bootstrap` folder.

* Create a new colima instance: `colima start --cpu 2 --memory 8`; in the `~/.colima/default/colima.yaml` update colima 'mounts' to also mount `/Volumes/T5/MediaServer` folder, remember to restart the VM afterwards: `colima restart`. Note: specifying mounts removes the default mounts, which are `/Users/x42` and `/tmp/colima`. `colima.yaml` can be copied over to `~/.colima/default/colima.yaml` for easy configuration.

* Access colima VM: `colima ssh`

* Create the following directories in the colima VM:
    * /provision/jellyfin/config
    * /provision/jellyfin/cache
    * /provision/qbittorrent/config
    * /provision/prowlarr/config

* Configure Jellyfin, learn more: `jellyfin_setup/README.md`.

* Configure pi-hole.