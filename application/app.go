package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guigateixeira/general-auth/handler"
	"github.com/guigateixeira/general-auth/internal/database"
	"github.com/guigateixeira/general-auth/repositories"
	"github.com/guigateixeira/general-auth/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router   http.Handler
	database *database.Queries
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

	// Initialize repositories
	userRepo := repositories.New(databaseConn)

	// Initialize services
	userSvc := services.New(userRepo)

	// Initialize handlers
	userHandler := handler.New(userSvc)

	app := &App{
		router:   loadRoutes(userHandler),
		database: databaseConn,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DB URL is not found")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	fmt.Println("Starting server in port", port)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
