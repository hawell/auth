package database

type Config struct {
	ConnectionString string `env:"DB_CONNECTION_STRING" json:"connection_string"`
}
