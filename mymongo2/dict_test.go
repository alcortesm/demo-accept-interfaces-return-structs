package mymongo2

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
)

type mockSession struct {
	clone func() Session
	db    func(string) DataBase
}

func (m mockSession) Clone() Session       { return m.clone() }
func (m mockSession) DB(s string) DataBase { return m.db(s) }

type mockDataBase struct {
	c func(s string) Collection
}

func (m mockDataBase) C(col string) Collection { return m.c(col) }

type mockCollection struct {
	find   func(interface{}) Query
	insert func(docs ...interface{}) error
}

func (m mockCollection) Find(i interface{}) Query         { return m.find(i) }
func (m mockCollection) Insert(docs ...interface{}) error { return m.insert(docs...) }

type mockQuery struct {
	one func(result interface{}) error
}

func (m mockQuery) One(result interface{}) error { return m.one(result) }

func TestAddThenLookUp(t *testing.T) {
	const (
		fixDB      = "test database"
		fixCol     = "test collection"
		fixID      = "test objectID"
		fixAbbr    = "test abbreviation"
		fixMeaning = "test meaning"
	)

	// mock a mongoDB that checks that fixAbbr and fixMeaning is added
	// and that fixAbbr is looked up.
	var session mockSession
	var database mockDataBase
	var collection mockCollection
	var query mockQuery
	{
		session = mockSession{
			clone: func() Session { return session },
			// checks that the database is fixDB
			db: func(s string) DataBase {
				if s != fixDB {
					t.Fatalf("want %q, got %q", fixDB, s)
				}
				return database
			},
		}
		database = mockDataBase{
			// checks that the collection is fixCol
			c: func(s string) Collection {
				if s != fixCol {
					t.Fatalf("want %q, got %q", fixCol, s)
				}
				return collection
			},
		}
		collection = mockCollection{
			// checks that fixAbbr is being requested
			find: func(q interface{}) Query {
				want := bson.M{"abbr": fixAbbr}
				if !reflect.DeepEqual(want, q) {
					t.Fatalf("want %#v, got %#v", want, q)
				}
				return query
			},
			// checks that fixAbbr & fixMeaning is being added
			insert: func(docs ...interface{}) error {
				if len(docs) != 1 {
					t.Fatalf("docs len was %d, want 1", len(docs))
				}
				want := Entry{
					Abbreviation: fixAbbr,
					Meaning:      fixMeaning,
				}
				if !reflect.DeepEqual(want, docs[0]) {
					t.Fatalf("want %#v, got %#v", want, docs[0])
				}
				return nil
			},
		}
		query = mockQuery{
			// mocks a query that returns fixAbbr and fixMeaning
			one: func(data interface{}) error {
				ret, ok := data.(*Entry)
				if !ok {
					t.Fatal("wrong data type: %T", data)
				}
				ret.ID = bson.ObjectId(fixID)
				ret.Abbreviation = fixAbbr
				ret.Meaning = fixMeaning
				return nil
			},
		}
	}

	dict := NewDict(session, fixDB, fixCol)
	if err := dict.Add(fixAbbr, fixMeaning); err != nil {
		t.Fatal(err)
	}
	got, err := dict.LookUp(fixAbbr)
	if err != nil {
		t.Error(err)
	}
	if got != fixMeaning {
		t.Errorf("want %q, got %q", fixMeaning, got)
	}
}
