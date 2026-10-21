package main

import (
	"context"
	"go-ecommerce-api/database"
	"go-ecommerce-api/handler"
	"go-ecommerce-api/repository"
	"go-ecommerce-api/router"
	"go-ecommerce-api/service"
	"go-ecommerce-api/worker"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(db)

	expirationService := service.NewOrderExpirationService(db, orderRepo, productRepo)
	expirationWorker := worker.NewOrderExpirationWorker(expirationService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	checkoutService := service.NewCheckoutService(db, cartRepo, productRepo, orderRepo, userRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go expirationWorker.Start(ctx)

	r := router.SetupRouter(productHandler, userHandler, authHandler, cartHandler, checkoutHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
