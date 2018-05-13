package main

import (
	"fmt"
	"local/demo-accept-interfaces-return-structs/myhttp3"
	"log"
	"net/http"
	"time"
)

func main() {
	inner := &http.Client{
		Timeout: 10 * time.Second,
	}
	contentType := "application/json"
	c := myhttp3.NewClient(inner, contentType)

	data, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.25s...", data)
}
