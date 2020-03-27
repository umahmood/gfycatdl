# Gfycatdl

Gfycat downloads gifs from [gfycat.com](https://gfycat.com/)

Gfycatdl is a [Go](https://golang.org/) library and command line tool.

**Note:** gfycat does have an api which requires an account and each request to
be authenticated. This tool simply scrapes the page and gets a download link.

# Installation

> $ go get github.com/umahmood/gfycatdl

# Usage

Command line:

> $ gfycatdl -url https://gfycat.com/violetsmartalleycat-sunset-dusk-nature

> $ gfycat -help
```
Download gfycat gifs.

Usage:

    -help - Print this message and exit.
    -url  - Gfycat url

Example usage:

     $ gfycatdl -url https://gfycat.com/violetsmartalleycat-sunset-dusk-nature
```

Library:

```
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
```

# Documentation

- https://pkg.go.dev/umahmood/gfycatdl

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
