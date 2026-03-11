# API Reference

Scrobblet exposes a small HTTP API for health checking and authentication flows.

**Default base URL:** `http://localhost:7276`

The port is configurable via `SCROBBLET_PORT`. See [Configuration Guide](configuration.md) for more details.

## Table of Contents
- [Endpoints](#endpoints)
  - [GET /health](#get-health)
  - [GET /spotify/login](#get-spotifylogin)
  - [GET /spotify/callback](#get-spotifycallback)
  - [GET /lastfm/login](#get-lastfmlogin)
  - [GET /lastfm/callback](#get-lastfmcallback)

## Endpoints

### GET /health

Returns the health status of the source and target clients.

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | All clients are healthy |
| `503 Service Unavailable` | One or more clients are unhealthy |

##### Body

```json
{
  "status": "healthy",
  "source": {
    "type": "Spotify",
    "status": "healthy",
    "timestamp": "2026-03-11T18:02:23.933059505Z"
  },
  "targets": [
    {
      "type": "Koito",
      "status": "healthy",
      "timestamp": "2026-03-11T18:01:03.990940128Z"
    }
  ]
}
```

| Field  | Description |
|---|---|
| `status` | Overall health status. `healthy` when all clients are healthy, otherwise `unhealthy`. |
| `source.type` | Source type  (e.g. `Spotify`) |
| `source.status` | Source health status. Options: `healthy`, `unhealthy` |
| `source.timestamp`  | Time of the last source health check |
| `targets[].type` | Target type  (e.g. `LastFm`, `ListenBrainz`, `Maloja`, `Koito`, `CSV`) |
| `targets[].status` | Target health status. Options: `healthy`, `unhealthy` |
| `targets[].timestamp` | Time of the last target health check |


### GET /spotify/login

Initiates the Spotify OAuth 2.0 authorization flow by redirecting the browser to Spotify's authorization page. 

Visit this URL in a browser to authenticate Scrobblet with your Spotify account. See [Configuration Guide](configuration.md#spotify) for setup instructions.

**Only available when `SCROBBLET_SOURCE=Spotify`.**

#### Response

##### Body

| Status | Description |
|---|---|
| `302 Found` | Redirect to Spotify authorization page |
| `500 Internal Server Error` | Failed to redirect |

### GET /spotify/callback

Handles the OAuth 2.0 callback from Spotify. Validates the state parameter and exchanges the authorization code for an access token.

This endpoint is called automatically by Spotify after the user authorizes the app.

**Only available when `SCROBBLET_SOURCE=Spotify`.**

#### Request

##### Query Parameters

| Parameter | Required | Description |
|---|---|---|
| `code` | Yes | Authorization code provided by Spotify |
| `state` | Yes | State parameter echoed back by Spotify (validated against stored value) |

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | Authentication successful |
| `400 Bad Request` | Invalid `code` or `state` parameter |
| `500 Internal Server Error` | Failed to exchange the authorization code for a token |

##### Body

```
Authentication successful! Feel free to close this browser tab.
```

### GET /lastfm/login

Initiates the Last.fm authentication flow by redirecting the browser to the Last.fm authorization page.

Visit this URL in a browser to authenticate Scrobblet with your Last.fm account. See [Configuration Guide](configuration.md#lastfm) for setup instructions.

**Only available when `SCROBBLET_TARGETS` includes `LastFm`.**

#### Response

##### Status

| Status | Description |
|---|---|
| `302 Found` | Redirect to Last.fm authorization page |
| `500 Internal Server Error` | Failed to redirect |

### GET /lastfm/callback

Handles the callback from Last.fm. Exchanges the authentication token for a session key.

This endpoint is called automatically by Last.fm after the user authorizes the app.

**Only available when `SCROBBLET_TARGETS` includes `LastFm`.**

#### Request

##### Query Parameters

| Parameter | Required | Description |
|---|---|---|
| `token` | Yes | Authorization token provided by Last.fm |

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | Authentication successful |
| `400 Bad Request` | Invalid `token` parameter |
| `500 Internal Server Error` | Failed to exchange the token for a session key |

##### Body

```
Authentication successful! Feel free to close this browser tab.
```