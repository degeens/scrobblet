# Configuration Guide

This guide covers all available Scrobblet configuration options, which are set through a YAML configuration file.

## Table of Contents
- [Configuration File](#configuration-file)
- [General Configuration](#general-configuration)
- [Source Configuration](#source-configuration)
  - [Spotify](#spotify)
  - [WiiM](#wiim)
- [Target Configuration](#target-configuration)
  - [Last.fm](#lastfm)
  - [ListenBrainz](#listenbrainz)
  - [Maloja](#maloja)
  - [Koito](#koito)
  - [CSV](#csv)
- [Full Example](#full-example)
- [Migrating from Environment Variables](#migrating-from-environment-variables)

## Configuration File

Scrobblet reads a single YAML configuration file:

- The path is taken from the `SCROBBLET_CONFIG` environment variable, defaulting to `/etc/scrobblet/scrobblet.yaml`.
- `${ENV_VAR}` references in the file are expanded from the environment before parsing, so secrets can stay in environment variables or Docker secrets. An unset variable expands to an empty string and is reported as a missing required field.
- Unknown keys are rejected to catch typos.
- Restart Scrobblet to apply configuration changes.

`SCROBBLET_CONFIG` is the only environment variable read directly; all other configuration lives in the YAML file.

## General Configuration

| Key | Required | Default | Description |
|-----|----------|---------|-------------|
| `port` | No | `7276` | Port the API server listens on |
| `dataPath` | No | `/etc/scrobblet` | Path where application data is stored |
| `logLevel` | No | `INFO` | The log level. Options: `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `rateLimit.rate` | No | `10` | Maximum requests per second |
| `rateLimit.burst` | No | `100` | Maximum burst size |
| `source` | Yes | — | Source to track (single object, see below) |
| `targets` | Yes | — | List of targets to scrobble to (see below) |

## Source Configuration

Exactly one source is configured via the `source` key.

### Spotify

> **ℹ️ Spotify Premium account required**: Scrobblet uses Spotify's Web API to fetch your currently playing track. Since February 11, 2026, all newly created Development Mode apps require a Spotify Premium account. As of March 9, 2026, this requirement also applies to all existing apps. See the [Spotify Developer Blog](https://developer.spotify.com/blog/2026-02-06-update-on-developer-access-and-platform-security) for full details.

```yaml
source:
  type: Spotify
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `clientId` | Yes | — | Your Spotify app Client ID |
| `clientSecret` | Yes | — | Your Spotify app Client secret |
| `redirectUrl` | Yes | — | OAuth 2.0 redirect URL, path must be `/api/spotify/callback` (must match your Spotify app Redirect URI) |

To set up Spotify:
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://127.0.0.1:7276/api/spotify/callback` (or your custom URL) to Redirect URIs
4. Copy the client ID and client secret
5. Start Scrobblet with the client ID, client secret, and redirect URL configured
6. Visit `http://localhost:7276/api/spotify/login` to authenticate

> **⚠️ Spotify Redirect URI**: Spotify only allows **HTTPS** Redirect URIs unless you use an explicit IPv4 or IPv6 loopback address (e.g. `http://127.0.0.1:7276` or `http://[::1]:7276`).
>
> If you cannot provide a valid redirect URI, you can intercept the callback and rewrite the address to your Scrobblet instance:
> 1. Run `docker run --rm -p 8080:8080 mendhak/http-https-echo:39` on the machine where your browser runs
> 2. Set the Spotify Redirect URI and `source.redirectUrl` to `http://127.0.0.1:8080/api/spotify/callback`
> 3. Visit `http://localhost:7276/api/spotify/login` to authenticate
> 4. Replace `127.0.0.1:8080` with the address of your Scrobblet instance and open the updated URL

### WiiM

```yaml
source:
  type: WiiM
  url: https://192.168.0.100
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `url` | Yes | — | Your WiiM device URL (e.g., `https://192.168.0.100`). Must be a local IP address |

To set up WiiM:
1. Open the WiiM Home mobile app
2. Open the settings for your WiiM device
3. Go to General → Network and copy the IP address
4. Start Scrobblet with `source.url` set to `https://` + the copied IP address (e.g., `https://192.168.0.100`)

> **ℹ️ Static IP recommended**: Since Scrobblet connects to your WiiM device using its IP address, it is recommended to configure a static IP for your WiiM device in the WiiM Home app so the address doesn't change over time.

## Target Configuration

Targets are configured as a **list** under the `targets` key. Every target entry has:

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `type` | Yes | — | `LastFm`, `ListenBrainz`, `Maloja`, `Koito` or `CSV` |
| `name` | No | Lowercase type (e.g. `listenbrainz`) | Instance name. Must match `[a-z0-9][a-z0-9-]*` |

Multiple targets of the same type are allowed as long as each instance has a unique name — e.g. scrobble to both the official ListenBrainz service and a self-hosted ListenBrainz-compatible service, or to two CSV files. The name identifies the instance in logs, the `/health` response, and the `instance` label of the submit metrics.

### Last.fm

```yaml
targets:
  - type: LastFm
    apiKey: ${LASTFM_API_KEY}
    sharedSecret: ${LASTFM_SHARED_SECRET}
    redirectUrl: http://127.0.0.1:7276/api/lastfm/callback
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `apiKey` | Yes | — | Your Last.fm API key |
| `sharedSecret` | Yes | — | Your Last.fm shared secret |
| `redirectUrl` | Yes | — | Redirect URL (must match your Last.fm API account callback URL). Path must be `/api/lastfm/callback` for the default name, `/api/lastfm/<name>/callback` for named instances |

To set up Last.fm:
1. Go to [Last.fm API Account Creation](https://www.last.fm/api/account/create)
2. Create an API account
3. Set callback URL to `http://127.0.0.1:7276/api/lastfm/callback` (or your custom URL)
4. Copy the API key and shared secret
5. Start Scrobblet with the API key and shared secret configured
6. Visit `http://localhost:7276/api/lastfm/login` to authenticate

Each Last.fm instance has its own login flow and session file:
- Default name (`lastfm`): log in via `/api/lastfm/login`, session stored in `lastfm_session.json`
- Named instance (e.g. `second-account`): log in via `/api/lastfm/second-account/login`, session stored in `lastfm_session_second-account.json`

### ListenBrainz

```yaml
targets:
  - type: ListenBrainz
    token: ${LISTENBRAINZ_TOKEN}
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `token` | Yes | — | Your ListenBrainz user token |
| `url` | No | `https://api.listenbrainz.org` | URL of a ListenBrainz-compatible service |

To set up ListenBrainz:
1. Go to [ListenBrainz User Settings](https://listenbrainz.org/settings/)
2. Copy your user token
3. Start Scrobblet with the user token configured

> ✅ **ListenBrainz-compatible API support**: Scrobblet can scrobble to any service that implements the ListenBrainz API specification. This includes the official ListenBrainz service, self-hosted ListenBrainz instances, or any other compatible implementation. Simply set `url` to point to your preferred service. Use a named second instance to scrobble to both.

### Maloja

```yaml
targets:
  - type: Maloja
    url: http://maloja:42010
    token: ${MALOJA_TOKEN}
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `url` | Yes | — | URL of your Maloja instance (e.g., `http://localhost:42010`) |
| `token` | Yes | — | Your Maloja API key |

Maloja speaks the ListenBrainz-compatible API; the path is forced to `/apis/listenbrainz`. Equivalently, you can use `type: ListenBrainz` with `url: http://maloja:42010/apis/listenbrainz`.

To set up Maloja:
1. Access your Maloja instance
2. Go to Settings → API Keys
3. Generate a new API key
4. Start Scrobblet with the API key configured

### Koito

```yaml
targets:
  - type: Koito
    url: http://koito:4110
    token: ${KOITO_TOKEN}
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `url` | Yes | — | URL of your Koito instance (e.g., `http://localhost:4110`) |
| `token` | Yes | — | Your Koito API key |

Koito speaks the ListenBrainz-compatible API; the path is forced to `/apis/listenbrainz`. Equivalently, you can use `type: ListenBrainz` with `url: http://koito:4110/apis/listenbrainz`.

To set up Koito:
1. Access your Koito instance
2. Log in
3. Go to Settings → API Keys
4. Generate a new API key
5. Start Scrobblet with the API key configured

### CSV

```yaml
targets:
  - type: CSV
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `filePath` | No | `<dataPath>/scrobbles.csv` for the default name, `<dataPath>/scrobbles-<name>.csv` for named instances | Path to the CSV file where scrobbles will be written |

Two CSV instances may not write to the same file.

The CSV target writes completed scrobbles to a CSV file with the following format:
- **Artists**: Multiple artists joined with ", "
- **Title**: Track title
- **Album**: Album name
- **Track Number**: Track number on the album
- **Duration**: Track duration in seconds
- **Started At**: ISO 8601 timestamp when tracking started
- **Ended At**: ISO 8601 timestamp when tracking ended

The CSV file is created automatically with headers on the first scrobble. Subsequent scrobbles are appended to the file.

## Full Example

```yaml
port: 7276
dataPath: /etc/scrobblet
logLevel: INFO
rateLimit:
  rate: 10
  burst: 100

source:
  type: Spotify
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback

targets:
  - type: ListenBrainz
    token: ${LISTENBRAINZ_TOKEN}

  - type: ListenBrainz
    name: work
    url: https://listenbrainz.example.com
    token: ${LISTENBRAINZ_WORK_TOKEN}

  - type: Maloja
    url: http://maloja:42010
    token: ${MALOJA_TOKEN}

  - type: Koito
    url: http://koito:4110
    token: ${KOITO_TOKEN}

  - type: LastFm
    apiKey: ${LASTFM_API_KEY}
    sharedSecret: ${LASTFM_SHARED_SECRET}
    redirectUrl: http://127.0.0.1:7276/api/lastfm/callback

  - type: LastFm
    name: second-account
    apiKey: ${LASTFM_2_API_KEY}
    sharedSecret: ${LASTFM_2_SHARED_SECRET}
    redirectUrl: http://127.0.0.1:7276/api/lastfm/second-account/callback

  - type: CSV

  - type: CSV
    name: archive
```

## Migrating from Environment Variables

Previous versions of Scrobblet were configured entirely through environment variables. Create a `scrobblet.yaml` and map your settings as follows. Secrets can keep living in the same environment variables, referenced with `${...}`.

| Environment variable | YAML path |
|----------------------|-----------|
| `SCROBBLET_PORT` | `port` |
| `SCROBBLET_DATA_PATH` | `dataPath` |
| `SCROBBLET_LOG_LEVEL` | `logLevel` |
| `SCROBBLET_RATE_LIMIT_RATE` | `rateLimit.rate` |
| `SCROBBLET_RATE_LIMIT_BURST` | `rateLimit.burst` |
| `SCROBBLET_SOURCE` | `source.type` |
| `SCROBBLET_TARGETS` | `targets[].type` (list entries) |
| `SPOTIFY_CLIENT_ID` | `source.clientId` |
| `SPOTIFY_CLIENT_SECRET` | `source.clientSecret` |
| `SPOTIFY_REDIRECT_URL` | `source.redirectUrl` |
| `WIIM_URL` | `source.url` |
| `LISTENBRAINZ_URL` | `targets[].url` |
| `LISTENBRAINZ_TOKEN` | `targets[].token` |
| `MALOJA_URL` | `targets[].url` |
| `MALOJA_TOKEN` | `targets[].token` |
| `KOITO_URL` | `targets[].url` |
| `KOITO_TOKEN` | `targets[].token` |
| `LASTFM_API_KEY` | `targets[].apiKey` |
| `LASTFM_SHARED_SECRET` | `targets[].sharedSecret` |
| `LASTFM_REDIRECT_URL` | `targets[].redirectUrl` |
| `CSV_FILE_PATH` | `targets[].filePath` |

Note for Last.fm users: an unnamed (`lastfm`) Last.fm target keeps the existing `/api/lastfm/login` and `/api/lastfm/callback` URLs and the existing `lastfm_session.json` session file, so you do not need to log in again.
