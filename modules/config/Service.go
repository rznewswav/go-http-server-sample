package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	utils "newswav/http-server-sample/services/utils"

	"github.com/joho/godotenv"
)

type ConfigValues struct {
	App      AppConfig
	Auth     AuthConfig
	Database DatabaseConfig
}

type ConfigService struct {
	Config ConfigValues
}

func (service *ConfigService) Init() error {
	godotenv.Load()

	typeof := reflect.TypeOf(service.Config)
	configPointerToStruct := reflect.ValueOf(&(service.Config)).Elem()

	failed := false

	for i := 0; i < typeof.NumField(); i++ {
		// Get the configValue, returns https://golang.org/pkg/reflect/#StructField
		configValue := typeof.Field(i)
		configValueByName := configPointerToStruct.FieldByName(configValue.Name)

		for j := 0; j < configValue.Type.NumField(); j++ {
			configValueField := configValue.Type.Field(j)
			configValueFieldStruct := configValueByName.FieldByName(configValueField.Name)

			envTagName := "env"
			isRequiredTagName := "required"
			defaultValueTagName := "default"

			// Get the field tag value
			env := configValueField.Tag.Get(envTagName)
			isRequired := configValueField.Tag.Get(isRequiredTagName)
			defaultValue := configValueField.Tag.Get(defaultValueTagName)
			currentValue := os.Getenv(env)

			if len(env) == 0 {
				failed = true
				fmt.Fprintf(os.Stderr, "Trying to set required value for %s but from env is not set in struct tag!\n", configValue.Name+"."+configValueField.Name)
				continue
			}

			valueToSet := utils.If(len(currentValue) == 0, defaultValue, currentValue)

			if len(valueToSet) == 0 && isRequired == "true" {
				failed = true
				fmt.Fprintf(os.Stderr, "Trying to set required value for %s from env %s but value is not set in .env!\n", configValue.Name+"."+configValueField.Name, env)
				continue
			}

			if configValueFieldStruct.CanSet() {
				if configValueFieldStruct.Kind() == reflect.Int64 {
					if x, err := strconv.ParseInt(valueToSet, 10, 64); err != nil {
						failed = true
						fmt.Fprintf(
							os.Stderr,
							"Trying to set value for %s with %s configured from %s but field cannot be set: Expected an Int64 but got %s instead!\n",
							configValue.Name+"."+configValueField.Name,
							valueToSet,
							utils.If(len(currentValue) == 0, "default value", ".env"),
							configValueFieldStruct.Kind().String(),
						)
					} else if !configValueFieldStruct.OverflowInt(x) {
						configValueFieldStruct.SetInt(x)
					}
				} else if configValueFieldStruct.Kind() == reflect.String {
					configValueFieldStruct.SetString(valueToSet)
				} else if configValueFieldStruct.Kind() == reflect.Bool {
					configValueFieldStruct.SetBool(valueToSet == "true")
				} else {
					failed = true
					fmt.Fprintf(
						os.Stderr,
						"Trying to set value for %s with %s configured from %s but field cannot be set: Unknown type %s!\n",
						configValue.Name+"."+configValueField.Name,
						valueToSet,
						utils.If(len(currentValue) == 0, "default value", ".env"),
						configValueFieldStruct.Kind().String(),
					)
				}
			} else {
				fmt.Fprintf(
					os.Stderr,
					"Trying to set value for %s with %s configured from %s but field cannot be set!\n",
					configValue.Name+"."+configValueField.Name,
					valueToSet,
					utils.If(len(currentValue) == 0, "default value", ".env"),
				)
			}

		}

	}

	if failed {
		return errors.New("failed to completely configure app configuration from .env")
	} else {
		return nil
	}
}
