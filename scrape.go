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
	Title     string `goquery:".chartlist-name a"`
	TitleUrl  string `goquery:".chartlist-name a,[href]"`
	Author    string `goquery:".chartlist-artist a"`
	AuthorUrl string `goquery:".chartlist-artist a,[href]"`
	CoverArt  string `goquery:".chartlist-image .cover-art img,[src]"`
}

// GetCurrentlyScrobbling fetches info from last.fm
func GetCurrentlyScrobbling(username string) (c CacheField, err error) {
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

	c = CacheField{
		SongTitle:   rawInfo.CurrentlyScrobblingSong.Title,
		SongUrl:     "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.TitleUrl,
		AuthorName:  rawInfo.CurrentlyScrobblingSong.Author,
		AuthorUrl:   "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.AuthorUrl,
		CoverArtUrl: rawInfo.CurrentlyScrobblingSong.CoverArt,
		AccountName: username,
		AccountUrl:  accountUrl,

		CacheTime: time.Now(),
	}

	return
}
