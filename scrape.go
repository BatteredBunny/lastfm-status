package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"astuart.co/goq"
)

type RawPageInfo struct {
	CurrentlyScrobblingSong *RawCurrentlyScrobbling `goquery:".chartlist-row--now-scrobbling"`
}

type RawCurrentlyScrobbling struct {
	SongTitle   string `goquery:".chartlist-name a"`
	SongUrl     string `goquery:".chartlist-name a,[href]"`
	AuthorName  string `goquery:".chartlist-artist a"`
	AuthorUrl   string `goquery:".chartlist-artist a,[href]"`
	CoverArtUrl string `goquery:".chartlist-image .cover-art img,[src]"`
}

// GetCurrentlyScrobbling fetches info from last.fm
func GetCurrentlyScrobbling(username string) (c UserCache, err error) {
	accountUrl := "https://www.last.fm/user/" + username
	resp, err := http.Get(accountUrl)
	if err != nil {
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rawInfo RawPageInfo
	if err = goq.NewDecoder(bytes.NewReader(b)).Decode(&rawInfo); err != nil {
		return
	}

	c = UserCache{
		SongTitle:   rawInfo.CurrentlyScrobblingSong.SongTitle,
		SongUrl:     "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.SongUrl,
		AuthorName:  rawInfo.CurrentlyScrobblingSong.AuthorName,
		AuthorUrl:   "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.AuthorUrl,
		CoverArtUrl: rawInfo.CurrentlyScrobblingSong.CoverArtUrl,
		AccountName: username,
		AccountUrl:  accountUrl,

		CacheTime: time.Now(),
	}

	return
}
