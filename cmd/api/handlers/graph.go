package handlers

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/degeens/scrobblet/internal/health"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

// PlaybackStateProvider returns the last known playback state from the tracker.
// Returns nil when nothing is currently playing.
type PlaybackStateProvider interface {
	LastPlaybackState() *sources.PlaybackState
}

//go:embed templates/*.html
var graphTemplateFS embed.FS

type graphNode struct {
	ID         string
	Label      string
	Role       string
	Healthy    bool
	ShowHealth bool
	BoxX       int
	BoxY       int
	Fill       string
	Stroke     string
}

type graphEdge struct {
	ID              string
	X1, Y1, X2, Y2 int
	MidX            int
	Active          bool
}

type nowPlayingData struct {
	Artist     string
	Title      string
	Album      string
	PositionMs int
	DurationMs int
	IsPlaying  bool
}

type jsTargetConfig struct {
	Label string `json:"label"`
	Index int    `json:"index"`
}

type jsEdgeConfig struct {
	X1   int `json:"x1"`
	Y1   int `json:"y1"`
	X2   int `json:"x2"`
	Y2   int `json:"y2"`
	MidX int `json:"midX"`
}

type graphData struct {
	Source        graphNode
	ScrobbletNode graphNode
	Targets       []graphNode
	SourceEdge    graphEdge
	Edges         []graphEdge
	SVGWidth      int
	SVGHeight     int
	Now           time.Time
	NowPlaying    *nowPlayingData
	JSTargets     template.JS
	JSEdges       template.JS
}

const graphNodeW, graphNodeH = 140, 50

func Graph(source sources.Source, provider PlaybackStateProvider, ts []targets.Target) http.HandlerFunc {
	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).ParseFS(graphTemplateFS,
		"templates/index.html",
		"templates/graph.html",
		"templates/graph-defs.html",
		"templates/graph-node.html",
		"templates/graph-edge.html",
		"templates/graph-legend.html",
		"templates/now-playing.html",
	))

	return func(w http.ResponseWriter, r *http.Request) {
		hc := health.CheckHealth(source, ts)

		var nowPlaying *nowPlayingData
		var isPlayingInit bool
		if state := provider.LastPlaybackState(); state != nil {
			nowPlaying = &nowPlayingData{
				Artist:     strings.Join(state.Track.Artists, ", "),
				Title:      state.Track.Title,
				Album:      state.Track.Album,
				PositionMs: int(state.Position / time.Millisecond),
				DurationMs: int(state.Track.Duration / time.Millisecond),
				IsPlaying:  state.IsPlaying,
			}
			isPlayingInit = state.IsPlaying
		}

		n := len(ts)
		svgWidth := 1000
		svgHeight := max(120, n*70+40)

		col1X := 30
		col2X := (svgWidth - graphNodeW) / 2
		col3X := svgWidth - graphNodeW - 30

		centerY := svgHeight/2 - graphNodeH/2

		sourceStroke := "#444"
		if isPlayingInit {
			sourceStroke = "#1db954"
		}

		sourceNode := graphNode{
			ID:         "source-rect",
			Label:      string(hc.Source.SourceType),
			Role:       "SOURCE",
			Healthy:    hc.Source.Healthy,
			ShowHealth: true,
			BoxX:       col1X,
			BoxY:       centerY,
			Fill:       "#0f1f2e",
			Stroke:     sourceStroke,
		}

		scrobbletNode := graphNode{
			Label:  "Scrobblet",
			BoxX:   col2X,
			BoxY:   centerY,
			Fill:   "#13131f",
			Stroke: "#666680",
		}

		sourceEdge := graphEdge{
			ID:     "source-edge",
			X1:     col1X + graphNodeW,
			Y1:     centerY + graphNodeH/2,
			X2:     col2X,
			Y2:     centerY + graphNodeH/2,
			MidX:   (col1X + graphNodeW + col2X) / 2,
			Active: isPlayingInit,
		}

		tx1 := col2X + graphNodeW
		ty1 := centerY + graphNodeH/2

		targetNodes := make([]graphNode, n)
		edges := make([]graphEdge, n)
		jsTargets := make([]jsTargetConfig, n)
		jsEdges := make([]jsEdgeConfig, n)

		for i, t := range hc.Targets {
			targetBoxY := 20 + i*70
			targetNodes[i] = graphNode{
				ID:         fmt.Sprintf("target-rect-%d", i),
				Label:      string(t.TargetType),
				Role:       "TARGET",
				Healthy:    t.Healthy,
				ShowHealth: true,
				BoxX:       col3X,
				BoxY:       targetBoxY,
				Fill:       "#13132a",
				Stroke:     "#2a2a4a",
			}
			x2 := col3X
			y2 := targetBoxY + graphNodeH/2
			edges[i] = graphEdge{
				ID:   fmt.Sprintf("edge-%d", i),
				X1:   tx1, Y1: ty1,
				X2:   x2, Y2: y2,
				MidX: (tx1 + x2) / 2,
			}
			jsTargets[i] = jsTargetConfig{Label: string(t.TargetType), Index: i}
			jsEdges[i] = jsEdgeConfig{X1: tx1, Y1: ty1, X2: x2, Y2: y2, MidX: (tx1 + x2) / 2}
		}

		marshalJS := func(v any) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		}

		data := graphData{
			Source:        sourceNode,
			ScrobbletNode: scrobbletNode,
			Targets:       targetNodes,
			SourceEdge:    sourceEdge,
			Edges:         edges,
			SVGWidth:      svgWidth,
			SVGHeight:     svgHeight,
			Now:           time.Now().UTC(),
			NowPlaying:    nowPlaying,
			JSTargets:     marshalJS(jsTargets),
			JSEdges:       marshalJS(jsEdges),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
