package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/server"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	srv := server.SetupRoutes()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "local")
	dbPassword := getEnv("DB_PASSWORD", "local")
	dbName := getEnv("DB_NAME", "storex")
	sslMode := getEnv("DB_SSLMODE", string(database.SSLModeDisable))
	serverPort := getEnv("SERVER_PORT", "8080")

	err := database.ConnectAndMigrate(
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		database.SSLMode(sslMode),
	)
	if err != nil {
		fmt.Printf("failed while initialize and migrate database:%v", err)
	}
	fmt.Println("server is running")
	ServerErr := http.ListenAndServe(":8080", srv)
	if ServerErr != nil {
		log.Fatal("")
		return
	}
	fmt.Println("server started at :8080", serverPort)
}
