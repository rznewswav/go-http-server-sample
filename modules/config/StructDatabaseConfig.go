package config

type DatabaseConfig struct {
	MongoURI    string `env:"MONGO_URI" required:"true"`
	MongoDbName string `env:"MONGO_DB" default:"golang-poc"`
}
