<p align="center">
  <img src="logo.png" alt="Scrobbling gopher" height="275"/>
</p>

# Scrobblet

[![CI](https://github.com/degeens/scrobblet/actions/workflows/ci.yml/badge.svg)](https://github.com/degeens/scrobblet/actions/workflows/ci.yml)
![Go version](https://img.shields.io/github/go-mod/go-version/degeens/scrobblet?label=go)
[![Go Report](https://goreportcard.com/badge/github.com/degeens/scrobblet)](https://goreportcard.com/report/github.com/degeens/scrobblet)
[![Release](https://img.shields.io/github/v/release/degeens/scrobblet?include_prereleases)](https://github.com/degeens/scrobblet/releases)
[![License](https://img.shields.io/github/license/degeens/scrobblet)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/degeens/scrobblet)](https://hub.docker.com/r/degeens/scrobblet)

Scrobblet is a lightweight scrobbler for self-hosters. It tracks your listening activity from a music source and scrobbles it to your preferred targets, and is easily extensible with new integrations.

Currently, Scrobblet supports **Spotify** as a source and can scrobble to **Last.fm**, **ListenBrainz** (including any ListenBrainz-compatible service), **Maloja**, **Koito**, and **CSV**. More sources and targets can be easily added. Feel free to [create a pull request](https://github.com/degeens/scrobblet/pulls) with your implementation or [open an issue](https://github.com/degeens/scrobblet/issues) to request a new integration.

> **⚠️ Warning**: This project is in early stages of development. Features and APIs may change without notice.

## Getting Started

Here's a minimal `docker-compose.yml` file to get started with scrobbling from Spotify to CSV:

```yaml
services:
  scrobblet:
    container_name: scrobblet
    image: degeens/scrobblet:latest
    volumes:
      - scrobblet-data:/etc/scrobblet
    ports:
      - 7276:7276
    environment:
      - SCROBBLET_SOURCE=Spotify
      - SCROBBLET_TARGETS=CSV
      - SPOTIFY_CLIENT_ID=your_spotify_client_id
      - SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
      - SPOTIFY_REDIRECT_URL=http://127.0.0.1:7276/spotify/callback
    restart: unless-stopped
volumes:
  scrobblet-data:
```

Set `SPOTIFY_CLIENT_ID`, `SPOTIFY_CLIENT_SECRET`, and `SPOTIFY_REDIRECT_URL` with your actual Spotify application credentials. See the [Configuration Guide](docs/configuration.md) for instructions on obtaining these.

Start the service with `docker-compose up -d`.

## Configuration

For detailed configuration instructions, see the [Configuration Guide](docs/configuration.md).

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
