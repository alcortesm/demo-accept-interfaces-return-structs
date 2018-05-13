package main

import (
	"fmt"
	"io"
	"local/demo-accept-interfaces-return-structs/myhttp4"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//main1()
	main2()
}

func main1() {
	inner := &http.Client{
		Timeout: 10 * time.Second,
	}
	c := myhttp4.NewClient(inner, "application/json")

	data, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.25s...\n", data)
}

func main2() {
	verbose := &withLogs{
		inner: &http.Client{Timeout: 10 * time.Second},
		log:   log.New(os.Stdout, "my app: ", log.LstdFlags),
	}
	c := myhttp4.NewClient(verbose, "application/json")

	data, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.25s...\n", data)
}

type withLogs struct {
	inner myhttp4.GetPoster
	log   interface {
		Printf(format string, v ...interface{})
	}
}

func (c *withLogs) Get(url string) (resp *http.Response, err error) {
	c.log.Printf("GET(%q)", url)
	return c.inner.Get(url)
}

func (c *withLogs) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	c.log.Printf("POST(%q)", url)
	return c.inner.Post(url, contentType, body)
}
