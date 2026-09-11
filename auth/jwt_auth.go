package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/basebytes/tools"
	"github.com/dgrijalva/jwt-go"
)

func Init(authConfig *Config) {
	Reload(authConfig)
}

func Reload(authConfig *Config) {
	lock.Lock()
	defer lock.Unlock()
	key = []byte(authConfig.Key)
	expireDuration = time.Second * authConfig.Expire
}

func getKey() []byte {
	lock.RLock()
	defer lock.RUnlock()
	return key
}

func getExpires() time.Duration {
	lock.RLock()
	defer lock.RUnlock()
	return expireDuration
}

func GenerateToken(info User) (tokenString string, expireAt int64, err error) {
	// Create the Claims
	if info.Token() == "" {
		expireAt = tools.GetLocalTime().Add(getExpires()).Unix()
	} else if info.IsAdmin() {
		expireAt = tools.GetLocalTime().Add(longExpireDuration).Unix()
	} else {
		expireAt = tools.GetLocalTime().Add(time.Duration(info.ExpiresInSeconds()) * time.Second).Unix()
	}
	claims := MyCustomClaims{
		info,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    info.Identity(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getKey())
	return
}

func ValidateToken(tokenString string, info User) (err error) {
	var _token *jwt.Token
	_token, err = jwt.ParseWithClaims(tokenString, &MyCustomClaims{User: info}, keyFunc)
	if _token != nil {
		if claims, ok := _token.Claims.(*MyCustomClaims); ok && _token.Valid {
			info = claims.User
		}
	}
	return
}

func RefreshToken(oldTokenString string, exp int64) (tokenString string, expireAt int64, err error) {
	var _token *jwt.Token
	_token, err = jwt.ParseWithClaims(oldTokenString, &MyCustomClaims{}, keyFunc)
	if err == nil {
		if claims, ok := _token.Claims.(*MyCustomClaims); ok && _token.Valid {
			expireAt = tools.GetLocalTime().Add(time.Second * time.Duration(exp)).Unix()
			newClaims := MyCustomClaims{
				claims.User,
				jwt.StandardClaims{
					ExpiresAt: expireAt,
					Issuer:    claims.User.Identity(),
					IssuedAt:  time.Now().Unix(),
				},
			}
			tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims).SignedString(getKey())
		} else {
			err = errors.New("validate tokenString failed ")
		}
	}
	return
}

var (
	key                []byte
	expireDuration     time.Duration
	longExpireDuration = time.Hour * 24 * 7
	lock               sync.RWMutex
	keyFunc            = func(token *jwt.Token) (any, error) { return getKey(), nil }
)
