# Scrobblet

![Go version](https://img.shields.io/github/go-mod/go-version/degeens/scrobblet?label=go)
[![Go Report](https://goreportcard.com/badge/github.com/degeens/scrobblet)](https://goreportcard.com/report/github.com/degeens/scrobblet)
[![Release](https://img.shields.io/github/v/release/degeens/scrobblet?include_prereleases)](https://github.com/degeens/scrobblet/releases)
[![License](https://img.shields.io/github/license/degeens/scrobblet)](LICENSE)


Scrobblet is a lightweight scrobbler for self-hosters. It tracks your listening activity from a music source and scrobbles it to your preferred target, and is easily extensible with new integrations.

> **⚠️ Warning**: This project is in early stages of development. Features and APIs may change without notice.

## Supported Sources and Targets

Sources:
- Spotify

Targets:
- Koito
- ListenBrainz
- Last.fm

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
      # Spotify (Required when SCROBBLET_SOURCE=Spotify)
      - SPOTIFY_CLIENT_ID=your_spotify_client_id
      - SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
      - SPOTIFY_REDIRECT_URL=http://127.0.0.1:7276/spotify/callback
      # Koito (Required when SCROBBLET_TARGET=Koito)
      - KOITO_URL=your_koito_url
      - KOITO_TOKEN=your_koito_token
      # ListenBrainz (Required when SCROBBLET_TARGET=ListenBrainz)
      - LISTENBRAINZ_TOKEN=your_listenbrainz_token
      # Last.fm (Required when SCROBBLET_TARGET=LastFm)
      - LASTFM_API_KEY=your_lastfm_api_key
      - LASTFM_SHARED_SECRET=your_lastfm_shared_secret
      - LASTFM_REDIRECT_URL=http://127.0.0.1:7276/lastfm/callback
    restart: unless-stopped
volumes:
  scrobblet-data:
```

Start the service with `docker-compose up -d`.

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

To set up Spotify:
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://127.0.0.1:7276/spotify/callback` (or your custom URL) to Redirect URIs
4. Copy the client ID and client secret
5. Start Scrobblet with the client ID and client secret configured
6. Visit `http://localhost:7276/spotify/login` to authenticate

### Koito Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `KOITO_URL` | Yes* | - | URL of your Koito instance (e.g., `http://localhost:4110`) |
| `KOITO_TOKEN` | Yes* | - | Your Koito API key |

*Required only when `SCROBBLET_TARGET=Koito`

To set up Koito:
1. Access your Koito instance
2. Log in
3. Go to Settings → API Keys
4. Generate a new API key
5. Start Scrobblet with the API key configured

### ListenBrainz Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LISTENBRAINZ_TOKEN` | Yes* | - | Your ListenBrainz user token |

*Required only when `SCROBBLET_TARGET=ListenBrainz`

To set up ListenBrainz:
1. Go to [ListenBrainz User Settings](https://listenbrainz.org/settings/)
2. Copy your user token
3. Start Scrobblet with the user token configured

### Last.fm Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LASTFM_API_KEY` | Yes* | - | Your Last.fm API key |
| `LASTFM_SHARED_SECRET` | Yes* | - | Your Last.fm shared secret |
| `LASTFM_REDIRECT_URL` | Yes* | - | Redirect URL (must match your Last.fm API account callback URL) |

*Required only when `SCROBBLET_TARGET=LastFm`

To set up Last.fm:
1. Go to [Last.fm API Account Creation](https://www.last.fm/api/account/create)
2. Create an API account
3. Set callback URL to `http://127.0.0.1:7276/lastfm/callback` (or your custom URL)
4. Copy the API key and shared secret
5. Start Scrobblet with the API key and shared secret configured
6. Visit `http://localhost:7276/lastfm/login` to authenticate

## License

See [LICENSE](LICENSE) file for details.
