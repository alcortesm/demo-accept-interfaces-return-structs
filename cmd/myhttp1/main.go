package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alcortesm/demo-accept-interfaces-return-structs/myhttp1"
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
	fmt.Printf("%v\n", c)
}
