package util

import (
	"fmt"
	"os"
	"time"
)

// DbConfig is an interface to manage Db configuration
type DbConfig interface {
	GetConnect() string
	GetType() string
}

// SqliteConfig is a Db config for SQLite
type SqliteConfig struct {
	DbFilePath string
}

// PostgresConfig is a Db config for PostgreSQL
type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     string
}

// Config contains all the config extractable from env
type Config struct {
	Db             DbConfig
	Port           string
	SecretPhrase   string
	Validity       time.Duration
	RefeshValidity time.Duration
	IdentityURL    string
	AttachmentURL  string
	IconURL        string
	StaticFilePath string
}

// InitConfig initialize a new Config object
func InitConfig() *Config {
	typeDb := getEnv("DB_TYPE", "sqlite")

	config := &Config{
		Port:           getEnv("PORT", "3000"),
		Validity:       time.Hour,
		RefeshValidity: time.Hour + 30*time.Minute,
		IdentityURL:    getEnv("WARDEN_IDENTITY_URL", "/identity"),
		AttachmentURL:  getEnv("WARDEN_ATTACHMENT_URL", "/attachments"),
		IconURL:        getEnv("WARDEN_ICONS_URL", "/icons"),
		SecretPhrase:   getEnv("WARDEN_SECRET_PHRASE", "This a secret ... sshhhshh"),
		StaticFilePath: getEnv("WARDEN_STATIC_PATH", "./fixtures/assets"),
	}
	if typeDb == "postgres" {
		// Postgres
		config.Db = PostgresConfig{
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Name:     getEnv("DB_NAME", ""),
			Port:     getEnv("DB_PORT", "5432"),
		}
	}

	// Default type is sqlite
	config.Db = SqliteConfig{
		DbFilePath: getEnv("DB_FILEPATH", "./fixtures/test.db"),
	}
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetConnect provide the Url for PostgreSQL
func (conf PostgresConfig) GetConnect() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Name)
}

//GetType for Posgres
func (conf PostgresConfig) GetType() string {
	return "postgres"
}

//GetConnect provide the filename for SQLite
func (conf SqliteConfig) GetConnect() string {
	return conf.DbFilePath
}

//GetType for SQLite
func (conf SqliteConfig) GetType() string {
	return "sqlite3"
}
