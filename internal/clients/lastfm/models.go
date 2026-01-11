package lastfm

import "encoding/json"

type UpdateNowPlayingRequest struct {
	Artist   string
	Track    string
	Album    string
	Duration int // In seconds
}

type UpdateNowPlayingResponse struct {
	NowPlaying struct {
		IgnoredMessage IgnoredMessage `json:"ignoredMessage"`
	} `json:"nowplaying"`
	Error   *int    `json:"error"`
	Message *string `json:"message"`
}

type ScrobbleRequest struct {
	Artist    string
	Track     string
	Album     string
	Duration  int   // In seconds
	Timestamp int64 // UNIX timestamp
}

type ScrobbleResponse struct {
	Scrobbles struct {
		Attr struct {
			Accepted int `json:"accepted"`
			Ignored  int `json:"ignored"`
		} `json:"@attr"`
		Scrobble ScrobbleList `json:"scrobble"`
	} `json:"scrobbles"`
	Error   *int    `json:"error"`
	Message *string `json:"message"`
}

// ScrobbleList handles Last.fm's inconsistent response format
type ScrobbleList []Scrobble

func (s *ScrobbleList) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as array first
	var arr []Scrobble
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = arr
		return nil
	}

	// If that fails, try as a single object
	var single Scrobble
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}

	*s = []Scrobble{single}

	return nil
}

type Scrobble struct {
	Track struct {
		Text string `json:"#text"`
	} `json:"track"`
	Artist struct {
		Text string `json:"#text"`
	} `json:"artist"`
	Album struct {
		Text string `json:"#text"`
	} `json:"album"`
	IgnoredMessage IgnoredMessage `json:"ignoredMessage"`
}

type IgnoredMessage struct {
	Code string `json:"code"`
	Text string `json:"#text"`
}
