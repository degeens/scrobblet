package spotify

type CurrentlyPlayingTrack struct {
	IsPlaying bool `json:"is_playing"`
	Progress  int  `json:"progress_ms"`
	Item      Item `json:"item"`
}

type Item struct {
	Album    Album    `json:"album"`
	Artists  []Artist `json:"artists"`
	Duration int      `json:"duration_ms"`
	Name     string   `json:"name"`
}

type Album struct {
	Name string `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}
