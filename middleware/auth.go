package middleware

import (
	"go-ecommerce-api/response"
	"go-ecommerce-api/utils"
	"net/http"
	"strings"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if result := !strings.HasPrefix(bearer, "Bearer "); result {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		tokenString := strings.TrimPrefix(bearer, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx := utils.SetUserID(r.Context(), claims.Subject)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		return
	})
}
