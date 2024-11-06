package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mummumgoodboy/verify"
	"github.com/onfirebyte/todo-dumb/internal/model"
	"github.com/onfirebyte/todo-dumb/internal/route"
	"github.com/onfirebyte/todo-dumb/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		log.Fatal("POSTGRES_URI is not set")
	}

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Todo{})

	publicKey := os.Getenv("JWT_PUBLIC_KEY")
	port := os.Getenv("PORT")

	log.Println("Database migrated")
	userService, err := service.NewTodoService(db)
	if err != nil {
		log.Fatal("Error creating user service", err)
	}

	verifier, err := verify.NewJWTVerifier(publicKey)
	if err != nil {
		log.Fatal("Error creating verifier", err)
	}

	route.CreateTodoRoute(userService, verifier)

	// start the server
	err = http.ListenAndServe(":"+port, nil)

	log.Println("Auth server started on port", port)
	if err != nil {
		log.Fatal(err)
	}
}
