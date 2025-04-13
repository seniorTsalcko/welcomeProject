package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"welcomeProject/pkg/utils"
)

func JWTAuth(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "missing token")
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "invalid token format")
				return
			}

			var token, err = jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				utils.ErrorResponse(w, http.StatusUnauthorized, "invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
