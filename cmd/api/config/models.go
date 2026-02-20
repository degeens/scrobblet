package config

import (
	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type Config struct {
	Port           string
	DataPath       string
	LogLevel       string
	RateLimitRate  int
	RateLimitBurst int
	Source         sources.SourceType
	Targets        []targets.TargetType
	Clients        clients.Config
}
