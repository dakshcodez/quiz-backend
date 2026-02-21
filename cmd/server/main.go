// Package main runs the Quiz Backend API server with Gin, Supabase (PostgreSQL), and JWT auth.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"quiz-backend/internal/config"
	"quiz-backend/internal/database"
	"quiz-backend/internal/handlers"
	"quiz-backend/internal/repository"
	"quiz-backend/internal/routes"
	"quiz-backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	quizService := services.NewQuizService(quizRepo, submissionRepo)

	// Handlers (dependency injection)
	authHandler := handlers.NewAuthHandler(authService)
	teacherHandler := handlers.NewTeacherHandler(quizService, cfg.DBTimeout)
	studentHandler := handlers.NewStudentHandler(quizService, cfg.DBTimeout)

	h := &routes.Handlers{
		Auth:    authHandler,
		Teacher: teacherHandler,
		Student: studentHandler,
	}

	router := routes.NewRouter()
	routes.Setup(router, h, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	log.Println("server stopped")
}
