package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var (
	jwtExpiration int
	signingKey    []byte
)

type JWTClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func init() {
	jwtExpiration = viper.GetInt("JWT.EXPIRE")
	signingKey = []byte(viper.GetString("JWT.PRIVATE_KEY"))
}
func GenerateToken(userId string) (string, error) {

	claims := JWTClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(jwtExpiration) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "Hypergro",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	return tokenString, err

}

func ValidateToken(encodedToken string) (string, error) {
	claims := &JWTClaims{}

	_, err := jwt.ParseWithClaims(encodedToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %s", t.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.UserId, nil
}
