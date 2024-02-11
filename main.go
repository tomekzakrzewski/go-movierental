package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tomekzakrzewski/go-movierental/api"
	"github.com/tomekzakrzewski/go-movierental/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	// init
	var (
		movieStore   = db.NewMovieStore(client)
		movieHandler = api.NewMovieHandler(movieStore)
		app          = fiber.New()
	)

	// handlers
	app.Post("/movies", movieHandler.HandlePostMovie)
	app.Get("/movies", movieHandler.HandleGetMovie)
	app.Put("/movies/:id", movieHandler.HandleUpdateMovie)
	app.Delete("/movies/:id", movieHandler.HandleDeleteMovie)

	app.Listen(os.Getenv("LISTEN_ADDR"))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
