package configuration

type Config struct {
	DatabaseConnectionURL string
	HTTPAddress           string
}

func New(databaseConnectionURL, httpAddress string) *Config {
	return &Config{
		DatabaseConnectionURL: databaseConnectionURL,
		HTTPAddress:           httpAddress,
	}
}
