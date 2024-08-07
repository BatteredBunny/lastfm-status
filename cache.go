package main

import (
	"log"
	"time"
)

type UserCache struct {
	SongTitle   string
	SongUrl     string
	AuthorName  string
	AuthorUrl   string
	CoverArtUrl string
	AccountName string
	AccountUrl  string

	CacheTime time.Time
}

func (u UserCache) Expired(CacheDuration time.Duration) bool {
	return u.CacheTime.Before(time.Now().Add(-CacheDuration))
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
