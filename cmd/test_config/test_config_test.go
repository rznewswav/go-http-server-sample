package main_test

import (
	config "newswav/http-server-sample/modules/config"
	"testing"
)

func TestService(t *testing.T) {
	service := config.ConfigService{}
	err := service.Init()
	if err != nil {
		t.Fail()
	}
}
