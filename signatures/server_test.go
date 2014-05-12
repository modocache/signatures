package signatures_test

import (
	. "github.com/modocache/signatures/signatures"

	"bytes"
	"encoding/json"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

/*
Convert JSON data into a slice.
*/
func sliceFromJSON(data []byte) []interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.([]interface{})
}

/*
Convert JSON data into a map.
*/
func mapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.(map[string]interface{})
}

/*
Server unit tests.
*/
var _ = Describe("Server", func() {
	var dbName string
	var session *DatabaseSession
	var server Server
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up a new server, connected to a test database,
		// before each test.
		dbName = "signatures_test"
		session = NewSession(dbName)
		server = NewServer(session)

		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// Clear the database after each test.
		session.DB(dbName).DropDatabase()
	})

	Describe("GET /signatures", func() {

		// Set up a new GET request before every test
		// in this describe block.
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/signatures", nil)
		})

		Context("when no signatures exist", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns a null body", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal("[]"))
			})
		})

		Context("when signatures exist", func() {

			// Insert two valid signatures into the database
			// before each test in this context.
			BeforeEach(func() {
				collection := session.DB(dbName).C("signatures")
				collection.Insert(gory.Build("signature"))
				collection.Insert(gory.Build("signature"))
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those signatures in the body", func() {
				server.ServeHTTP(recorder, request)

				peopleJSON := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(peopleJSON)).To(Equal(2))

				personJSON := peopleJSON[0].(map[string]interface{})
				Expect(personJSON["first_name"]).To(Equal("Jane"))
				Expect(personJSON["last_name"]).To(Equal("Doe"))
				Expect(personJSON["age"]).To(Equal(float64(27)))
				Expect(personJSON["message"]).To(Equal("I agree!"))
				Expect(personJSON["email"]).To(
					ContainSubstring("jane-doe"))
			})
		})
	})

	Describe("POST /signatures", func() {

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("signatureTooYoung"))
				request, _ = http.NewRequest(
					"POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with valid JSON", func() {

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("signature"))
				request, _ = http.NewRequest(
					"POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(201))
			})

			It("returns the inserted signature", func() {
				server.ServeHTTP(recorder, request)

				personJSON := mapFromJSON(recorder.Body.Bytes())
				Expect(personJSON["first_name"]).To(Equal("Jane"))
				Expect(personJSON["last_name"]).To(Equal("Doe"))
				Expect(personJSON["age"]).To(Equal(float64(27)))
				Expect(personJSON["message"]).To(Equal("I agree!"))
				Expect(personJSON["email"]).To(
					ContainSubstring("jane-doe"))
			})
		})

		Context("with JSON containing a duplicate email", func() {
			BeforeEach(func() {
				signature := gory.Build("signature")
				session.DB(dbName).C("signatures").Insert(signature)

				body, _ := json.Marshal(signature)
				request, _ = http.NewRequest(
					"POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})
})
