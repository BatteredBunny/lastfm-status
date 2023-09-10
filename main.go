package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	_ "embed"

	"astuart.co/goq"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed status.gohtml
var templateFile embed.FS

type TemplateInput struct {
	Title     string `goquery:".chartlist-name a"`
	TitleUrl  string `goquery:".chartlist-name a,[href]"`
	Author    string `goquery:".chartlist-artist a"`
	AuthorUrl string `goquery:".chartlist-artist a,[href]"`
	CoverArt  string `goquery:".chartlist-image .cover-art img,[src]"`

	AccountUrl  string
	AccountName string
}

type CacheField struct {
	TemplateInput
	CacheTime time.Time
}

// Input is the username of the scrobbler
var Cache = make(map[string]CacheField)

var CacheDuration = time.Minute

func CacheExpired(t time.Time) bool {
	return t.Before(time.Now().Add(-CacheDuration))
}

func CacheCleaner() {
	log.Println("Starting cache cleaner")

	for {
		time.Sleep(time.Hour)
		log.Println("Starting hourly cleaning")
		var deleteCounter int
		for name, v := range Cache {
			if CacheExpired(v.CacheTime) {
				delete(Cache, name)
				deleteCounter++
			}
		}

		log.Println("Deleted", deleteCounter, "scrobble caches")
	}
}

type RawPageInfo struct {
	Song *TemplateInput `goquery:".chartlist-row--now-scrobbling"`
}

// GetInfo fetches info from last.fm
func GetInfo(username string) (c CacheField, err error) {
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

	rawInfo.Song.AuthorUrl = "https://www.last.fm" + rawInfo.Song.AuthorUrl
	rawInfo.Song.TitleUrl = "https://www.last.fm" + rawInfo.Song.TitleUrl
	rawInfo.Song.AccountUrl = accountUrl
	rawInfo.Song.AccountName = username

	c = CacheField{
		TemplateInput: *rawInfo.Song,
		CacheTime:     time.Now(),
	}

	return
}

var StatusTemplate *template.Template

func main() {
	port := flag.Uint("port", 8080, "port to run server on")
	flag.DurationVar(&CacheDuration, "cache-length", CacheDuration, "how long to cache an entry for")
	ratelimiting := flag.Bool("ratelimit", true, "enables ratelimiting for /status api")
	flag.Parse()

	go CacheCleaner()

	var err error
	StatusTemplate, err = template.ParseFS(templateFile, "status.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Last.fm status running")
	})

	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	if *ratelimiting {
		ratelimit := tollbooth.NewLimiter(4, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
		http.Handle("/status", tollbooth.LimitFuncHandler(ratelimit, statusHandler))
	} else {
		http.HandleFunc("/status", statusHandler)
	}

	log.Printf("Starting server on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Please insert an username")
		return
	}

	var err error
	v, valid := Cache[username]
	if !valid {
		log.Printf("No cache for %s, refreshing data\n", username)
		v, err = GetInfo(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		Cache[username] = v
	} else if CacheExpired(v.CacheTime) {
		log.Printf("Cache is too old for %s, refreshing data\n", username)
		v, err = GetInfo(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		Cache[username] = v
	} else {
		log.Printf("Getting info for %s from cache\n", username)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	NotPlaying := v.TemplateInput.Author == ""
	if !NotPlaying {
		if err = StatusTemplate.Execute(w, v.TemplateInput); err != nil {
			log.Println("WARNING:", err)
		}
	}
}
