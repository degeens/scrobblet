# API Reference

Scrobblet exposes a small HTTP API for health checking and authentication flows.

**Default base URL:** `http://localhost:7276`

The port is configurable via `SCROBBLET_PORT`. See [Configuration Guide](configuration.md) for more details.

## Table of Contents
- [Endpoints](#endpoints)
  - [GET /health](#get-health)

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

