package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
)

func (app *Application) SetupHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Last.fm status running")
	})

	http.Handle("/static/", http.FileServer(http.FS(StaticFiles)))

	if *&app.Config.RateLimiting {
		ratelimit := tollbooth.NewLimiter(4, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
		http.Handle("/status", tollbooth.LimitFuncHandler(ratelimit, app.StatusHandler))
	} else {
		http.HandleFunc("/status", app.StatusHandler)
	}
}

type TemplateInput struct {
	CacheField

	Refresh float64
	Light   bool
	Dark    bool
	Dynamic bool
}

func (app *Application) StatusHandler(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	var light, dark, dynamic bool
	switch theme {
	case "light":
		light = true
	case "dark":
		dark = true
	case "dynamic":
		dynamic = true
	default:
		dynamic = true
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Please insert an username")
		return
	}

	var err error
	cache, valid := app.Cache[username]
	if !valid {
		log.Printf("No cache for %s, refreshing data\n", username)
		cache, err = GetCurrentlyScrobbling(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		app.Cache[username] = cache
	} else if cache.Expired(app.Config.CacheDuration) {
		log.Printf("Cache has expired for %s, refreshing data\n", username)
		cache, err = GetCurrentlyScrobbling(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		app.Cache[username] = cache
	} else {
		log.Printf("Getting info for %s from cache\n", username)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	NotPlaying := cache.AuthorName == ""
	if !NotPlaying {
		tmpi := TemplateInput{
			Refresh: app.Config.CacheDuration.Seconds(),
			Light:   light,
			Dark:    dark,
			Dynamic: dynamic,

			CacheField: cache,
		}
		if err = app.StatusTemplate.Execute(w, tmpi); err != nil {
			log.Println("WARNING:", err)
		}
	}
}
