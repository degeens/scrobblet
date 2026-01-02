# Scrobblet

![Go version](https://img.shields.io/github/go-mod/go-version/degeens/scrobblet?label=go)
[![Go Report](https://goreportcard.com/badge/github.com/degeens/scrobblet)](https://goreportcard.com/report/github.com/degeens/scrobblet)
![Release](https://img.shields.io/github/v/release/degeens/scrobblet?include_prereleases)
[![License](https://img.shields.io/github/license/degeens/scrobblet)](LICENSE)


Scrobblet is a lightweight scrobbler for self-hosters. It tracks your listening activity from a music source and scrobbles it to your preferred target, and is easily extensible with new integrations.

> **⚠️ Warning**: This project is in early stages of development. Features and APIs may change without notice.

## Supported Sources and Targets

Sources:
- Spotify

Targets:
- Koito
- ListenBrainz

More sources and targets can be easily added! Feel free to [create a pull request](https://github.com/degeens/scrobblet/pulls) with your implementation or [open an issue](https://github.com/degeens/scrobblet/issues) to request a new integration.

## Getting Started

Create a `docker-compose.yml` file with your configuration:

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
      - SCROBBLET_PORT=7276
      - SCROBBLET_DATA_PATH=/etc/scrobblet
      - SCROBBLET_SOURCE=Spotify
      - SCROBBLET_TARGET=Koito
      - SPOTIFY_CLIENT_ID=your_spotify_client_id
      - SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
      - SPOTIFY_REDIRECT_URL=http://127.0.0.1:7276/callback
      - KOITO_URL=your_koito_url
      - KOITO_TOKEN=your_koito_token
      - LISTENBRAINZ_TOKEN=your_listenbrainz_token

volumes:
  scrobblet-data:
```

Start the service with `docker-compose up -d`. If using Spotify as a source, visit `http://localhost:7276/login` to authenticate with your Spotify account.

## Configuration

All configuration is done through environment variables.

### Core Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SCROBBLET_PORT` | No | `7276` | Port the API server listens on |
| `SCROBBLET_DATA_PATH` | No | `/etc/scrobblet` | Path where application data is stored |
| `SCROBBLET_SOURCE` | Yes | - | Source to track (see [Supported Sources and Targets](#supported-sources-and-targets)) |
| `SCROBBLET_TARGET` | Yes | - | Target to scrobble to (see [Supported Sources and Targets](#supported-sources-and-targets)) |

### Spotify Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SPOTIFY_CLIENT_ID` | Yes* | - | Your Spotify app Client ID |
| `SPOTIFY_CLIENT_SECRET` | Yes* | - | Your Spotify app Client secret |
| `SPOTIFY_REDIRECT_URL` | Yes* | - | OAuth 2.0 redirect URL (must match your Spotify app Redirect URI) |

*Required only when `SCROBBLET_SOURCE=Spotify`

To get Spotify credentials:
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://127.0.0.1:7276/callback` (or your custom URL) to Redirect URIs
4. Copy the Client ID and Client secret

### Koito Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `KOITO_URL` | Yes* | - | URL of your Koito instance (e.g., `http://localhost:4110`) |
| `KOITO_TOKEN` | Yes* | - | Your Koito API key |

*Required only when `SCROBBLET_TARGET=Koito`

To get a Koito token:
1. Access your Koito instance
2. Log in
3. Go to Settings → API Keys
4. Generate a new API key

### ListenBrainz Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LISTENBRAINZ_TOKEN` | Yes* | - | Your ListenBrainz user token |

*Required only when `SCROBBLET_TARGET=ListenBrainz`

To get a ListenBrainz token:
1. Go to [ListenBrainz User Settings](https://listenbrainz.org/settings/)
2. Copy your User token

## License

See [LICENSE](LICENSE) file for details.
