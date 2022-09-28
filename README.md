# Sprinter

![](preview.png)

## Features

- Continuously monitors network availability with Pings to configurated servers
- Web UI with interactive charts

Built with [Go](https://go.dev/) and [SQLite](https://sqlite.org/).

## Usage (Docker)

```
docker pull ghcr.io/romanzipp/sprinter:latest
```

See [repository](https://github.com/romanzipp/Sprinter/pkgs/container/sprinter) for more information.

### Environment variables

- `INTERVAL` Ping interval in seconds (default: `60`)
- `PING_HOSTS` Destination servers, separated by comma (default: `1.1.1.1,google.com`)
- `PING_TIMEOUT` Ping timeout in seconds (default: `5`)
- `PING_PRIVILEGED` Privileged UDP ping mode (see [Troubleshooting](#troubleshooting))

---

### Local

#### Build

```shell
docker build -t sprinter:latest .
```

#### Run

```shell
docker run \
  -v "$(pwd)/data/:/data/" \
  -v "$(pwd)/.env:/.env" \
  -p 8080:8080 \
  sprinter:latest
```

## Development

### Requirements

- Go 1.19+
- Yarn
- _Docker_

### Go app

#### Install dependencies

```
go mod download
```

#### Build & hot reload

```shell
gow -e=go,html run .
```

### Frontend

#### Install dependencies

```
yarn install
```

#### Build & hot reload

```shell
yarn watch
```
## Roadmap

- [ ] Make ping timeout configurable

## Troubleshooting

**macOS**: You will need to run the executable as root and enable privileged mode by setting `PING_PRIVILEGED` to true.

## License

Released under the [MIT License](LICENSE.md).

## Authors

- [Roman Zipp](https://github.com/romanzipp)
