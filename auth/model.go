package auth

import (
	"fmt"
	"time"

	"github.com/basebytes/tools"
	"github.com/dgrijalva/jwt-go"
)

const (
	defaultExpire = Expire8H
	Expire8H      = 28800
)

var defaultKey = fmt.Sprintf("data-%s", tools.CurrentDataCompactStr())

func NewConfig() *Config {
	return &Config{Key: defaultKey, Expire: defaultExpire}
}

type Config struct {
	Key    string
	Expire time.Duration //second
}

type MyCustomClaims struct {
	User
	jwt.StandardClaims
}

type User interface {
	Identity() string
	Token() string
	ExpiresInSeconds() int
	IsAdmin() bool
}
