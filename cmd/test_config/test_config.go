package main

import (
	"encoding/json"
	config "newswav/http-server-sample/modules/config"
)

func main() {
	service := config.ConfigService{}
	service.Init()

	println("Available env:")

	if b, err := json.MarshalIndent(service.Config, "", "  "); err != nil {
		println("Unable to parse config into JSON:", err.Error())
	} else {
		println(string(b))
	}

}
