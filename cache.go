package main

import (
	"log"
	"time"
)

type CacheField struct {
	SongTitle   string
	SongUrl     string
	AuthorName  string
	AuthorUrl   string
	CoverArtUrl string
	AccountName string
	AccountUrl  string

	CacheTime time.Time
}

func (c CacheField) Expired(CacheDuration time.Duration) bool {
	return c.CacheTime.Before(time.Now().Add(-CacheDuration))
}

func (app *Application) CacheCleaner() {
	log.Println("Starting cache cleaner")

	for {
		time.Sleep(time.Hour)
		log.Println("Starting hourly cleaning")
		var deleteCounter int
		for name, cache := range app.Cache {
			if cache.Expired(app.Config.CacheDuration) {
				delete(app.Cache, name)
				deleteCounter++
			}
		}

		log.Println("Deleted", deleteCounter, "scrobble caches")
	}
}
