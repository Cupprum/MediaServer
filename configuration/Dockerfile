FROM golang:1.25-alpine as builder

RUN apk add --no-cache busybox-static

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/configure_services .
RUN CGO_ENABLED=0 GOOS=linux go test -c -o /app/bin/qbittorrent.test ./qbittorrent
RUN CGO_ENABLED=0 GOOS=linux go test -c -o /app/bin/prowlarr.test ./prowlarr
RUN CGO_ENABLED=0 GOOS=linux go test -c -o /app/bin/jellyfin.test ./jellyfin

FROM scratch

COPY --from=builder /bin/busybox.static /bin/sh
COPY --from=builder /app/bin /app/bin

ENTRYPOINT ["/bin/sh", "-c"]