# Configuration Guide

All configuration is done through environment variables.

## Table of Contents
- [General Configuration](#general-configuration)
- [Client Configuration](#client-configuration)
  - [Spotify Configuration](#spotify-configuration)
  - [Koito Configuration](#koito-configuration)
  - [ListenBrainz Configuration](#listenbrainz-configuration)
  - [Last.fm Configuration](#lastfm-configuration)
  - [CSV Configuration](#csv-configuration)

## General Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SCROBBLET_PORT` | No | `7276` | Port the API server listens on |
| `SCROBBLET_DATA_PATH` | No | `/etc/scrobblet` | Path where application data is stored |
| `SCROBBLET_SOURCE` | Yes | - | Source to track. Options: `Spotify` |
| `SCROBBLET_TARGET` | Yes | - | Target to scrobble to. Options: `Koito`, `ListenBrainz`, `LastFm`, `CSV` |

## Client Configuration

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

### CSV Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `CSV_FILE_PATH` | No | `${SCROBBLET_DATA_PATH}/scrobbles.csv` | Path to the CSV file where scrobbles will be written |

*Required only when `SCROBBLET_TARGET=CSV`

The CSV target writes completed scrobbles to a CSV file with the following format:
- **Artist(s)**: Multiple artists joined with ", "
- **Title**: Track title
- **Album**: Album name
- **Started At**: ISO 8601 timestamp when tracking started
- **Ended At**: ISO 8601 timestamp when tracking ended

The CSV file is created automatically with headers on the first scrobble. Subsequent scrobbles are appended to the file.
