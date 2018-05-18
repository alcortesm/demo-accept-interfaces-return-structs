package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/demo-accept-interfaces-return-structs/mymongo1"

	"github.com/globalsign/mgo"
)

const (
	url        = "localhost:27017/test"
	database   = "test"
	collection = "abbreviations"
)

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("error dialing mongo at %q: %s", url, err)
	}
	defer session.Close()
	dict := mymongo1.NewDict(session, database, collection)

	abbreviation := "NATO"
	meaning := "North Atlantic Treaty Organization"
	if err := dict.Add(abbreviation, meaning); err != nil {
		log.Fatalf("error adding meaning for %q: %s", abbreviation, err)
	}
	fmt.Printf("successfully added definition for %q\n", abbreviation)

	returned, err := dict.LookUp(abbreviation)
	if err != nil {
		log.Fatalf("error retrieving %q meaning: %s", abbreviation, err)
	}
	fmt.Printf("%q means %q\n", abbreviation, returned)
}
