package mymongo1

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Dic is a dictionary of abbreviations and their meanings.
type Dict struct {
	session    *mgo.Session
	database   string
	collection string
}

func NewDict(session *mgo.Session, db, col string) *Dict {
	return &Dict{
		session:    session,
		database:   db,
		collection: col,
	}
}

// Entry is the internal mongo schema for dictionary entries.
type Entry struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Abbreviation string        `bson:"abbr"`
	Meaning      string        `bson:"data"`
}

func (d *Dict) LookUp(a string) (string, error) {
	col := d.session.Clone().DB(d.database).C(d.collection)
	query := bson.M{"abbr": a}
	var result Entry
	if err := col.Find(query).One(&result); err != nil {
		return "", fmt.Errorf("looking for %q: %s", a, err)
	}
	return result.Meaning, nil
}

func (d *Dict) Add(a, m string) error {
	col := d.session.Clone().DB(d.database).C(d.collection)
	doc := Entry{
		Abbreviation: a,
		Meaning:      m,
	}
	if err := col.Insert(doc); err != nil {
		return fmt.Errorf("inserting %q: %s", doc, err)
	}
	return nil
}
