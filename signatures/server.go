package signatures

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

type Server *martini.ClassicMartini

func NewServer(session *DatabaseSession) Server {
	m := Server(martini.Classic())
	m.Use(render.Renderer(render.Options{IndentJSON: true}))
	m.Use(session.Database())

	m.Get("/signatures", func(r render.Render, db *mgo.Database) {
		r.JSON(200, fetchAllSignatures(db))
	})

	m.Post("/signatures", binding.Json(Signature{}),
		func(signature Signature, r render.Render, db *mgo.Database) {
			if signature.valid(db) {
				// signature is valid, insert into database
				err := db.C("signatures").Insert(signature)
				if err == nil {
					// insert successful, 201 Created
					r.JSON(201, signature)
				} else {
					// insert failed, 400 Bad Request
					r.JSON(400, map[string]string{"error": err.Error()})
				}
			} else {
				// signature is invalid, 400 Bad Request
				r.JSON(400, map[string]string{"error": "Not a valid signature"})
			}
		})

	return m
}
