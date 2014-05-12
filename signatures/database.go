package signatures

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

/*
I want to use a different database for my tests,
so I'll embed *mgo.Session and store the database name.
*/
type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

/*
Connect to the local MongoDB and set up the database.
*/
func NewSession(name string) *DatabaseSession {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	addIndexToSignatureEmails(session.DB(name))
	return &DatabaseSession{session, name}
}

/*
Add a unique index on the "email" field.
This doesn't prevent users from signing twice,
since they can still enter
"dudebro+signature2@exmaple.com". But if they're
that clever, I say they deserve the extra signature.
*/
func addIndexToSignatureEmails(db *mgo.Database) {
	index := mgo.Index{
		Key:      []string{"email"},
		Unique:   true,
		DropDups: true,
	}
	indexErr := db.C("signatures").EnsureIndex(index)
	if indexErr != nil {
		panic(indexErr)
	}
}

/*
Martini lets you inject parameters for routing handlers
by using `context.Map()`. I'll pass each route handler
a instance of a *mgo.Database, so they can retrieve
and insert signatures to and from that database.

For more information, check out:
http://blog.gopheracademy.com/day-11-martini
*/
func (session *DatabaseSession) Database() martini.Handler {
	return func(context martini.Context) {
		s := session.Clone()
		context.Map(s.DB(session.databaseName))
		defer s.Close()
		context.Next()
	}
}
