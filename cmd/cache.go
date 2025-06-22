package cmd

import (
	"log"
	"time"
)

type UserListeningCache struct {
	SongTitle   string
	SongUrl     string
	AuthorName  string
	AuthorUrl   string
	CoverArtUrl string
	AccountName string
	AccountUrl  string

	CacheTime time.Time
}

func (u UserListeningCache) Expired(CacheDuration time.Duration) bool {
	return u.CacheTime.Before(time.Now().Add(-CacheDuration))
}

func (app *Application) CacheCleaner() {
	log.Println("Starting cache cleaner")

	for {
		time.Sleep(time.Hour)
		log.Println("Starting hourly cleaning")

		var deleteCounter int
		for name, cache := range app.UserListeningCache {
			if cache.Expired(app.Config.CacheDuration) {
				delete(app.UserListeningCache, name)
				deleteCounter++
			}
		}

		log.Println("Deleted", deleteCounter, "scrobble caches")
	}
}
