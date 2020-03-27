package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/umahmood/gfycatdl"
)

const usage = `Download gfycat gifs.

Usage:

    -help - Print this message and exit.
    -url  - Gfycat url
`
const examples = `Example usage:

     $ gfycatdl -url https://gfycat.com/violetsmartalleycat-sunset-dusk-nature
`

func printUsage() {
	fmt.Println(usage)
	fmt.Println(examples)
}

func main() {
	var (
		help bool
		url  string
	)
	flag.Usage = func() {
		printUsage()
	}
	flag.BoolVar(&help, "help", false, "Print this message and exit")
	flag.StringVar(&url, "url", "", "gfycat url to download")
	flag.Parse()
	if flag.NFlag() == 0 {
		printUsage()
		os.Exit(1)
	}
	if help {
		printUsage()
		os.Exit(0)
	}
	g, err := gfycatdl.New(url)
	if err != nil {
		log.Fatalln(err)
	}
	scrapedURL, err := g.ScrapeVideoSource()
	if err != nil {
		log.Fatalln(err)
	}
	err = gfycatdl.DownloadFile(g.ResourceName, scrapedURL)
	if err != nil {
		log.Fatalln(err)
	}
}
