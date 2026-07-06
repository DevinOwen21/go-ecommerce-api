package main

import (
	"go-ecommerce-api/database"
	"go-ecommerce-api/handler"
	"go-ecommerce-api/repository"
	"go-ecommerce-api/router"
	"go-ecommerce-api/service"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewProductRepository(db)
	productService := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productService)
	r := router.SetupRouter(productHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
