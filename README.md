# Sprinter

![](preview.png)

## Features

- Monitor network uptime
- Pings & Speed test

Built with [Go](https://go.dev/)and [SQLite](https://sqlite.org/).

## Docker

#### Build

```shell
docker build -t sprinter:latest .
```

#### Run

```shell
docker run \
  -v "$(pwd)/data.db:/data.db" \
  -v "$(pwd)/.env:/.env" \
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
go get
```

#### Build & hot reload

```shell
gow -e=go,html run .
```

## License

Released under the [MIT License](LICENSE.md).

## Authors

- [Roman Zipp](https://github.com/romanzipp)
