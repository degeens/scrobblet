package wiim

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type PlayerStatus struct {
	Title           string `json:"title"`
	Artist          string `json:"artist"`
	Album           string `json:"album"`
	CurrentPosition int    `json:"curpos,string"`
	TotalLength     int    `json:"totlen,string"`
}

func (p *PlayerStatus) UnmarshalJSON(data []byte) error {
	// Use an alias so json.Unmarshal does not call this method recursively.
	type playerStatusAlias PlayerStatus
	var raw playerStatusAlias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var err error
	raw.Title, err = hexToText(raw.Title)
	if err != nil {
		return fmt.Errorf("decode title: %w", err)
	}

	raw.Artist, err = hexToText(raw.Artist)
	if err != nil {
		return fmt.Errorf("decode artist: %w", err)
	}

	raw.Album, err = hexToText(raw.Album)
	if err != nil {
		return fmt.Errorf("decode album: %w", err)
	}

	*p = PlayerStatus(raw)
	return nil
}

func hexToText(s string) (string, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
