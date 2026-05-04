package handlers

import (
	"embed"
	"net/http"
)

//go:embed static
var staticFS embed.FS

func Static() http.Handler {
	return http.FileServer(http.FS(staticFS))
}
