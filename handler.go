package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) SetupRouter() (err error) {
	app.Router = gin.Default()
	app.Router.ForwardedByClientIP = app.Config.BehindReverseProxy
	app.Router.SetTrustedProxies([]string{app.Config.TrustedProxy})

	app.Router.SetHTMLTemplate(template.Must(template.ParseFS(Templates, "template/*.gohtml")))

	app.Router.StaticFileFS("/", "static/html/main.html", http.FS(StaticFiles))

	var sub fs.FS
	sub, err = fs.Sub(StaticFiles, "static/css")
	if err != nil {
		return
	}

	app.Router.StaticFS("/css", http.FS(sub))

	if app.Config.RateLimiting {
		log.Println("Enabling ratelimiting")
		app.Router.GET("/status", app.RatelimiterMiddleware(), app.StatusHandler)
		app.Router.GET("/monthly", app.RatelimiterMiddleware(), app.MonthlyHandler)
	} else {
		app.Router.GET("/status", app.StatusHandler)
		app.Router.GET("/monthly", app.MonthlyHandler)
	}

	return
}

type TemplateInput struct {
	UserListeningCache

	Refresh float64
	Light   bool
	Dark    bool
	Dynamic bool
}

type StatusQuery struct {
	Theme    string `form:"theme"` // light, dark, dynamic
	Username string `form:"username" binding:"required"`
}

func (app *Application) StatusHandler(c *gin.Context) {
	var input StatusQuery
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
	cache, exists := app.UserListeningCache[input.Username]
	if !exists || cache.Expired(app.Config.CacheDuration) {
		log.Printf("Refreshing cache for %s\n", input.Username)
		cache, err = GetCurrentlyScrobbling(input.Username)
		if err != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		app.UserListeningCache[input.Username] = cache
	} else {
		log.Printf("Getting info for %s from cache\n", input.Username)
	}

	NotPlaying := cache.AuthorName == ""
	if !NotPlaying {
		c.HTML(http.StatusOK, "status.gohtml", TemplateInput{
			Refresh: app.Config.CacheDuration.Seconds(),
			Light:   light,
			Dark:    dark,
			Dynamic: dynamic,

			UserListeningCache: cache,
		})
	}
}

type MonthQuery struct {
	Theme    string `form:"theme"` // bright
	Username string `form:"username" binding:"required"`
}

type MonthlyTemplateInput struct {
	UserMonthlyAlbumsCache

	Refresh float64

	Bright bool
}

func (app *Application) MonthlyHandler(c *gin.Context) {
	var input MonthQuery
	if err := c.BindQuery(&input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	bright := input.Theme == "bright"

	var err error
	cache, exists := app.UserMonthlyAlbumsCache[input.Username]
	if !exists || cache.Expired(app.Config.CacheDuration) {
		log.Printf("Refreshing monthly plays cache for %s\n", input.Username)
		cache, err = GetMonthlyArtists(input.Username)
		if err != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		app.UserMonthlyAlbumsCache[input.Username] = cache
	} else {
		log.Printf("Getting info for %s from monthly cache\n", input.Username)
	}

	c.HTML(http.StatusOK, "monthly.gohtml", MonthlyTemplateInput{
		Refresh: app.Config.MonthlyCacheDuration.Seconds(),
		Bright:  bright,

		UserMonthlyAlbumsCache: cache,
	})
}
