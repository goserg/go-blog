package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/goserg/microblog/models"
)

func Hash(bytes []byte) string {
	h := sha1.New()
	h.Write(bytes)
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	IPAddress = strings.Split(IPAddress, ":")[0]
	return IPAddress
}

var tokenKey = []byte(os.Getenv("TOKEN_PASS"))

type Claims struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	jwt.StandardClaims
}

func GenerateNewToken(user models.User) string {
	claims := &Claims{
		Username:     user.UserName,
		PasswordHash: user.PasswordHash,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(tokenKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return ""
	}
	return tokenString
}

func ParseJWT(tokenString string) (string, string, int64, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", 0, jwt.ValidationError{}
		}
		return "", "", 0, jwt.ValidationError{}
	}
	if !tkn.Valid {
		return "", "", 0, jwt.ValidationError{}
	}
	return claims.Username, claims.PasswordHash, claims.StandardClaims.ExpiresAt, nil
}
