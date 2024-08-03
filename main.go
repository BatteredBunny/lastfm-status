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
	Cache  map[string]UserCache
	Router *gin.Engine

	Ratelimiter *limiter.Limiter

	Config Config
}

type Config struct {
	CacheDuration time.Duration
	Port          uint

	RateLimiting bool
}

func ParseConfig() (cfg Config) {
	flag.UintVar(&cfg.Port, "port", 8080, "port to run server on")
	flag.DurationVar(&cfg.CacheDuration, "cache-length", time.Minute, "how long to cache an entry for")
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
		Cache:  make(map[string]UserCache),
		Config: ParseConfig(),
	}

	app.SetupRatelimiter()
	if err := app.SetupRouter(); err != nil {
		log.Fatal(err)
	}

	go app.CacheCleaner()

	log.Fatal(app.Router.Run(fmt.Sprintf(":%d", app.Config.Port)))
}
