package gfycatdl

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ErrBadScrape failed to find gfycat downloadable media
var ErrBadScrape = errors.New("failed to find gfycat downloadable media")

// ErrBadDomain url is not a valid gfycat url i.e. https://gfycat.com/violetsmartalleycat-sunset-dusk-nature
var ErrBadDomain = errors.New("url is not a valid gfycat url")

// Gfycatdl instance
type Gfycatdl struct {
	ResourceName string // Name gyfcat resource name
	url          string
}

// New constructs a Gfycatdl instance
func New(gfycatURL string) (*Gfycatdl, error) {
	u, err := url.Parse(gfycatURL)
	if err != nil {
		return nil, ErrBadDomain
	}
	// allow localhost domain to aid in testing
	if u.Host != "gfycat.com" && !strings.Contains(u.Host, "127.0.0.1") &&
		!strings.Contains(u.Host, "::1") {
		return nil, ErrBadDomain
	}
	return &Gfycatdl{url: gfycatURL}, nil
}

// ScrapeVideoSource scrapes <source> tags from a gfycat web page
func (g *Gfycatdl) ScrapeVideoSource() (string, error) {
	resp, err := scrape(g.url)
	if err != nil {
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	var videoURL string
	doc.Find("source").Each(func(index int, item *goquery.Selection) {
		src := item.AttrOr("src", "")
		if src != "" {
			if !strings.Contains(src, "mobile") &&
				strings.Contains(src, "gfycat.com") &&
				strings.HasSuffix(src, ".mp4") {
				videoURL = src
			}
		}
	})
	if videoURL == "" {
		return "", ErrBadScrape
	}
	res, err := url.Parse(videoURL)
	if err != nil {
		return "", ErrBadScrape
	}
	g.ResourceName = res.Path[1:]
	return videoURL, nil
}

// DownloadFile downloads a file given a url pointing to a resource.
func DownloadFile(filepath, url string) error {
	client := &http.Client{Timeout: 1 * time.Minute}
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return nil
}

// scrape a url and returns the body of the response
func scrape(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
