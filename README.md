# unbound-container

[![Build and Push](https://github.com/trankimtung/unbound-container/actions/workflows/build.yml/badge.svg)](https://github.com/trankimtung/unbound-container/actions/workflows/build.yml)

A minimal container image for [Unbound](https://nlnetlabs.nl/projects/unbound/), a validating, recursive, caching DNS resolver.

Unbound is compiled from source with secure defaults suitable for use as a private DNS resolver in home labs and self-hosted infrastructure.

## Features

Out of the box, this Unbound container will:

- Run as non-root.
- Listen on UDP and TCP port 5335.
- Act as a recursive resolver with caching.
- Validate DNSSEC signatures.
- Only allow queries from RFC 1918 / RFC 4193 private ranges.
- Supported architectures: `amd64`, `arm64`.

## Images

Images are published to GitHub Container Registry:

```
ghcr.io/trankimtung/unbound-container
```

| Tag | Description |
|-----|-------------|
| `latest` | Latest stable version |
| `1.24.2` | Unbound 1.24.2 |

## Quick Start

```sh
docker run -d \
  --name unbound \
  -p 5335:5335/udp \
  -p 5335:5335/tcp \
  ghcr.io/trankimtung/unbound-container:latest
```

Test resolution:

```sh
dig @127.0.0.1 example.com
dig @127.0.0.1 example.com +dnssec
```

## Docker Compose

```yaml
services:
  unbound:
    image: ghcr.io/trankimtung/unbound-container:latest
    container_name: unbound
    restart: unless-stopped
    ports:
      - "53:53/udp"
      - "53:53/tcp"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `UNBOUND_PORT` | `5335` | Port unbound listens on |

## Access Control

By default, the resolver only accepts queries from private network ranges:

| Range | Description |
|-------|-------------|
| `127.0.0.1/32` | Loopback |
| `10.0.0.0/8` | RFC 1918 |
| `172.16.0.0/12` | RFC 1918 |
| `192.168.0.0/16` | RFC 1918 |
| `fc00::/7` | RFC 4193 (ULA) |

All other sources are refused. If your host or container network uses a different subnet, mount a custom `unbound.conf` with the appropriate `access-control` directives.

## Custom Configuration

To use a custom Unbound configuration, mount your `unbound.conf` to `/opt/unbound/etc/unbound/unbound.conf` in the container:

```sh
docker run -d \
  --name unbound \
  -p 5335:5335/udp \
  -p 5335:5335/tcp \
  -v /path/to/your/unbound.conf:/opt/unbound/etc/unbound/unbound.conf:ro \
  ghcr.io/trankimtung/unbound-container:latest
```

The example configuration from the Unbound build is available at `/opt/unbound/etc/unbound/unbound.conf.example` inside the container.

## Build

Each Unbound version has its own directory. To build locally:

```sh
docker build -t unbound:1.24.2 1.24.2/
```

Multi-arch build:

```sh
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t unbound:1.24.2 \
  1.24.2/
```

## License

This project is licensed under the BSD 3-Clause License. See [LICENSE](LICENSE.md) for details.
