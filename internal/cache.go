package internal

import (
	"log"
	"time"

	"github.com/BatteredBunny/lastfm-status/internal/lastfm"
)

type userListeningCache struct {
	lastfm.Scrobble
	CacheTime time.Time
}

func (u userListeningCache) Expired(CacheDuration time.Duration) bool {
	return u.CacheTime.Before(time.Now().Add(-CacheDuration))
}

func (app *Application) cacheCleaner() {
	log.Println("Starting cache cleaner")

	for {
		time.Sleep(time.Hour)
		log.Println("Starting hourly cleaning")

		var deleteCounter int
		for name, cache := range app.userListeningCache {
			if cache.Expired(app.Config.CacheDuration) {
				delete(app.userListeningCache, name)
				deleteCounter++
			}
		}

		log.Println("Deleted", deleteCounter, "scrobble caches")
	}
}
