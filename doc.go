/*
Package gfycatdl downloads gifs from gfycat given a URL

Usage:
    package main

    import (
        "log"
        "github.com/umahmood/gfycatdl"
    )
    func main() {
        g, err := gfycatdl.New("https://gfycat.com/violetsmartalleycat-sunset-dusk-nature")
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
*/
package gfycatdl
