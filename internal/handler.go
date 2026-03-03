package internal

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	embed "github.com/BatteredBunny/lastfm-status"
	"github.com/BatteredBunny/lastfm-status/internal/lastfm"
	"github.com/gin-gonic/gin"
)

func (app *Application) setupRouter() (err error) {
	app.router = gin.Default()
	app.router.ForwardedByClientIP = app.Config.BehindReverseProxy
	app.router.SetTrustedProxies([]string{app.Config.TrustedProxy})

	app.router.SetHTMLTemplate(template.Must(template.ParseFS(embed.Templates, "template/*.gohtml")))

	app.router.StaticFileFS("/", "static/html/main.html", http.FS(embed.StaticFiles))

	var sub fs.FS
	sub, err = fs.Sub(embed.StaticFiles, "static/css")
	if err != nil {
		return
	}

	app.router.StaticFS("/css", http.FS(sub))

	if app.Config.RateLimiting {
		log.Println("Enabling ratelimiting")
		app.router.GET("/status", app.ratelimiterMiddleware(), app.statusHandler)
	} else {
		app.router.GET("/status", app.statusHandler)
	}

	return
}

type templateInput struct {
	userListeningCache

	Refresh float64
	Light   bool
	Dark    bool
	Dynamic bool
}

type statusQuery struct {
	Theme    string `form:"theme"` // light, dark, dynamic
	Username string `form:"username" binding:"required"`
}

func (app *Application) statusHandler(c *gin.Context) {
	var input statusQuery
	if err := c.BindQuery(&input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var light, dark, dynamic bool
	switch input.Theme {
	case "light":
		light = true
	case "dark":
		dark = true
	case "dynamic":
		dynamic = true
	default:
		dynamic = true
	}

	var err error
	cache, exists := app.userListeningCache[input.Username]
	if !exists || cache.Expired(app.Config.CacheDuration) {
		log.Printf("Refreshing cache for %s\n", input.Username)
		var scrobble lastfm.Scrobble
		scrobble, err = lastfm.GetCurrentlyScrobbling(input.Username)
		cache = userListeningCache{
			Scrobble:  scrobble,
			CacheTime: time.Now(),
		}
		if err != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		app.userListeningCache[input.Username] = cache
	} else {
		log.Printf("Getting info for %s from cache\n", input.Username)
	}

	NotPlaying := cache.AuthorName == ""
	if !NotPlaying {
		c.HTML(http.StatusOK, "status.gohtml", templateInput{
			Refresh: app.Config.CacheDuration.Seconds(),
			Light:   light,
			Dark:    dark,
			Dynamic: dynamic,

			userListeningCache: cache,
		})
	}
}
