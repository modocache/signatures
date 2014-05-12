package signatures

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

/*
Wrap the Martini server struct.
*/
type Server *martini.ClassicMartini

/*
Create a new *martini.ClassicMartini server.
We'll use a JSON renderer and our MongoDB
database handler. We define two routes:
"GET /signatures" and "POST /signatures".
*/
func NewServer(session *DatabaseSession) Server {
	// Create the server and set up middleware.
	m := Server(martini.Classic())
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Use(session.Database())

	// Define the "GET /signatures" route.
	m.Get("/signatures", func(r render.Render, db *mgo.Database) {
		r.JSON(200, fetchAllSignatures(db))
	})

	// Define the "POST /signatures" route.
	m.Post("/signatures", binding.Json(Signature{}),
		func(signature Signature,
			r render.Render,
			db *mgo.Database) {

			if signature.valid() {
				// signature is valid, insert into database
				err := db.C("signatures").Insert(signature)
				if err == nil {
					// insert successful, 201 Created
					r.JSON(201, signature)
				} else {
					// insert failed, 400 Bad Request
					r.JSON(400, map[string]string{
						"error": err.Error(),
					})
				}
			} else {
				// signature is invalid, 400 Bad Request
				r.JSON(400, map[string]string{
					"error": "Not a valid signature",
				})
			}
		})

	// Return the server. Call Run() on the server to
	// begin listening for HTTP requests.
	return m
}
