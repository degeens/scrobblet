package koito

type SubmitListens struct {
	ListenType string    `json:"listen_type"`
	Payload    []Payload `json:"payload"`
}

type Payload struct {
	ListenedAt    int64         `json:"listened_at"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
	ArtistName  string `json:"artist_name"`
	TrackName   string `json:"track_name"`
	ReleaseName string `json:"release_name"`
}
