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
		rentStore    = db.NewRentStore(client)
		movieStore   = db.NewMovieStore(client)
		movieHandler = api.NewMovieHandler(movieStore, rentStore)
		userStore    = db.NewUserStore(client)
		userHandler  = api.NewUserHandler(userStore)
		rentHandler  = api.NewRentHandler(rentStore)
		authHandler  = api.NewAuthHandler(userStore)
		app          = fiber.New()
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", api.JWTAuthentication(userStore))
	)

	// USER ONLY RATE, RENT, GET MOVIES, POST AUTH
	// auth handler
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// movie handlers
	apiv1.Post("/movies", movieHandler.HandlePostMovie)
	apiv1.Get("/movies", movieHandler.HandleGetMovie)
	apiv1.Put("/movies/:id", movieHandler.HandleUpdateMovie)
	apiv1.Delete("/movies/:id", movieHandler.HandleDeleteMovie)
	apiv1.Get("/movies/:id", movieHandler.HandleGetMovieByID)
	apiv1.Put("/movies/:id/rate", movieHandler.HandleUpdateMovieRating)
	// user handlers
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	//rent handlers
	apiv1.Post("/rents", rentHandler.HandlePostRent)
	apiv1.Get("/rents", rentHandler.HandleGetRents)
	apiv1.Post("/rents/:id/movie", movieHandler.HandleRentMovie)

	app.Listen(os.Getenv("LISTEN_ADDR"))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
