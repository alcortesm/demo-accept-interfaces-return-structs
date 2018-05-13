package main

import (
	"fmt"
	"local/demo-accept-interfaces-return-structs/myhttp1"
	"log"
	"time"
)

func main() {
	timeout := 10 * time.Second
	contentType := "application/json"
	c := myhttp1.NewClient(timeout, contentType)

	data, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.25s...\n", data)
}
