package router

import (
	"go-ecommerce-api/handler"
	"go-ecommerce-api/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(productHandler *handler.ProductHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProductByID)
	r.Post("/products", productHandler.CreateProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)
	return r
}
