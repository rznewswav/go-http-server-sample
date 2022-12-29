package main

import (
	"newswav/http-server-sample/cmd/server"

	"github.com/ztrue/shutdown"
)

func main() {
	shutdownHook := server.Prepare()

	shutdown.AddWithKey("server hook", func() {
		(*shutdownHook)()
	})

	go func() {
		server.Start()
	}()

	shutdown.Listen()
}
