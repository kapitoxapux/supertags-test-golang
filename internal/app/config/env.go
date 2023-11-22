package config

import (
	"net/http"
	"time"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
	DB      string `env:"DATABASE_DSN"`
	Cookie  *http.Cookie
}

var ServerConfig Config

const serverAdress = "localhost:8080"
const baseURL = "http://localhost:8080"
const storageDB = "host= port= user= password= dbname=supertags sslmode=disable"

var cookie = http.Cookie{
	Name:    "SESSION",
	Value:   "",
	Path:    "/",
	Expires: time.Now(),
}

func SetConfig() Config {
	ServerConfig.Address = serverAdress
	ServerConfig.Base = baseURL
	ServerConfig.DB = storageDB
	ServerConfig.Cookie = &cookie

	return ServerConfig
}

func GetAddress() string {

	return ServerConfig.Address
}

func GetDB() string {

	return ServerConfig.DB
}

func GetCookie() *http.Cookie {

	return ServerConfig.Cookie
}

func SetCookieConfig(cookie string) *http.Cookie {

	ServerConfig.Cookie.Value = cookie

	return ServerConfig.Cookie
}
