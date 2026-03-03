package internal

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/gin-gonic/gin"
)

type Application struct {
	userListeningCache map[string]userListeningCache

	router      *gin.Engine
	ratelimiter *limiter.Limiter
	Config      Config
}

type Config struct {
	CacheDuration time.Duration // Currently listening cache duration

	Port uint

	RateLimiting bool

	BehindReverseProxy bool
	TrustedProxy       string
}

func parseConfig() (cfg Config) {
	flag.UintVar(&cfg.Port, "port", 8080, "port to run server on")
	flag.BoolVar(&cfg.BehindReverseProxy, "reverse-proxy", false, "Set true if behind reverse proxy to make the ratelimiter work")
	flag.StringVar(&cfg.TrustedProxy, "trusted-proxy", "", "trusted proxy for reverse prox")
	flag.DurationVar(&cfg.CacheDuration, "cache-length", time.Minute, "how long to cache an entry for")
	flag.BoolVar(&cfg.RateLimiting, "ratelimit", true, "enables ratelimiting for /status api")
	flag.Parse()

	return cfg
}

func (app *Application) setupRatelimiter() {
	if app.Config.RateLimiting {
		app.ratelimiter = tollbooth.NewLimiter(4, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	}

	if app.Config.BehindReverseProxy {
		app.ratelimiter.SetIPLookup(limiter.IPLookup{
			Name:           "X-Forwarded-For",
			IndexFromRight: 0,
		})
	} else {
		app.ratelimiter.SetIPLookup(limiter.IPLookup{
			Name:           "RemoteAddr",
			IndexFromRight: 0,
		})
	}
}

func NewApplication() Application {
	return Application{
		userListeningCache: make(map[string]userListeningCache),
		Config:             parseConfig(),
	}
}

func (app *Application) Run() {
	app.setupRatelimiter()
	if err := app.setupRouter(); err != nil {
		log.Fatal(err)
	}

	go app.cacheCleaner()

	log.Fatal(app.router.Run(fmt.Sprintf(":%d", app.Config.Port)))
}
