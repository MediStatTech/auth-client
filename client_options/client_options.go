package client_options

import (
	"time"

	"github.com/MediStatTech/env"
	log "github.com/MediStatTech/logger"
)

const (
	ContextTimeout            = 15 * time.Second
	DefaultProductionAddress  = "localhost:50051"
	DefaultDevelopmentAddress = "localhost:8080"
	DevPort                   = 8080 // Port for development environment
)

type Options struct {
	Log         *log.Logger
	ENV         *env.ENV
	AddressName string // The address to connect to (e.g., "localhost:8080")
}
