package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guigateixeira/general-auth/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router   http.Handler
	database *database.Queries
	// rdb    *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}

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
	app.database = database.New(conn)

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

	// err := a.rdb.Ping(ctx).Err()
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to redis: %w", err)
	// }

	// defer func() {
	// 	if err := a.rdb.Close(); err != nil {
	// 		fmt.Println("Failed to close redis connection", err)
	// 	}
	// }()

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
