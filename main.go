package main

import "github.com/modocache/signatures/signatures"

func main() {
	session := signatures.NewSession("signatures")
	server := signatures.NewServer(session)
	server.Run()
}
