package signatures_test

import (
	"github.com/modocache/signatures/signatures"

	"bytes"
	"encoding/json"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Signatures", func() {
	var databaseName string
	var session *signatures.DatabaseSession
	var server signatures.Server
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		databaseName = "signatures_test"
		session = signatures.NewSession(databaseName)
		server = signatures.NewServer(session)
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		session.DB(databaseName).DropDatabase()
	})

	Describe("GET /signatures", func() {
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
				Expect(recorder.Body.String()).To(Equal("null"))
			})
		})

		Context("when signatures exist", func() {
			BeforeEach(func() {
				signatures := session.DB(databaseName).C("signatures")
				signatures.Insert(gory.Build("signature"))
				signatures.Insert(gory.Build("signature"))
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those signatures in the body", func() {
				server.ServeHTTP(recorder, request)

				var body interface{}
				json.Unmarshal(recorder.Body.Bytes(), &body)
				peopleJSON := body.([]interface{})
				Expect(len(peopleJSON)).To(Equal(2))

				personJSON := peopleJSON[0].(map[string]interface{})
				Expect(personJSON["first_name"]).To(Equal("Jane"))
				Expect(personJSON["last_name"]).To(Equal("Doe"))
				Expect(personJSON["age"]).To(Equal(float64(27)))
				Expect(personJSON["message"]).To(Equal("I wholeheartedly support this petition!"))
				Expect(personJSON["email"]).To(ContainSubstring("jane-doe"))
			})
		})
	})

	Describe("POST /signatures", func() {
		Context("with invalid JSON", func() {
			BeforeEach(func() {
				body, _ := json.Marshal(`{ "first_name": "Bill" }`)
				request, _ = http.NewRequest("POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with valid JSON", func() {
			BeforeEach(func() {
				body, _ := json.Marshal(gory.Build("signature"))
				request, _ = http.NewRequest("POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(201))
			})

			It("returns the inserted signature", func() {
				server.ServeHTTP(recorder, request)

				var body interface{}
				json.Unmarshal(recorder.Body.Bytes(), &body)
				personJSON := body.(map[string]interface{})

				Expect(personJSON["first_name"]).To(Equal("Jane"))
				Expect(personJSON["last_name"]).To(Equal("Doe"))
				Expect(personJSON["age"]).To(Equal(float64(27)))
				Expect(personJSON["message"]).To(Equal("I wholeheartedly support this petition!"))
				Expect(personJSON["email"]).To(ContainSubstring("jane-doe"))
			})
		})

		Context("with JSON containing a duplicate email", func() {
			BeforeEach(func() {
				signature := gory.Build("signature").(*signatures.Signature)
				session.DB(databaseName).C("signatures").Insert(signature)

				body, _ := json.Marshal(signature)
				request, _ = http.NewRequest("POST", "/signatures", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})
})
