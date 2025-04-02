package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"time"

	_ "embed"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

//go:embed static
var StaticFiles embed.FS

//go:embed template
var Templates embed.FS

type Application struct {
	UserListeningCache     map[string]UserListeningCache
	UserMonthlyAlbumsCache map[string]UserMonthlyAlbumsCache

	Router *gin.Engine

	Ratelimiter *limiter.Limiter

	Config Config
}

type Config struct {
	CacheDuration        time.Duration // Currently listening cache duration
	MonthlyCacheDuration time.Duration // Monthly top artists cache duration

	Port uint

	RateLimiting bool
}

func ParseConfig() (cfg Config) {
	flag.UintVar(&cfg.Port, "port", 8080, "port to run server on")
	flag.DurationVar(&cfg.CacheDuration, "cache-length", time.Minute, "how long to cache an entry for")
	flag.DurationVar(&cfg.MonthlyCacheDuration, "monthly-cache-length", time.Hour, "how long to cache an entry for")
	flag.BoolVar(&cfg.RateLimiting, "ratelimit", true, "enables ratelimiting for /status api")
	flag.Parse()

	return cfg
}

func (app *Application) SetupRatelimiter() {
	if app.Config.RateLimiting {
		app.Ratelimiter = tollbooth.NewLimiter(4, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	}
}

func main() {
	app := Application{
		UserListeningCache:     make(map[string]UserListeningCache),
		UserMonthlyAlbumsCache: make(map[string]UserMonthlyAlbumsCache),
		Config:                 ParseConfig(),
	}

	app.SetupRatelimiter()
	if err := app.SetupRouter(); err != nil {
		log.Fatal(err)
	}

	go app.CacheCleaner()

	log.Fatal(app.Router.Run(fmt.Sprintf(":%d", app.Config.Port)))
}
