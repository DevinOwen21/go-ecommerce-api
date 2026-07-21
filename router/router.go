package router

import (
	"go-ecommerce-api/handler"
	"go-ecommerce-api/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(productHandler *handler.ProductHandler, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProductByID)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)

		r.Get("/users/me", userHandler.GetProfile)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Use(middleware.AdminOnly)

		r.Post("/products", productHandler.CreateProduct)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProduct)
	})
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	return r
}
