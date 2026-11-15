# Installation Guide

This guide covers installing Scrobblet with Docker Compose using either a prebuilt image (recommended) or a local build.

## Table of Contents
- [Use a prebuilt image (recommended)](#use-a-prebuilt-image-recommended)
- [Build the image locally](#build-the-image-locally)

## Use a prebuilt image (recommended)

1. Create `scrobblet.yaml`:

```yaml
source:
  type: Spotify
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback

targets:
  - type: CSV
```

2. Create `docker-compose.yml`:

```yaml
services:
  scrobblet:
    container_name: scrobblet
    image: degeens/scrobblet:latest
    volumes:
      - scrobblet-data:/etc/scrobblet
      - ./scrobblet.yaml:/etc/scrobblet/scrobblet.yaml:ro
    ports:
      - 7276:7276
    environment:
      - SPOTIFY_CLIENT_ID=your_spotify_client_id
      - SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
    restart: unless-stopped
volumes:
  scrobblet-data:
```

3. Edit `scrobblet.yaml` and `docker-compose.yml` based on the [Configuration Guide](configuration.md). Keep secrets in `environment:` (or Docker secrets) and reference them from `scrobblet.yaml` with `${...}`.

4. Start Scrobblet:

```bash
docker compose up -d
```

5. Check the logs to verify Scrobblet is running:

```bash
docker logs scrobblet
```

## Build the image locally

1. Clone the repository and enter the directory:

```bash
git clone https://github.com/degeens/scrobblet.git
cd scrobblet
```

2. Edit the included `scrobblet.yaml` to set the desired source and targets based on the [Configuration Guide](configuration.md).

3. Create a `.env` holding the secrets referenced from `scrobblet.yaml`:

```bash
SPOTIFY_CLIENT_ID=
SPOTIFY_CLIENT_SECRET=
LISTENBRAINZ_TOKEN=
MALOJA_TOKEN=
KOITO_TOKEN=
LASTFM_API_KEY=
LASTFM_SHARED_SECRET=
```

4. Build and start Scrobblet with the desired Compose file. Example (CSV):

```bash
docker compose -f docker-compose.dev.csv.yml up -d --build
```

5. Check the logs to verify Scrobblet is running:

```bash
docker logs scrobblet
```
