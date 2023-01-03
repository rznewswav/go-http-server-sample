package main

import (
	"newswav/http-server-sample/cmd/server"
)

func main() {
	shutdownHook := server.Prepare()
	defer shutdownHook()
	server.Start()
}
