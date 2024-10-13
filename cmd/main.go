package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/yeboahd24/workout-tracker/config"
	"github.com/yeboahd24/workout-tracker/router"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database connection
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Initialize router
	r := router.SetupRouter(db, cfg)

	// Start server
	log.Printf("Server starting on port %d", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
