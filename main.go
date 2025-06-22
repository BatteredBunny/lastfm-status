package main

import (
	"github.com/BatteredBunny/lastfm-status/cmd"
)

func main() {
	app := cmd.NewApplication()
	app.Run()
}
