package signatures

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

func NewSession(name string) *DatabaseSession {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	addIndexToSignatureEmails(session.DB(name))
	return &DatabaseSession{session, name}
}

func (session *DatabaseSession) Database() martini.Handler {
	return func(context martini.Context) {
		s := session.Clone()
		context.Map(s.DB(session.databaseName))
		defer s.Close()
		context.Next()
	}
}

func addIndexToSignatureEmails(db *mgo.Database) {
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	indexErr := db.C("signatures").EnsureIndex(index)
	if indexErr != nil {
		panic(indexErr)
	}
}
