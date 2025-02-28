# Configure Jellyfin setup

* Run the setup, this step should be executed from mac, but its also fine from VM: `docker compose up -d`

* Access Jellyfin: [http://localhost:8096](http://localhost:8096)

* Enroll Jellyfin as a server, configure movie directory and so on.

* Access QBittorrent: [http://localhost:8080](http://localhost:8080) username is `admin` and password can be seen in logs: `docker logs qbittorrent`.

* Change password.

* Configure QBittorrent > Options > Downloads > `Default Save Path` to `/data/torrents`. This will download the files onto the SSD.

* Configure QBittorrent > Options > BitTorrent > `Seeding Limits`. Enable all three options, let change the last field to `Remove torrent`.

* Access Prowlarr: [http://localhost:9696](http://localhost:9696) and setup username and password (Note: pro-tip disable authentication on local network in general settings)

* Configure Prowlarr to use QBittorrent as a Download client, it is possible to configure qbittorrent to download films sequantally, so they can be streamed while downloading.

* Configure Prowlarr with indexers. Possible list:
    * 1337x: Note -> change base url as the more basic ones are usually blocked.
    * Badass Torrents
    * EZTV
    * LimeTorrents
    * The Pirate Bay
    * YTS
