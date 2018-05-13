package main

import (
	"fmt"
	"local/demo-accept-interfaces-return-structs/myhttp2"
	"log"
	"time"
)

func main() {
	timeout := 10 * time.Second
	contentType := "application/json"
	c := myhttp2.NewClient(timeout, contentType)

	data, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.25s...", data)
}
