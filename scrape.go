package main

import (
	"bytes"
	"fmt"
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
func GetCurrentlyScrobbling(username string) (c UserListeningCache, err error) {
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

	c = UserListeningCache{
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

type RawAlbum struct {
	Name     string `goquery:".link-block-target"`
	CoverArt string `goquery:".grid-items-cover-image-image img,[src]"`

	ArtistName string `goquery:".grid-items-item-aux-block"`
	ArtistUrl  string `goquery:".grid-items-item-aux-block,[href]"`

	RawPlays []string `goquery:".grid-items-item-aux-text a"`
}

type RawAlbumInfo struct {
	Albums []RawAlbum `goquery:".grid-items .grid-items-item"`
}

func GetMonthlyArtists(username string) (c UserMonthlyAlbumsCache, err error) {
	albumsUrl := fmt.Sprintf("https://www.last.fm/user/%s/partial/albums?albums_date_preset=LAST_30_DAYS", username)

	resp, err := http.Get(albumsUrl)
	if err != nil {
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rawInfo RawAlbumInfo
	if err = goq.NewDecoder(bytes.NewReader(b)).Decode(&rawInfo); err != nil {
		return
	}

	for _, album := range rawInfo.Albums {
		c.Albums = append(c.Albums, Album{
			Name:       album.Name,
			CoverArt:   album.CoverArt,
			ArtistName: album.ArtistName,
			ArtistUrl:  "https://www.last.fm" + album.ArtistUrl,
			Plays:      album.RawPlays[1],
		})
	}

	c.CacheTime = time.Now()

	return
}
