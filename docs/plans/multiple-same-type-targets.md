# Plan: YAML Configuration & Multiple Same-Type Targets

**Status:** Implemented
**Date:** 2026-10-06

## 1. Goal

Allow configuring multiple scrobble targets of the same type — e.g. scrobble to both the official ListenBrainz service and a self-hosted ListenBrainz-compatible service, or to two CSV files — by replacing the current environment-variable-based configuration with a YAML configuration file where targets are a **list of named instances** instead of a fixed set of types.

## 2. Background: why this isn't possible today

The current design conflates *target type* with *target instance* in five places:

| # | Location | Problem |
|---|----------|---------|
| 1 | `internal/config/validators.go` (`validateTargets`) | Explicitly rejects duplicate target types |
| 2 | Env var naming (`internal/config/constants.go`) | `LISTENBRAINZ_TOKEN`, `KOITO_URL`, … are type-scoped; no namespace for a second instance |
| 3 | `internal/clients/config.go` | `clients.Config` has exactly one slot per type (`Koito`, `Maloja`, `ListenBrainz`, `LastFm`, `CSV`) |
| 4 | Observability (`internal/scrobbler/submitter.go`, `internal/api/handlers/health.go`) | `TargetType()` is the only identity: Prometheus `target` label, health JSON `type` field, and log lines would silently mix instances of the same type |
| 5 | Last.fm auth (`internal/api/router.go`, `internal/clients/lastfm`) | Hardcoded `/api/lastfm/login\|callback` routes bound to the *first* Last.fm client; fixed session file `lastfm_session.json`; redirect URL validated against a fixed path |

Note: Koito and Maloja are already ListenBrainz-compatible services modeled as separate *types* sharing one client — evidence the type-per-service approach doesn't scale.

## 3. Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Config mechanism | Single YAML file, **replacing** env-var config (hard cutover) | Project is early-stage (README warns APIs may change); dual env+YAML support doubles complexity for little gain |
| Config file location | `SCROBBLET_CONFIG` env var, default `/etc/scrobblet/scrobblet.yaml` | Default data path is already a mounted volume in every Compose example; `SCROBBLET_CONFIG` is the *only* remaining env var read directly |
| Secrets | `${ENV_VAR}` expansion applied once to the raw file bytes before parsing (`os.ExpandEnv`) | Tokens stay in env vars / Docker secrets; single-pass expansion means `$` inside env var *values* is safe |
| YAML library | `gopkg.in/yaml.v3` (new dependency) | De-facto standard, supports custom unmarshalling for type-dispatched target configs |
| Instance identity | Every target entry has optional `name`; defaults to lowercase type (`listenbrainz`, `csv`, …) | Names are always non-empty downstream → no special-casing in metrics/health/logs; a duplicate type+name is a config error |
| Name charset | `[a-z0-9][a-z0-9-]*` | Used in URL paths (Last.fm routes) and file names (session/CSV files) |
| Metrics | Add separate `instance` label to `scrobblet_submits_total` and `scrobblet_submit_duration_seconds` | User decision; **breaking metrics change** (existing series gain a label), called out in release notes |
| Sources | Stay singleton (`source:` is a single object, not a list) | Out of scope; validator already enforces one source |
| Koito/Maloja behavior | Unchanged | They remain convenience presets that force the `/apis/listenbrainz` path; users can now also express them as `type: ListenBrainz` with a custom `url` |

## 4. Configuration schema

### Minimal example (equivalent to today's quick start)

```yaml
source:
  type: Spotify
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback

targets:
  - type: CSV
```

### Full example (the motivating use case)

```yaml
port: 7276                    # optional, default 7276
dataPath: /etc/scrobblet      # optional, default /etc/scrobblet
logLevel: INFO                # optional, DEBUG|INFO|WARN|ERROR, default INFO
rateLimit:                    # optional
  rate: 10                    #   default 10
  burst: 100                  #   default 100

source:
  type: Spotify               # Spotify | WiiM
  # Spotify fields (required when type=Spotify):
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback
  # WiiM field (required when type=WiiM):
  # url: https://192.168.0.100

targets:
  - type: ListenBrainz        # name defaults to "listenbrainz"
    token: ${LISTENBRAINZ_TOKEN}

  - type: ListenBrainz
    name: work                # second instance of the same type
    url: https://listenbrainz.example.com   # optional, default https://api.listenbrainz.org
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

  - type: CSV                 # filePath optional, default <dataPath>/scrobbles.csv
  - type: CSV
    name: archive             # filePath optional, default <dataPath>/scrobbles-archive.csv
```

### Per-type field reference

| Type | Fields | Required | Optional (default) |
|------|--------|----------|--------------------|
| `ListenBrainz` | `url`, `token` | `token` | `url` (`https://api.listenbrainz.org`) |
| `Maloja` | `url`, `token` | both | — (path forced to `/apis/listenbrainz`) |
| `Koito` | `url`, `token` | both | — (path forced to `/apis/listenbrainz`) |
| `LastFm` | `apiKey`, `sharedSecret`, `redirectUrl` | all | — |
| `CSV` | `filePath` | — | `filePath` (`<dataPath>/scrobbles[-<name>].csv`) |
| `Spotify` (source) | `clientId`, `clientSecret`, `redirectUrl` | all | — |
| `WiiM` (source) | `url` | yes | — |

## 5. Implementation

### Phase 1 — YAML config foundation (breaking change)

**`internal/config`** (largest change):

- `models.go` — new YAML model types with `yaml` tags mirroring the schema above, kept *separate* from client config structs. Custom `UnmarshalYAML` on the target/source entries dispatches on `type` into typed per-type option structs.
- `config.go` — `LoadConfig()` rewritten:
  1. Resolve path: `SCROBBLET_CONFIG` or `<default>/scrobblet.yaml`
  2. Read file, apply `os.ExpandEnv` once, `yaml.Unmarshal`
  3. Apply defaults (port, dataPath, logLevel, rate limit, ListenBrainz URL, CSV filePath, target names)
  4. Validate (below)
  5. Map YAML models → construction specs
- `validators.go` — ported to validate loaded structs instead of strings; env-var helpers in `utils.go` deleted (`setURLPath` stays).
- `constants.go` — reduced to defaults + `SCROBBLET_CONFIG` key; all per-integration env keys deleted.

**Dependency graph fix** — delete `internal/clients/config.go` (the monolithic one-slot-per-type struct). Define construction specs where they're consumed, avoiding import cycles (config imports sources/targets; they must not import config):

```go
// internal/sources/source.go
type Spec struct {
    Type    SourceType
    Spotify *spotify.Config   // set iff Type == SourceSpotify
    WiiM    *wiim.Config      // set iff Type == SourceWiiM
}
func New(spec Spec) (any, Source, error)

// internal/targets/target.go
type Spec struct {
    Type         TargetType
    Name         string
    ListenBrainz *listenbrainz.Config
    LastFm       *lastfm.Config
    CSV          *csv.Config
}
func New(spec Spec, scrobbletVersion string) (Instance, error)

type Instance struct {         // replaces today's parallel ([]any, []Target) returns
    Name   string
    Type   TargetType
    Client any
    Target Target
}
func NewMultiple(specs []Spec, scrobbletVersion string) ([]Instance, error)
```

**`cmd/scrobblet/main.go`** — pass `cfg.Source` / `cfg.Targets` specs; derive `[]targets.Target` and client list from `[]targets.Instance`.

### Phase 2 — Instance identity in observability

- `internal/targets/target.go` — add `Name() string` to the `Target` interface; implement in all three target wrappers (name stored at construction, already defaulted during config load).
- `internal/metrics/metrics.go` — add `instance` label to `Submits` and `SubmitDuration` vecs.
- `internal/scrobbler/submitter.go` — pass `target.Name()` into the new label and add `"instance"` to log attrs.
- `internal/health/health.go` — add `Name` to `TargetHealthCheck`.
- `internal/api/handlers/health.go` — add `name` to `ClientHealth` JSON for targets (source unchanged).

### Phase 3 — Per-instance auth flows & storage

**Last.fm** (the only target with an OAuth flow):

- `internal/clients/lastfm/client.go` — `NewClient` takes the instance name; session file becomes `lastfm_session.json` for the default name, `lastfm_session_<name>.json` otherwise.
- `internal/api/router.go` — replace `getLastFmClient([]any)` scan with iteration over `[]targets.Instance`; for each Last.fm instance register:
  - default name → `GET /api/lastfm/login`, `GET /api/lastfm/callback` (unchanged paths)
  - named → `GET /api/lastfm/{name}/login`, `GET /api/lastfm/{name}/callback`
- `internal/config/validators.go` — `validateRedirectURL` for Last.fm checks the path matching the instance (`/api/lastfm/callback` vs `/api/lastfm/<name>/callback`).

**CSV** — config validation rejects two CSV instances resolving to the same `filePath`.

### Phase 4 — Docs, examples, tooling

- `docs/configuration.md` — rewrite as YAML schema reference + env-to-YAML migration table (`LISTENBRAINZ_TOKEN` → `targets[].token` with `${LISTENBRAINZ_TOKEN}`, etc.).
- `docs/installation.md`, `README.md` — quick start mounts a config file; secrets stay in `environment:`.
- `example.env` → replace with example `scrobblet.yaml` (plus a slim `.env` holding only secrets referenced via `${...}`).
- `docker-compose.dev.*.yml` — mount the example config: `./scrobblet.yaml:/etc/scrobblet/scrobblet.yaml:ro`.
- `docs/api.md` — health response gains `targets[].name`; document named Last.fm routes; metrics tables gain the `instance` label.
- `Dockerfile` — no change needed (config arrives via the existing `/etc/scrobblet` volume).

## 6. Validation rules & edge cases

| Rule | Error behavior |
|------|----------------|
| Config file missing/unreadable | Fatal at startup with path in message |
| Unknown `type` in source/targets | Fatal, listing valid types |
| Duplicate `(type, name)` pair | Fatal naming both entries |
| Invalid `name` charset | Fatal with charset explanation |
| `${VAR}` referencing unset env var | Expands to empty → caught by required-field validation, error names the field |
| Required per-type field missing | Fatal naming target (`targets[2] (ListenBrainz "work"): token is required`) |
| Spotify/LastFm `redirectUrl` | Existing scheme/host/path validation; Last.fm path depends on instance name |
| WiiM `url` | Existing private-IP-only validation (unchanged) |
| Two CSV instances, same resolved `filePath` | Fatal |
| Unknown YAML keys | Reject via `yaml.Decoder.KnownFields(true)` to catch typos |
| `logLevel` unparseable | Keep current behavior: warn and fall back to INFO |

## 7. Testing plan

- `internal/config` — new table-driven tests (replacing/extending `validators_test.go`):
  - schema: defaults applied, name defaulting, `${VAR}` expansion (set/unset), unknown-key rejection
  - validation: duplicate (type,name), invalid names, per-type required fields, redirect URL paths (default + named Last.fm), CSV path conflicts, WiiM private IP
  - golden test: the full example in §4 loads into the expected specs
- `internal/targets` — `New`/`NewMultiple` construct instances with correct names/types from specs
- `internal/api` — router registers named Last.fm routes; health JSON includes `name`
- Existing scrobbler/rules tests — untouched (pipeline logic unchanged)
- Manual E2E: config with two `ListenBrainz` instances against a mock compatible service (e.g. the existing dev Compose files) → `/health` lists both by name, `/metrics` shows distinct `instance` label values

## 8. Migration notes for users (to include in release notes)

1. Breaking: all configuration env vars except `SCROBBLET_CONFIG` are removed; create `scrobblet.yaml`
2. Mapping table env var → YAML path
3. Secrets can stay in env vars and be referenced as `${VAR}`
4. Breaking metrics change: `scrobblet_submits_total` and `scrobblet_submit_duration_seconds` gain an `instance` label — update dashboards/alerts
5. Last.fm: single unnamed target keeps existing auth URLs and session file; named targets use `/api/lastfm/<name>/login` and need matching `redirectUrl`

## 9. Open questions

1. **Top-level `version: 1` key** in the YAML schema to ease future migrations? (Lean: yes, optional and ignored for now.) Answer: No
2. **Koito/Maloja deprecation**: keep as presets forever, or deprecate in favor of `type: ListenBrainz` + `url`? (Lean: keep — zero cost, nicer UX.) Answer: Keep
3. **Live reload** of the config file? (Lean: no — restart to reload, document it.) Answer: No - restart to reload

## 10. Out of scope

- Multiple sources (validator still enforces exactly one)
- Retry/backoff for failed submits (existing `todo` in `submitter.go`)
- Per-target enable/disable of now-playing vs scrobbles
