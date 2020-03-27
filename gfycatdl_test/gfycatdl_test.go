package gfycatdl_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/umahmood/gfycatdl"
)

func readHTML(html string) (string, error) {
	file, err := os.Open(html)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func startTestServer() *httptest.Server {
	s := make(chan struct{})
	var ts *httptest.Server

	go func() {
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path[1:]
			html, err := readHTML(path)
			if err != nil {
				fmt.Fprintf(w, "error accessing "+path)
			} else {
				fmt.Fprintf(w, html)
			}
		}))
		s <- struct{}{}
	}()
	_ = <-s
	return ts
}

func TestScrapeVideoURL(t *testing.T) {
	server := startTestServer()
	defer server.Close()
	g, err := gfycatdl.New(server.URL + "/test_valid_video_source.html")
	if err != nil {
		t.Errorf("%s", err)
	}
	scrapedURL, err := g.ScrapeVideoSource()
	if err != nil {
		t.Errorf("%s", err)
	}
	wantURL := "https://zippy.gfycat.com/AcademicPrestigiousGrasshopper.mp4"
	if scrapedURL != wantURL {
		t.Errorf("fail got %s want %s", scrapedURL, wantURL)
	}
	wantName := "AcademicPrestigiousGrasshopper.mp4"
	if g.ResourceName != wantName {
		t.Errorf("fail got name %s want name %s", g.ResourceName, wantName)
	}
}

func TestScrapeVideoURLNotFound(t *testing.T) {
	server := startTestServer()
	defer server.Close()
	g, err := gfycatdl.New(server.URL + "/test_no_valid_video_source.html")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = g.ScrapeVideoSource()
	if err != gfycatdl.ErrBadScrape {
		t.Errorf("fail got err %s want err %s", err, gfycatdl.ErrBadScrape)
	}
	if g.ResourceName != "" {
		t.Errorf("fail got name %s - field should be empty.", g.ResourceName)
	}
}

func TestInvalidURL(t *testing.T) {
	_, err := gfycatdl.New("http://foo.com/bar")
	if err != gfycatdl.ErrBadDomain {
		t.Errorf("fail got err %s want err %s", err, gfycatdl.ErrBadDomain)
	}
}
