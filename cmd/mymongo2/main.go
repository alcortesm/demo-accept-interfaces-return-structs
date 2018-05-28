package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/demo-accept-interfaces-return-structs/mymongo2"

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
	dict := mymongo2.NewDict(mSession{session}, database, collection)

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

// Adaptors from Mongo types to mymongo types
type mSession struct{ *mgo.Session }
type mDataBase struct{ *mgo.Database }
type mCollection struct{ *mgo.Collection }
type mQuery struct{ *mgo.Query }

func (m mSession) Clone() mymongo2.Session {
	return mSession{m.Session.Clone()}
}

func (m mSession) DB(db string) mymongo2.DataBase {
	return mDataBase{m.Session.DB(db)}
}

func (m mDataBase) C(col string) mymongo2.Collection {
	return mCollection{m.Database.C(col)}
}

func (m mCollection) Find(i interface{}) mymongo2.Query {
	return mQuery{m.Collection.Find(i)}
}

func (m mCollection) Insert(docs ...interface{}) error {
	return m.Collection.Insert(docs...)
}

func (m mQuery) One(result interface{}) error {
	return m.Query.One(result)
}
