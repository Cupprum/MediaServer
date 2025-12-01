- MediaServer:
    - qBittorrent:
        widget:
            type: "qbittorrent"
            url: "http://qbittorrent:8080"
            username: "$QBITTORRENT_USERNAME"
            password: "$QBITTORRENT_PASSWORD"
    - Prowlarr:
        widget:
            type: "prowlarr"
            url: "http://prowlarr:9696"
            key: "$PROWLARR_APIKEY"
    - Jellyfin:
        widget:
            type: jellyfin
            url: "http://jellyfin:8096"
            key: "$JELLYFIN_APIKEY"