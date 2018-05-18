package mymongo1

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Dict struct {
	session    *mgo.Session
	database   string
	collection string
}

func NewDict(s *mgo.Session, d, c string) *Dict {
	return &Dict{
		session:    s,
		database:   d,
		collection: c,
	}
}

func (c *Dict) LookUp(a string) (string, error) {
	col := c.session.Clone().DB(c.database).C(c.collection)
	query := bson.M{"abbreviation": a}
	var result struct {
		Meaning string
	}
	if err := col.Find(query).One(&result); err != nil {
		return "", fmt.Errorf("looking for %q: %s", a, err)
	}
	return result.Meaning, nil
}

func (c *Dict) Add(a, m string) error {
	col := c.session.Clone().DB(c.database).C(c.collection)
	entry := struct {
		Abbreviation string
		Meaning      string
	}{
		Abbreviation: a,
		Meaning:      m,
	}
	if err := col.Insert(entry); err != nil {
		return fmt.Errorf("inserting %q: %s", entry, err)
	}
	return nil
}
