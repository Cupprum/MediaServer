# Configure Media Server

The Media Server is based on Jellyfin, QBittorrent, Prowlarr and Heimdall.

## Installation and configuration

* The server required presence of some folders for storing its configuration, these can be created using the `./bootstrap.sh` shell script.

* Install the ArgoCD applications, which configure the services using `setup.sh` shell script, specifically `./setup.sh install`.

  TODO: change IPs to routes

* Access Jellyfin: [http://jellyfin.pi.local/](http://jellyfin.pi.local/)

* Enroll Jellyfin as a server, configure movie directory and so on.

* Install Jellyfin app on Apple TV; configure it to access `http://jellyfin.pi.local/`. The username is `admin` and password can be found in stored passwords in Safari.

* Access QBittorrent: [http://qbittorrent.pi.local/](http://qbittorrent.pi.local/) username is `admin` and password can be seen in logs: `kubectl logs qbittorrent`.

* Change password for QBittorrent.

* Configure QBittorrent > Options > BitTorrent > `Seeding Limits`. Enable all three options, let change the last field to `Remove torrent`.

* Access Prowlarr: [http://prowlarr.pi.local/](http://prowlarr.pi.local/) and setup username and password (Note: pro-tip - disable authentication on local network in general settings)

* Configure Prowlarr to use QBittorrent as a Download client, use the following host: `qbittorrent.server.svc.cluster.local`, it is possible to configure qbittorrent to download films sequantally, so they can be streamed while downloading.

* Configure Prowlarr with indexers. Possible list:
    * 1337x: Note -> change base url as the more basic ones are usually blocked.
    * Badass Torrents: Note -> seems broken
    * EZTV
    * LimeTorrents
    * The Pirate Bay
    * YTS

* Access Heimdall: [http://heimdall.pi.local/](http://heimdall.pi.local/) and configure access to the rest of the services; add the following `applications` to the `application list`:
    * Prowlarr -> Application Type: Prowlarr -> Set URL to: `http://prowlarr.pi.local/`
    * QBittorrent -> Application Type: QBittorrent -> Set URL to: `http://qbittorrent.pi.local/`
    * Jellyfin -> Application Type: Jellyfin -> Set URL to: `http://jellyfin.pi.local/`
    * Argo CD -> Application Type: Argo CD -> Set URL to: `http://argocd.pi.local/`
    * Grafana -> Application Type: Grafana -> Set URL to: `http://grafana.pi.local/`

## Deletion process

Delete the argocd apps with their configuration: `./setup.sh delete`