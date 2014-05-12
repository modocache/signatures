package signatures

import "labix.org/v2/mgo"

/*
Each signature is composed of a first name, last name,
email, age, and short message. When represented in
JSON, ditch TitleCase for snake_case.
*/
type Signature struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Message   string `json:"message"`
}

/*
I want to make sure all these fields are present. The message
is optional, but if it's present it has to be less than
140 characters--it's a short blurb, not your life story.
*/
func (signature *Signature) valid() bool {
	return len(signature.FirstName) > 0 &&
		len(signature.LastName) > 0 &&
		len(signature.Email) > 0 &&
		signature.Age >= 18 && signature.Age <= 180 &&
		len(signature.Message) < 140
}

/*
I'll use this method when displaying all signatures for
"GET /signatures". Consult the mgo docs for more info:
http://godoc.org/labix.org/v2/mgo
*/
func fetchAllSignatures(db *mgo.Database) []Signature {
	signatures := []Signature{}
	err := db.C("signatures").Find(nil).All(&signatures)
	if err != nil {
		panic(err)
	}

	return signatures
}
