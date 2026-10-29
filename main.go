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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		expirationWorker.Start(ctx)
	}()

	r := router.SetupRouter(productHandler, userHandler, authHandler, cartHandler, checkoutHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Http server error:", err)
		}
	}()
	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("HTTP server shutdown error:", err)
	}
	log.Println("HTTP server shutdown completed")

	wg.Wait()
	log.Println("Expiration worker stopped")
}
