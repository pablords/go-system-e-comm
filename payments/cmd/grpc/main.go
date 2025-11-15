package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"payments-go/internal/infra/database"
	grpcHandler "payments-go/internal/infra/grpc/handler"
	"payments-go/internal/infra/repository"
	"payments-go/internal/usecase"
	pb "payments-go/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Payment Service")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Get database configuration from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3307")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "payments_db")
	grpcPort := getEnv("GRPC_PORT", "50051")

	// Initialize database connection
	db, err := database.NewMySQL(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Database connection established")

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepositoryMySQL(db.GetDB())

	// Initialize use cases
	processPaymentUC := usecase.NewProcessPaymentUseCase(paymentRepo)
	getPaymentUC := usecase.NewGetPaymentUseCase(paymentRepo)
	cancelPaymentUC := usecase.NewCancelPaymentUseCase(paymentRepo)
	listPaymentsUC := usecase.NewListPaymentsUseCase(paymentRepo)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register payment service
	paymentServiceServer := grpcHandler.NewPaymentServiceServer(
		processPaymentUC,
		getPaymentUC,
		cancelPaymentUC,
		listPaymentsUC,
	)
	pb.RegisterPaymentServiceServer(grpcServer, paymentServiceServer)

	// Register reflection service (useful for grpcurl and debugging)
	reflection.Register(grpcServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		slog.Error("Failed to listen", "port", grpcPort, "error", err)
		os.Exit(1)
	}

	slog.Info("gRPC server listening", "port", grpcPort)

	// Handle graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Failed to serve gRPC", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
