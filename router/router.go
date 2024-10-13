package router

import (
	"database/sql"
	"github.com/yeboahd24/workout-tracker/config"
	"github.com/yeboahd24/workout-tracker/handler"
	"github.com/yeboahd24/workout-tracker/middleware"
	"github.com/yeboahd24/workout-tracker/repository"
	"net/http"
)

func SetupRouter(db *sql.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db)

	// Create handlers
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	exerciseHandler := handler.NewExerciseHandler(exerciseRepo)
	workoutHandler := handler.NewWorkoutHandler(workoutRepo)

	// Auth routes
	mux.HandleFunc("/signup", authHandler.SignUp)
	mux.HandleFunc("/login", authHandler.Login)

	// Exercise routes
	mux.Handle("/exercises",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(exerciseHandler.GetAll)))
	mux.Handle("/exercises/create",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(exerciseHandler.Create)))

	// Workout routes
	mux.Handle("/workouts",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(workoutHandler.GetByUser)))
	mux.Handle("/workouts/create",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(workoutHandler.Create)))
	mux.Handle("/workouts/update",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(workoutHandler.Update)))
	mux.Handle("/workouts/delete",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			http.HandlerFunc(workoutHandler.Delete)))
	mux.Handle("/workouts/report", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(workoutHandler.GenerateReport)))

	return mux
}
