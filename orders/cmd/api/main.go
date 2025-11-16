package main

import (
	"log/slog"
	"net/http"
	"orders-go/internal/infra/database"
	grpcClient "orders-go/internal/infra/grpc/client"
	"orders-go/internal/infra/http/handler"
	infraRepo "orders-go/internal/infra/repository"
	"orders-go/internal/usecase"
	"os"
	"time"

	_ "orders-go/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Orders API
// @version 1.0
// @description API robusta de gerenciamento de pedidos com carrinho de compras
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@orders-api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Orders API")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Connect to database
	db, err := database.NewMySQLConnection()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connected successfully")

	// Connect to Payment Service via gRPC
	paymentServiceAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentServiceAddr == "" {
		paymentServiceAddr = "localhost:50051"
	}

	paymentClient, err := grpcClient.NewPaymentClient(paymentServiceAddr, logger)
	if err != nil {
		slog.Error("Failed to connect to payment service", "error", err)
		os.Exit(1)
	}
	defer paymentClient.Close()
	slog.Info("Connected to payment service successfully", "addr", paymentServiceAddr)

	// Initialize repositories
	productRepo := infraRepo.NewProductRepository(db, logger)
	orderRepo := infraRepo.NewOrderRepository(db, logger)
	itemRepo := infraRepo.NewItemRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, logger)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, logger)
	cartUseCase := usecase.NewCartUseCase(orderRepo, productRepo, logger)
	createOrderWithPaymentUseCase := usecase.NewCreateOrderUseCase(orderRepo, itemRepo, productRepo, paymentClient, logger)
	cancelOrderUseCase := usecase.NewCancelOrderUseCase(orderRepo, paymentClient, logger)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase, logger)
	orderHandler := handler.NewOrderHandler(orderUseCase, logger)
	cartHandler := handler.NewCartHandler(cartUseCase, logger)
	orderWithPaymentHandler := handler.NewOrderWithPaymentHandler(createOrderWithPaymentUseCase, cancelOrderUseCase, logger)

	// Setup router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.List)
			r.Post("/", productHandler.Create)
			r.Get("/{id}", productHandler.GetByID)
			r.Put("/{id}", productHandler.Update)
			r.Delete("/{id}", productHandler.Delete)
		})

		// Order routes
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", orderHandler.List)
			r.Get("/{id}", orderHandler.GetByID)
			r.Delete("/{id}", orderHandler.Delete)

			// Order with payment integration
			r.Post("/with-payment", orderWithPaymentHandler.CreateOrderWithPayment)
			r.Post("/{id}/cancel", orderWithPaymentHandler.CancelOrder)
		})

		// Cart routes
		r.Route("/cart", func(r chi.Router) {
			r.Post("/", cartHandler.CreateCart)
			r.Get("/{id}", cartHandler.GetCart)
			r.Post("/{id}/items", cartHandler.AddItem)
			r.Delete("/{id}/items/{itemId}", cartHandler.RemoveItem)
			r.Put("/{id}/items/{itemId}", cartHandler.UpdateItemQuantity)
			r.Get("/{id}/calculate", cartHandler.CalculateTotal)
			r.Put("/{id}/status", cartHandler.UpdateStatus)
		})
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
