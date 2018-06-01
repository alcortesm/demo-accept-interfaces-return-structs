package mymongo2

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// Dict is a dictionary of abbreviations and their meanings.
type Dict struct {
	session    Session
	database   string
	collection string
}

// Interfaces to github.com/globalsign/mgo types.
type Session interface {
	Clone() Session
	DB(string) DataBase
}
type DataBase interface {
	C(string) Collection
}
type Collection interface {
	Find(interface{}) Query
	Insert(docs ...interface{}) error
}
type Query interface {
	One(result interface{}) error
}

// NewDict returns a new dictionary ready to use.
func NewDict(session Session, db, col string) *Dict {
	return &Dict{
		session:    session,
		database:   db,
		collection: col,
	}
}

// entry is the internal mongo schema for dictionary entries.
type entry struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Abbreviation string        `bson:"abbr"`
	Meaning      string        `bson:"data"`
}

// Add adds an entry to the dictionary, using abbreviation a and meaning m.
func (c *Dict) Add(a, m string) error {
	col := c.session.Clone().DB(c.database).C(c.collection)
	doc := entry{
		Abbreviation: a,
		Meaning:      m,
	}
	if err := col.Insert(doc); err != nil {
		return fmt.Errorf("inserting %q: %s", doc, err)
	}
	return nil
}

// LookUp returns the meaning for the abbreviation a.
func (c *Dict) LookUp(a string) (string, error) {
	col := c.session.Clone().DB(c.database).C(c.collection)
	query := bson.M{"abbr": a}
	var result entry
	if err := col.Find(query).One(&result); err != nil {
		return "", fmt.Errorf("looking for %q: %s", a, err)
	}
	return result.Meaning, nil
}
