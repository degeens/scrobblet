# Installation Guide

This guide covers installing Scrobblet with Docker Compose using either a prebuilt image (recommended) or a local build.

## Table of Contents
- [Use a prebuilt image (recommended)](#use-a-prebuilt-image-recommended)
- [Build the image locally](#build-the-image-locally)

## Use a prebuilt image (recommended)

1. Create `docker-compose.yml`:

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

2. Edit `docker-compose.yml` to set the required environment variables based on the [Configuration Guide](configuration.md).

3. Start Scrobblet:

```bash
docker compose up -d
```

4. Check the logs to verify Scrobblet is running:

```bash
docker logs scrobblet
```

## Build the image locally

1. Clone the repository and enter the directory:

```bash
git clone https://github.com/degeens/scrobblet.git
cd scrobblet
```

2. Copy the example environment file:

```bash
cp example.env .env
```

3. Edit `.env` to set the required environment variables based on the [Configuration Guide](configuration.md).

4. Build and start Scrobblet with the desired Compose file. Example (CSV):

```bash
docker compose -f docker-compose.dev.csv.yml up -d --build
```

5. Check the logs to verify Scrobblet is running:

```bash
docker logs scrobblet
```