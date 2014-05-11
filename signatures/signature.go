package signatures

import "labix.org/v2/mgo"

type Signature struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Message   string `json:"message"`
}

func (signature *Signature) valid(db *mgo.Database) bool {
	return len(signature.FirstName) > 0 &&
		len(signature.LastName) > 0 &&
		len(signature.Email) > 0 &&
		signature.Age > 21 && signature.Age < 180 &&
		len(signature.Message) < 140
}

func fetchAllSignatures(db *mgo.Database) []Signature {
	var signatures []Signature
	db.C("signatures").Find(nil).All(&signatures)
	return signatures
}
