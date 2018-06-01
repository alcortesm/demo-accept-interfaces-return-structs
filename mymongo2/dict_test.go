package mymongo2_test

import (
	"reflect"
	"testing"

	"github.com/alcortesm/demo-accept-interfaces-return-structs/mymongo2"
	"github.com/globalsign/mgo/bson"
)

// mocks to the github.com/globalsign/mgo types.
type mockSession struct {
	clone func() mymongo2.Session
	db    func(string) mymongo2.DataBase
}

func (m mockSession) Clone() mymongo2.Session       { return m.clone() }
func (m mockSession) DB(s string) mymongo2.DataBase { return m.db(s) }

type mockDataBase struct {
	c func(s string) mymongo2.Collection
}

func (m mockDataBase) C(col string) mymongo2.Collection { return m.c(col) }

type mockCollection struct {
	find   func(interface{}) mymongo2.Query
	insert func(docs ...interface{}) error
}

func (m mockCollection) Find(i interface{}) mymongo2.Query { return m.find(i) }
func (m mockCollection) Insert(docs ...interface{}) error  { return m.insert(docs...) }

type mockQuery struct {
	one func(result interface{}) error
}

func (m mockQuery) One(result interface{}) error { return m.one(result) }

func TestAdd(t *testing.T) {
	const (
		fixDB      = "test database"
		fixCol     = "test collection"
		fixAbbr    = "test abbreviation"
		fixMeaning = "test meaning"
	)

	// mock a mongoDB that checks that a new dictionary entry for fixAbbr and
	// fixMeaning is being added.
	var session mockSession
	{
		// checks that a dictionary entry (fixAbbr & fixMeaning) with
		// the correct format is being added.
		collection := mockCollection{
			insert: func(docs ...interface{}) error {
				if len(docs) != 1 {
					t.Fatalf("docs len was %d, want 1", len(docs))
				}
				want := bson.M{
					"abbr": fixAbbr,
					"data": fixMeaning,
				}
				if !equalBSON(t, want, docs[0]) {
					t.Fatalf("Different serialized values: want %#v, got %#v",
						want, docs[0])
				}
				return nil
			},
		}
		// checks that the requested collection is fixCol
		database := mockDataBase{
			c: func(s string) mymongo2.Collection {
				if s != fixCol {
					t.Fatalf("want %q, got %q", fixCol, s)
				}
				return collection
			},
		}
		// checks that the requested database is fixDB
		session = mockSession{
			clone: func() mymongo2.Session { return session },
			db: func(s string) mymongo2.DataBase {
				if s != fixDB {
					t.Fatalf("want %q, got %q", fixDB, s)
				}
				return database
			},
		}
	}

	dict := mymongo2.NewDict(session, fixDB, fixCol)
	if err := dict.Add(fixAbbr, fixMeaning); err != nil {
		t.Fatal(err)
	}
}

func equalBSON(t *testing.T, a, b interface{}) bool {
	t.Helper()
	aSerialized, err := bson.Marshal(a)
	if err != nil {
		t.Fatalf("serializing %#v: %v\n", a, err)
	}
	bSerialized, err := bson.Marshal(b)
	if err != nil {
		t.Fatalf("serializing %#v: %v\n", b, err)
	}
	return reflect.DeepEqual(aSerialized, bSerialized)
}

func TestLookUp(t *testing.T) {
	const (
		fixDB      = "test database"
		fixCol     = "test collection"
		fixAbbr    = "test abbreviation"
		fixMeaning = "test meaning"
	)
	fixEntry, err := bson.Marshal(bson.M{
		"_id":  bson.NewObjectId(),
		"abbr": fixAbbr,
		"data": fixMeaning,
	})
	if err != nil {
		t.Fatal(err)
	}

	// mock a mongoDB that checks that fixAbbr is being looked up and returns
	// fixEntry.
	var session mockSession
	{
		// mocks a query that returns fixEntry
		query := mockQuery{
			one: func(data interface{}) error {
				if err := bson.Unmarshal(fixEntry, data); err != nil {
					t.Fatal(err)
				}
				return nil
			},
		}
		// checks that fixAbbr is being requested
		collection := mockCollection{
			find: func(got interface{}) mymongo2.Query {
				want := bson.M{"abbr": fixAbbr}
				if !reflect.DeepEqual(want, got) {
					t.Fatalf("want %#v, got %#v", want, got)
				}
				return query
			},
		}
		// checks that the requested collection is fixCol
		database := mockDataBase{
			c: func(s string) mymongo2.Collection {
				if s != fixCol {
					t.Fatalf("want %q, got %q", fixCol, s)
				}
				return collection
			},
		}
		// checks that the requested database is fixDB
		session = mockSession{
			clone: func() mymongo2.Session { return session },
			db: func(s string) mymongo2.DataBase {
				if s != fixDB {
					t.Fatalf("want %q, got %q", fixDB, s)
				}
				return database
			},
		}
	}

	dict := mymongo2.NewDict(session, fixDB, fixCol)
	got, err := dict.LookUp(fixAbbr)
	if err != nil {
		t.Error(err)
	}
	if got != fixMeaning {
		t.Errorf("want %q, got %q", fixMeaning, got)
	}
}
