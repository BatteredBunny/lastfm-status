package lastfm

import (
	"bytes"
	"io"
	"net/http"

	"astuart.co/goq"
)

type rawPageInfo struct {
	CurrentlyScrobblingSong *rawCurrentlyScrobbling `goquery:".chartlist-row--now-scrobbling"`
}

type rawCurrentlyScrobbling struct {
	SongTitle   string `goquery:".chartlist-name a"`
	SongUrl     string `goquery:".chartlist-name a,[href]"`
	AuthorName  string `goquery:".chartlist-artist a"`
	AuthorUrl   string `goquery:".chartlist-artist a,[href]"`
	CoverArtUrl string `goquery:".chartlist-image .cover-art img,[src]"`
}

type Scrobble struct {
	SongTitle   string
	SongUrl     string
	AuthorName  string
	AuthorUrl   string
	CoverArtUrl string
	AccountName string
	AccountUrl  string
}

// GetCurrentlyScrobbling fetches info from last.fm
func GetCurrentlyScrobbling(username string) (c Scrobble, err error) {
	accountUrl := "https://www.last.fm/user/" + username
	resp, err := http.Get(accountUrl)
	if err != nil {
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rawInfo rawPageInfo
	if err = goq.NewDecoder(bytes.NewReader(b)).Decode(&rawInfo); err != nil {
		return
	}

	c = Scrobble{
		SongTitle:   rawInfo.CurrentlyScrobblingSong.SongTitle,
		SongUrl:     "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.SongUrl,
		AuthorName:  rawInfo.CurrentlyScrobblingSong.AuthorName,
		AuthorUrl:   "https://www.last.fm" + rawInfo.CurrentlyScrobblingSong.AuthorUrl,
		CoverArtUrl: rawInfo.CurrentlyScrobblingSong.CoverArtUrl,
		AccountName: username,
		AccountUrl:  accountUrl,
	}

	return
}
