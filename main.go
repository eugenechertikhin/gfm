package main

import (
	"flag"
	"fmt"
	"gfm/app"
)

var (
	ascii  = flag.Bool("a", false, "switch to use ascii border, not utf")
	scheme = flag.String("scheme", "colour", "colour scheme (custom, colour, bw)")
)

func main() {
	flag.Parse()
	
	app.Init(*ascii, *scheme)
	defer app.Finish()

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
