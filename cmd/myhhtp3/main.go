package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alcortesm/demo-accept-interfaces-return-structs/myhttp3"
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
