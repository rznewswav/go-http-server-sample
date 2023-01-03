package server

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"newswav/http-server-sample/modules/auth"
	"newswav/http-server-sample/modules/config"
	"newswav/http-server-sample/modules/mongodb"
)

var depManager = dig.New()
var router = gin.Default()

func Prepare() func() {
	depManager.Provide(func() (*config.ConfigModule, error) {
		module := config.ConfigModule{}
		module.Service = &config.ConfigService{}
		err := module.Service.Init()
		return &module, err
	})

	depManager.Provide(func(config *config.ConfigModule) (*mongodb.MongodbModule, error) {
		module := mongodb.MongodbModule{}
		module.Service = &mongodb.MongodbService{}
		module.Service.Init(config.Service.Config.Database.MongoURI, config.Service.Config.Database.MongoDbName)
		return &module, nil
	})

	depManager.Provide(func(config *config.ConfigModule, database *mongodb.MongodbModule) (*auth.AuthModule, error) {
		module := auth.AuthModule{}
		module.Controller = &auth.AuthController{}
		service := &auth.AuthService{}
		service.Database = database.Service
		service.JWTSecret = config.Service.Config.Auth.JwtSecret

		module.Service = service
		module.Controller.Service = module.Service

		return &module, nil
	})

	router.SetTrustedProxies([]string{"localhost"})

	// Auth route
	depManager.Invoke(func(authModule *auth.AuthModule) error {
		router.GET("api/v1/auth", authModule.Controller.GetAuthenticateToken)
		router.POST("api/v1/auth", authModule.Controller.PostGenerateAuthToken)
		return nil
	})

	shutdownHook := func() {
		fmt.Println("Closing database...")
		err := depManager.Invoke(func(database *mongodb.MongodbModule) error {
			err := database.Service.Client.Disconnect(context.Background())
			return err
		})
		if err == nil {
			fmt.Println("Database is closed.")
		} else {
			fmt.Fprintln(os.Stderr, "Unable to close database:", err.Error())
		}
		fmt.Println("Server is closed.")
	}
	return shutdownHook
}

func Start() {
	depManager.Invoke(func(config *config.ConfigModule) error {
		err := router.Run(":" + config.Service.Config.App.AppPort)
		return err
	})
}
