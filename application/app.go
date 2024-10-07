package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/guigateixeira/general-auth/handler"
	"github.com/guigateixeira/general-auth/internal/database"
	"github.com/guigateixeira/general-auth/kafka"
	"github.com/guigateixeira/general-auth/repositories"
	"github.com/guigateixeira/general-auth/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router      http.Handler
	dbConn      *sql.DB
	database    *database.Queries
	kafkaClient *kafka.KafkaClient
	// rdb    *redis.Client
}

func New() *App {
	godotenv.Load(".env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DB URL is not found")
	}

	// Connect to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Initialize the queries object and set it to the App struct
	databaseConn := database.New(conn)

	// Initialize Kafka client
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaClient, err := kafka.NewKafkaClient(kafkaBrokers)
	if err != nil {
		log.Fatal("failed to create Kafka client:", err)
	}

	// Initialize repositories
	userRepo := repositories.New(databaseConn)

	// Initialize services
	userSvc := services.New(userRepo)

	// Initialize handlers
	userHandler := handler.New(userSvc)

	app := &App{
		router:      loadRoutes(userHandler),
		dbConn:      conn,
		database:    databaseConn,
		kafkaClient: kafkaClient,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	fmt.Println("Starting server on port", port)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Shutting down server...")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the HTTP server
		if err := server.Shutdown(timeout); err != nil {
			fmt.Printf("Error during server shutdown: %v\n", err)
		}

		// Call the App's Shutdown method to clean up other resources
		if err := a.Shutdown(); err != nil {
			fmt.Printf("Error during app shutdown: %v\n", err)
		}

		return nil
	}
}

func (a *App) Shutdown() error {
	// Close the Kafka client
	if a.kafkaClient != nil {
		if err := a.kafkaClient.Close(); err != nil {
			return fmt.Errorf("error closing Kafka client: %w", err)
		}
	}

	// Close the database connection
	if a.database != nil {
		if err := a.dbConn.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
	}
	return nil
}
