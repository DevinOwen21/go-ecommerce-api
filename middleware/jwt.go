package middleware

import (
	"go-ecommerce-api/response"
	"go-ecommerce-api/utils"
	"net/http"
	"strings"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if result := !strings.HasPrefix(authHeader, "Bearer "); result {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		user := utils.UserContext{
			UserID: claims.Subject,
			Role:   claims.Role,
		}
		ctx := utils.SetUserContext(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		return
	})
}
