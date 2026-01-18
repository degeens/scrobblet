package listenbrainz

type SubmitListens struct {
	ListenType string    `json:"listen_type"`
	Payload    []Payload `json:"payload"`
}

type Payload struct {
	ListenedAt    *int64        `json:"listened_at,omitempty"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
	ArtistName     string         `json:"artist_name"`
	TrackName      string         `json:"track_name"`
	ReleaseName    string         `json:"release_name"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type AdditionalInfo struct {
	ArtistNames             []string `json:"artist_names"`
	SubmissionClient        string   `json:"submission_client"`
	SubmissionClientVersion string   `json:"submission_client_version"`
}
