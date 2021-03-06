package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(userId uint) (string, error) {

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "aaa.com",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {

		return jwtkey, nil
	})
	return token, claims, err
}
