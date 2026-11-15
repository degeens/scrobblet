# API Reference

This reference covers all HTTP endpoints exposed by Scrobblet for health checking, metrics, and OAuth authentication flows.

## Table of Contents
- [Endpoints](#endpoints)
  - [GET /health](#get-health)
  - [GET /metrics](#get-metrics)
  - [GET /api/spotify/login](#get-apispotifylogin)
  - [GET /api/spotify/callback](#get-apispotifycallback)
  - [GET /api/lastfm/login](#get-apilastfmlogin)
  - [GET /api/lastfm/callback](#get-apilastfmcallback)
  - [GET /api/lastfm/{name}/login](#get-apilastfmnamelogin)
  - [GET /api/lastfm/{name}/callback](#get-apilastfmnamecallback)

## Endpoints

### GET /health

Returns the health status of the source and target clients.

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | All clients are healthy |
| `503 Service Unavailable` | One or more clients are unhealthy |

##### Headers

| Header | Value |
|---|---|
| `Content-Type` | `application/json` |

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
      "name": "koito",
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
| `targets[].name` | Target instance name (e.g. `koito`, `work`) |
| `targets[].status` | Target health status. Options: `healthy`, `unhealthy` |
| `targets[].timestamp` | Time of the last target health check |


### GET /metrics

Returns Prometheus metrics, including application, Go runtime, and process metrics.

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | Metrics returned |

##### Headers

| Header | Value |
|---|---|
| `Content-Type` | `text/plain` |

##### Body

The table below documents the application metrics:

| Metric Name | Type | Help | Labels |
|---|---|---|---|
| `scrobblet_build_info` | GaugeVec | Build information. | `version` |
| `scrobblet_polling_interval_seconds` | Gauge | Current polling interval (in seconds). | — |
| `scrobblet_polls_total` | CounterVec | Total number of polls. | `source`, `status` |
| `scrobblet_submits_total` | CounterVec | Total number of submits. | `target`, `instance`, `kind`, `status` |
| `scrobblet_scrobble_evaluations_total` | CounterVec | Total number of scrobble evaluations. | `result` |
| `scrobblet_poll_duration_seconds` | HistogramVec | Duration (in seconds) of polls. | `source` |
| `scrobblet_submit_duration_seconds` | HistogramVec | Duration (in seconds) of submits. | `target`, `instance`, `kind` |

### GET /api/spotify/login

Initiates the Spotify OAuth 2.0 authorization flow by redirecting the browser to Spotify's authorization page. 

Visit this URL in a browser to authenticate Scrobblet with your Spotify account. See [Configuration Guide](configuration.md#spotify) for setup instructions.

**Only available when the source type is `Spotify`.**

#### Response

##### Status

| Status | Description |
|---|---|
| `302 Found` | Redirect to Spotify authorization page |
| `500 Internal Server Error` | Failed to redirect |

### GET /api/spotify/callback

Handles the OAuth 2.0 callback from Spotify. Validates the state parameter and exchanges the authorization code for an access token.

This endpoint is called automatically by Spotify after the user authorizes the app.

**Only available when the source type is `Spotify`.**

#### Request

##### Query Parameters

| Parameter | Required | Description |
|---|---|---|
| `code` | Yes | Authorization code provided by Spotify |
| `state` | Yes | State parameter echoed back by Spotify |

#### Response

##### Status

| Status | Description |
|---|---|
| `200 OK` | Authentication successful |
| `400 Bad Request` | Invalid `code` or `state` parameter |
| `500 Internal Server Error` | Failed to exchange the authorization code for a token |

##### Headers

| Header | Value |
|---|---|
| `Content-Type` | `text/plain` |

##### Body

On success (`200 OK`):

```
Authentication successful! Feel free to close this browser tab.
```

On failure (`400 Bad Request`, `500 Internal Server Error`):

```
Authentication failed. Please try again.
```

### GET /api/lastfm/login

Initiates the Last.fm authentication flow by redirecting the browser to the Last.fm authorization page.

Visit this URL in a browser to authenticate Scrobblet with your Last.fm account. See [Configuration Guide](configuration.md#lastfm) for setup instructions.

**Only available when a `LastFm` target with the default name (`lastfm`) is configured.** Named instances use `/api/lastfm/{name}/login` instead.

#### Response

##### Status

| Status | Description |
|---|---|
| `302 Found` | Redirect to Last.fm authorization page |
| `500 Internal Server Error` | Failed to redirect |

### GET /api/lastfm/callback

Handles the callback from Last.fm. Exchanges the authentication token for a session key.

This endpoint is called automatically by Last.fm after the user authorizes the app.

**Only available when a `LastFm` target with the default name (`lastfm`) is configured.** Named instances use `/api/lastfm/{name}/callback` instead.

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

##### Headers

| Header | Value |
|---|---|
| `Content-Type` | `text/plain` |

##### Body

On success (`200 OK`):

```
Authentication successful! Feel free to close this browser tab.
```

On failure (`400 Bad Request`, `500 Internal Server Error`):

```
Authentication failed. Please try again.
```
### GET /api/lastfm/{name}/login

Initiates the Last.fm authentication flow for a named Last.fm target instance (e.g. `second-account`), by redirecting the browser to the Last.fm authorization page.

Visit this URL in a browser to authenticate the named instance with its Last.fm account. Each instance has its own session, stored in `lastfm_session_<name>.json`.

**Only available for each `LastFm` target configured with an explicit `name`.**

#### Response

##### Status

| Status | Description |
|---|---|
| `302 Found` | Redirect to Last.fm authorization page |
| `404 Not Found` | No Last.fm target with this name configured |
| `500 Internal Server Error` | Failed to redirect |

### GET /api/lastfm/{name}/callback

Handles the callback from Last.fm for a named Last.fm target instance. Exchanges the authentication token for a session key.

This endpoint is called automatically by Last.fm after the user authorizes the app. The `redirectUrl` of a named Last.fm target must use the path `/api/lastfm/<name>/callback`.

**Only available for each `LastFm` target configured with an explicit `name`.**

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
| `404 Not Found` | No Last.fm target with this name configured |
| `500 Internal Server Error` | Failed to exchange the token for a session key |

##### Headers

| Header | Value |
|---|---|
| `Content-Type` | `text/plain` |

##### Body

On success (`200 OK`):

```
Authentication successful! Feel free to close this browser tab.
```

On failure (`400 Bad Request`, `500 Internal Server Error`):

```
Authentication failed. Please try again.
```
