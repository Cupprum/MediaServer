# Setup MacOS to host torrent downloading docker based center

* Start with setting up mac, learn more in `bootstrap` folder.

* Create a new colima instance: `colima start --cpu 2 --memory 8`; in the `~/.colima/default/colima.yaml` update colima 'mounts' to also mount `/Volumes/T5/MediaServer` folder, remember to restart the VM afterwards: `colima restart`. Note: specifying mounts removes the default mounts, which are `/Users/x42` and `/tmp/colima`. `colima.yaml` can be copied over to `~/.colima/default/colima.yaml` for easy configuration.

* Access colima VM: `colima ssh`

* Create the following directories in the colima VM:
    * /provision/jellyfin/config
    * /provision/jellyfin/cache
    * /provision/qbittorrent/config

* Run the setup, this step should be executed from mac, but its also fine from VM: `docker compose up -d`

* Access Jellyfin: [http://localhost:8096](http://localhost:8096)

* Enroll Jellyfin as a server, configure movie directory and so on.

* Access QBittorrent: [http://localhost:8080](http://localhost:8080) username is `admin` and password can be seen in logs: `docker logs qbittorrent`.

* Change password.

* Configure QBittorrent > Options > Downloads > `Default Save Path` to `/data/torrents`. This will download the files onto the SSD.

* Configure QBittorrent > Options > BitTorrent > `Seeding Limits`. Enable all three options, let change the last field to `Remove torrent`.