package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"time"

	_ "embed"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
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

	BehindReverseProxy bool
	TrustedProxy       string
}

func ParseConfig() (cfg Config) {
	flag.UintVar(&cfg.Port, "port", 8080, "port to run server on")
	flag.BoolVar(&cfg.BehindReverseProxy, "reverse-proxy", false, "Set true if behind reverse proxy to make the ratelimiter work")
	flag.StringVar(&cfg.TrustedProxy, "trusted-proxy", "", "trusted proxy for reverse prox")
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

	if app.Config.BehindReverseProxy {
		app.Ratelimiter.SetIPLookup(limiter.IPLookup{
			Name:           "X-Forwarded-For",
			IndexFromRight: 0,
		})
	} else {
		app.Ratelimiter.SetIPLookup(limiter.IPLookup{
			Name:           "RemoteAddr",
			IndexFromRight: 0,
		})
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
