package main

import (
	"UserRepository/docs"
	"UserRepository/internal/config"
	"UserRepository/internal/handler"
	"UserRepository/internal/repository/postgres"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title User Repository API
// @version 1.0
// @description API сервер для работы с пользователями
// @host localhost:8080
// @BasePath /api
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	docs.SwaggerInfo.Title = "User Repository API"
	docs.SwaggerInfo.Description = "API сервер для работы с пользователями"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.GetByID)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
		r.Get("/", userHandler.List)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
