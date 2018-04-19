package ddpportal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
)

type JwtProvider struct {
}

func (j JwtProvider) CheckToken(jwtToken string) bool {
	tk, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("monsupersecretvraimentsupersecret"), nil
	})

	if err != nil {
		log.WithError(err).Error("Unable to parse token")
		return false
	}
	return tk.Valid
}

func NewJwtAuthMiddleware() func(http.Handler) http.Handler {
	jwtProvider := JwtProvider{}
	middleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			if bearer != "" {
				token := strings.Split(bearer, " ")[1]
				if token != "" && jwtProvider.CheckToken(token) {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "Not Authorized", 401)
				}
			} else {
				http.Error(w, "Not Authorized", 401)
			}
		}
		return http.HandlerFunc(fn)
	}

	return middleware
}
