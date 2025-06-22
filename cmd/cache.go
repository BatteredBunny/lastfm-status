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

type Album struct {
	Name     string
	CoverArt string
	Plays    string

	ArtistName string
	ArtistUrl  string
}

type UserMonthlyAlbumsCache struct {
	Albums    []Album
	CacheTime time.Time
}

func (u UserMonthlyAlbumsCache) Expired(CacheDuration time.Duration) bool {
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

		deleteCounter = 0
		for name, cache := range app.UserMonthlyAlbumsCache {
			if cache.Expired(app.Config.MonthlyCacheDuration) {
				delete(app.UserMonthlyAlbumsCache, name)
				deleteCounter++
			}
		}

		log.Println("Deleted", deleteCounter, "monthly scrobble caches")
	}
}
