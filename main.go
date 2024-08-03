package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "embed"
)

//go:embed static/*
var StaticFiles embed.FS

//go:embed status.gohtml
var TemplateFile embed.FS

type Application struct {
	StatusTemplate *template.Template
	Cache          map[string]CacheField

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

func (app *Application) SetupTemplates() (err error) {
	app.StatusTemplate, err = template.ParseFS(TemplateFile, "status.gohtml")
	return
}

func main() {
	app := Application{
		Cache:  make(map[string]CacheField),
		Config: ParseConfig(),
	}

	if err := app.SetupTemplates(); err != nil {
		log.Fatal(err)
		return
	}

	go app.CacheCleaner()

	app.SetupHandlers()

	log.Printf("Starting server on :%d\n", app.Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil))
}
