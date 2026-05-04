package events

import "sync"

type EventType string

const (
	EventScrobble    EventType = "scrobble"
	EventNowPlaying  EventType = "now_playing"
	EventTrackChange EventType = "track_change"
)

type Event struct {
	Type       EventType `json:"type"`
	Target     string    `json:"target"`
	Success    bool      `json:"success"`
	Artist     string    `json:"artist"`
	Title      string    `json:"title"`
	Album      string    `json:"album"`
	DurationMs int64     `json:"durationMs"`
	IsPlaying  bool      `json:"isPlaying"`
	PositionMs int64     `json:"positionMs"`
}

type Bus struct {
	mu   sync.Mutex
	subs []chan Event
}

func NewBus() *Bus {
	return &Bus{}
}

func (b *Bus) Publish(e Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}

func (b *Bus) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 16)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()

	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, sub := range b.subs {
			if sub == ch {
				b.subs = append(b.subs[:i], b.subs[i+1:]...)
				close(ch)
				return
			}
		}
	}
}
