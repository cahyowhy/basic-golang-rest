package config

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

type Config struct {
	DB *DBConfig
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Username: "root",
			Password: "root",
			Name:     "todoapp",
			Charset:  "utf8",
		},
	}
}