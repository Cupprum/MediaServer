curl -sSfi "http://jellyfin.pi.local/System/Info"

curl -sSfi -X POST "http://jellyfin.pi.local/Startup/Configuration" \
  -H "Content-Type: application/json" \
  -d '{
    "UICulture": "en-US",
    "MetadataCountryCode": "US",
    "PreferredMetadataLanguage": "en"
  }'

curl -sSfi "http://jellyfin.pi.local/Startup/User"

curl -sSfi -X POST "http://jellyfin.pi.local/Startup/User" \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "admin",
    "Password": "jellyfin"
  }'

curl -sSfi -X POST "http://jellyfin.pi.local/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies" \
  -H "Content-Type: application/json" \
  -d '{
    "LibraryOptions": {
      "TypeOptions": [
        {
          "Type": "Movie",
          "MetadataFetchers": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "MetadataFetcherOrder": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "ImageFetchers": [
            "TheMovieDb",
            "The Open Movie Database",
            "Embedded Image Extractor",
            "Screen Grabber"
          ],
          "ImageFetcherOrder": [
            "TheMovieDb",
            "The Open Movie Database",
            "Embedded Image Extractor",
            "Screen Grabber"
          ]
        }
      ],
      "PathInfos": [
        {
          "Path": "/media/movies"
        }
      ]
    }
  }'


curl -sSfi -X POST "http://jellyfin.pi.local/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows" \
  -H "Content-Type: application/json" \
  -d '{
    "LibraryOptions": {
      "TypeOptions": [
        {
          "Type": "Series",
          "MetadataFetchers": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "MetadataFetcherOrder": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "ImageFetchers": [
            "TheMovieDb"
          ],
          "ImageFetcherOrder": [
            "TheMovieDb"
          ]
        },
        {
          "Type": "Season",
          "MetadataFetchers": [
            "TheMovieDb"
          ],
          "MetadataFetcherOrder": [
            "TheMovieDb"
          ],
          "ImageFetchers": [
            "TheMovieDb"
          ],
          "ImageFetcherOrder": [
            "TheMovieDb"
          ]
        },
        {
          "Type": "Episode",
          "MetadataFetchers": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "MetadataFetcherOrder": [
            "TheMovieDb",
            "The Open Movie Database"
          ],
          "ImageFetchers": [
            "TheMovieDb",
            "The Open Movie Database",
            "Embedded Image Extractor",
            "Screen Grabber"
          ],
          "ImageFetcherOrder": [
            "TheMovieDb",
            "The Open Movie Database",
            "Embedded Image Extractor",
            "Screen Grabber"
          ]
        }
      ],
      "PathInfos": [
        {
          "Path": "/media/tv"
        }
      ]
    }
  }'



curl -sSfi -X POST "http://jellyfin.pi.local/Startup/RemoteAccess" \
  -H "Content-Type: application/json" \
  -d '{
    "EnableRemoteAccess":false,
    "EnableAutomaticPortMapping":false
  }'

curl -sSfi -X POST "http://jellyfin.pi.local/Startup/Complete"