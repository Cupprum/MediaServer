# Configure Media Server

The Media Server is based on Jellyfin, QBittorrent, Prowlarr and Heimdall.

## Installation process

* Install the ArgoCD applications, which configure the services using `install.sh` shell script.

* Access Jellyfin: [http://192.168.0.100:8096](http://192.168.0.100:8096)

* Enroll Jellyfin as a server, configure movie directory and so on.

* Install Jellyfin app on Apple TV; configure it to access `http://192.168.0.100:8096`. The username is `root` and password can be found in stored passwords in Safari.

* Access QBittorrent: [http://192.168.0.100:8080](http://192.168.0.100:8080) username is `admin` and password can be seen in logs: `kubectl logs qbittorrent`.

* Change password.

* Configure QBittorrent > Options > Downloads > `Default Save Path` to `/data/torrents`. This will download the files onto the SSD.

* Configure QBittorrent > Options > BitTorrent > `Seeding Limits`. Enable all three options, let change the last field to `Remove torrent`.

* Access Prowlarr: [http://192.168.0.100:9696](http://192.168.0.100:9696) and setup username and password (Note: pro-tip - disable authentication on local network in general settings)

* Configure Prowlarr to use QBittorrent as a Download client, it is possible to configure qbittorrent to download films sequantally, so they can be streamed while downloading.

* Configure Prowlarr with indexers. Possible list:
    * 1337x: Note -> change base url as the more basic ones are usually blocked.
    * Badass Torrents
    * EZTV
    * LimeTorrents
    * The Pirate Bay
    * YTS

* Access Heimdall: [http://192.168.0.100:8088](http://192.168.0.100:8088) and configure access to the rest of the services; add the following `applications` to the `application list`:
    * Prowlarr -> Application Type: Prowlarr -> Set URL to: `http://192.168.0.100:9696`
    * QBittorrent -> Application Type: QBittorrent -> Set URL to: `http://192.168.0.100:8080`
    * Jellyfin -> Application Type: Jellyfin -> Add Application name and download logo from internet -> Set URL to: `http://192.168.0.100:8096`

## Deletion process

Delete the argocd apps with their configuration: `./delete.sh`