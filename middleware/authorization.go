package middleware

import (
	"go-ecommerce-api/response"
	"go-ecommerce-api/utils"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := utils.GetUserContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if user.Role != "admin" {
			response.Error(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
