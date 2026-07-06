package router

import (
	"go-ecommerce-api/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(productHandler *handler.ProductHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProductByID)
	return r
}
