package main_test

import (
	"newswav/http-server-sample/cmd/server"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	os.Setenv("JWT_SECRET", "abcdef")
	os.Setenv("MONGO_URI", "mongodb://127.0.0.1:27017/")
	server.Prepare()
	server.Start()
}
