package main

import "github.com/modocache/signatures/signatures"

/*
Create a new MongoDB session, using a database
named "signatures". Create a new server using
that session, then begin listening for HTTP requests.
*/
func main() {
	session := signatures.NewSession("signatures")
	server := signatures.NewServer(session)
	server.Run()
}
