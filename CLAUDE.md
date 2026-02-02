# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Scrobblet is a lightweight music scrobbler for self-hosters written in Go. It tracks listening activity from music sources (currently Spotify) and scrobbles to targets (Koito, ListenBrainz). The project is designed to be easily extensible with new source and target integrations.

## Build and Development Commands

### Building
```bash
# Build the application
go build -o scrobblet ./cmd/api

# Build Docker image
docker build -t scrobblet .
```

### Running
```bash
# Run locally (requires environment variables to be set)
go run ./cmd/api

# Run with Docker Compose
docker-compose up -d
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Testing
The project currently does not have a test suite set up.

## Architecture

### High-Level Design

Scrobblet uses a **producer-consumer architecture** with goroutines and channels:

1. **Main Application** (`cmd/api/main.go`): Initializes configuration, sets up source/target clients, starts the scrobbler in a goroutine, and runs an HTTP server for OAuth callbacks
2. **Scrobbler** (`internal/scrobbler/scrobbler.go`): Coordinates two concurrent workers via channels
3. **Tracker** (`internal/scrobbler/tracker.go`): Polls the music source, detects track changes/replays, tracks listening duration, and sends tracks to channels
4. **Submitter** (`internal/scrobbler/submitter.go`): Consumes from channels and submits "now playing" and "played" tracks to the target

### Channel-Based Communication

The scrobbler uses two channels for communication:
- `playingTrackChan` (buffer size 1): Carries the currently playing track for "now playing" updates
- `playedTrackChan` (buffer size 10): Queues completed tracks that meet the scrobble threshold

### Core Components

#### Sources (`internal/sources/`)
- **Source Interface**: Defines `GetPlaybackState()` method
- **SpotifySource**: Implements polling of Spotify's Currently Playing API
- Sources return `PlaybackState` containing track info, position, and timestamp

#### Targets (`internal/targets/`)
- **Target Interface**: Defines `SubmitPlayingTrack()` and `SubmitPlayedTrack()` methods
- **KoitoTarget**: Submits to Koito API
- **ListenBrainzTarget**: Submits to ListenBrainz API

#### Clients (`internal/clients/`)
- **Spotify Client** (`spotify/`): Handles OAuth 2.0 flow, token refresh, and API requests
- **Koito Client** (`koito/`): HTTP client for Koito API
- **ListenBrainz Client** (`listenbrainz/`): HTTP client for ListenBrainz API

#### Common Models (`internal/common/`)
- `Track`: Represents a music track (artists, title, album, duration)
- `TrackedTrack`: Extends Track with tracking metadata (duration listened, timestamps, last position)
- Scrobble threshold logic (track must be played for 50% of duration or 4 minutes, whichever comes first)

### Tracker Behavior

The Tracker implements intelligent polling:
- **Active Polling**: 10s interval when music is playing
- **Inactive Polling**: 30s interval after 5 minutes of inactivity
- **Drift Tolerance**: 2s tolerance for position drift to detect seeks/pauses
- **Track Change Detection**: Detects track changes or replays (track restarted from beginning)

### HTTP API

The API server (`cmd/api/`) provides OAuth endpoints when using Spotify:
- `GET /login`: Initiates Spotify OAuth flow
- `GET /callback`: Handles OAuth callback and stores tokens

## Configuration

All configuration is done via environment variables (see README.md for full list). Key variables:
- `SCROBBLET_SOURCE`: Source type (e.g., "Spotify")
- `SCROBBLET_TARGETS`: Target types (e.g., "Koito", "ListenBrainz")
- `SCROBBLET_DATA_PATH`: Where tokens/data are persisted (default: `/etc/scrobblet`)

## Adding New Integrations

### Adding a New Source
1. Create client in `internal/clients/<source>/`
2. Implement `Source` interface in `internal/sources/<source>_source.go`
3. Add source type constant to `internal/sources/source.go`
4. Update factory function in `sources.New()`
5. Add configuration loading in `cmd/api/config.go`

### Adding a New Target
1. Create client in `internal/clients/<target>/`
2. Implement `Target` interface in `internal/targets/<target>_target.go`
3. Add target type constant to `internal/targets/target.go`
4. Update factory function in `targets.New()`
5. Add configuration loading in `cmd/api/config.go`

## OAuth Token Management

Spotify tokens are stored in `SCROBBLET_DATA_PATH/spotify_token.json`. The `spotify.Client` automatically handles token refresh using the `golang.org/x/oauth2` package.
