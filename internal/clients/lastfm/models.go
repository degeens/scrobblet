package lastfm

// UpdateNowPlayingRequest represents the request for track.updateNowPlaying
type UpdateNowPlayingRequest struct {
	Artist      string
	Track       string
	Album       string
	Duration    int // in seconds
	TrackNumber int
	AlbumArtist string
}

// ScrobbleRequest represents a single scrobble
type ScrobbleRequest struct {
	Artist      string
	Track       string
	Timestamp   int64 // UNIX timestamp
	Album       string
	TrackNumber int
	Duration    int // in seconds
	AlbumArtist string
}

// ScrobbleResponse represents the JSON response from track.scrobble
type ScrobbleResponse struct {
	Scrobbles struct {
		Attr struct {
			Accepted string `json:"accepted"`
			Ignored  string `json:"ignored"`
		} `json:"@attr"`
		Scrobble []Scrobble `json:"scrobble"`
	} `json:"scrobbles"`
	Error   *int    `json:"error"`
	Message *string `json:"message"`
}

// NowPlayingResponse represents the JSON response from track.updateNowPlaying
type NowPlayingResponse struct {
	NowPlaying struct {
		Track struct {
			Text      string `json:"#text"`
			Corrected string `json:"corrected"`
		} `json:"track"`
		Artist struct {
			Text      string `json:"#text"`
			Corrected string `json:"corrected"`
		} `json:"artist"`
		Album struct {
			Text      string `json:"#text"`
			Corrected string `json:"corrected"`
		} `json:"album"`
		AlbumArtist struct {
			Text      string `json:"#text"`
			Corrected string `json:"corrected"`
		} `json:"albumArtist"`
		IgnoredMessage IgnoredMessage `json:"ignoredMessage"`
	} `json:"nowplaying"`
	Error   *int    `json:"error"`
	Message *string `json:"message"`
}

// Scrobble represents a single scrobble in the response
type Scrobble struct {
	Track struct {
		Text      string `json:"#text"`
		Corrected string `json:"corrected"`
	} `json:"track"`
	Artist struct {
		Text      string `json:"#text"`
		Corrected string `json:"corrected"`
	} `json:"artist"`
	Album struct {
		Text      string `json:"#text"`
		Corrected string `json:"corrected"`
	} `json:"album"`
	AlbumArtist struct {
		Text      string `json:"#text"`
		Corrected string `json:"corrected"`
	} `json:"albumArtist"`
	Timestamp      string         `json:"timestamp"`
	IgnoredMessage IgnoredMessage `json:"ignoredMessage"`
}

// IgnoredMessage represents the ignored message in responses
type IgnoredMessage struct {
	Code string `json:"code"`
	Text string `json:"#text"`
}

// ErrorResponse represents an error response from the API (legacy, now handled inline)
type ErrorResponse struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
}
