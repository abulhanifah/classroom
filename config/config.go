package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config string

type DBConfig struct {
	DSN    string
	Driver string
}

var Default = map[string]Config{
	"APP_ENV":                  "local",
	"APP_PORT":                 "6191",
	"APP_URL":                  "http://declassroom.site",
	"APP_KEY":                  "3a1d279da65a7b813a2844516e26be5c",
	"API_IS_DEBUG":             "false",
	"OAUTH2_ACCESS_EXPIRE_IN":  "1h",   // 1 hour
	"OAUTH2_REFRESH_EXPIRE_IN": "720h", // 30 days
	"DB_DRIVER":                "mysql",
	"DB_HOST":                  "localhost",
	"DB_PORT":                  "3306",
	"DB_NAME":                  "declassroom",
	"DB_USER":                  "root",
	"DB_PASSWORD":              "root",
	"DB_SSL_MODE":              "disable",
	"DB_MAX_OPEN_CONNS":        "100",
	"DB_MAX_IDLE_CONNS":        "2",
	"DB_CONN_MAX_LIFETIME":     "0ms",
	"DB_IS_DEBUG":              "false",
}

func init() {
	godotenv.Load()
}

func GetConfig(key string) Config {
	value := Config(os.Getenv(key))
	if value == "" {
		value = Default[key]
	}
	return value
}

func (c Config) String() string {
	return string(c)
}

func (c Config) Int() int {
	v, err := strconv.Atoi(c.String())
	if err != nil {
		return 0
	}
	return v
}

func (c Config) Bool() bool {
	if strings.ToLower(c.String()) == "true" {
		return true
	}
	return false
}

func (c Config) Duration() time.Duration {
	v, err := time.ParseDuration(c.String())
	if err != nil {
		return 0
	}
	return v
}

func CreateDSN() (string, error) {
	host := GetConfig("DB_HOST").String()
	port := GetConfig("DB_PORT").String()
	name := GetConfig("DB_NAME").String()
	user := GetConfig("DB_USER").String()
	password := GetConfig("DB_PASSWORD").String()

	driver := GetConfig("DB_DRIVER").String()
	var dsn string
	if driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, name)
	} else if driver == "postgres" {
		ssl := GetConfig("DB_SSL_MODE").String()
		if ssl == "" && GetConfig("APP_ENV") == "local" {
			ssl = "disable"
		}
		if ssl != "" {
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, name, ssl)
		} else {
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
				host, port, user, password, name)
		}
	} else {
		return "", errors.New("Invalid driver")
	}
	return dsn, nil
}

func GetDBConfig() (DBConfig, error) {
	var dbconf DBConfig
	var err error
	dbconf.DSN, err = CreateDSN()
	dbconf.Driver = GetConfig("DB_DRIVER").String()
	return dbconf, err
}
