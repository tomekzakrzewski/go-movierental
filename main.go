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

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	// init
	var (
		store = &db.Store{
			User:  db.NewUserStore(client),
			Movie: db.NewMovieStore(client),
			Rent:  db.NewRentStore(client),
		}
		rentStore    = db.NewRentStore(client)
		userStore    = db.NewUserStore(client)
		movieHandler = api.NewMovieHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		rentHandler  = api.NewRentHandler(rentStore)
		authHandler  = api.NewAuthHandler(userStore)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin        = apiv1.Group("/admin", api.AdminAuth)
	)

	// USER ONLY RATE, RENT, GET MOVIES, POST AUTH
	// auth handler
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// movie handlers
	apiv1.Get("/movies/:id", movieHandler.HandleGetMovieByID)
	apiv1.Put("/movies/:id/rate", movieHandler.HandleUpdateMovieRating)
	apiv1.Post("/movies/:id/rent", movieHandler.HandleRentMovie)
	apiv1.Post("/movies/rented", movieHandler.HandleGetRentedMovies)
	apiv1.Get("/movies", movieHandler.HandleGetMovies)
	admin.Post("/movies", movieHandler.HandlePostMovie)
	admin.Put("/movies/:id", movieHandler.HandleUpdateMovie)
	admin.Delete("/movies/:id", movieHandler.HandleDeleteMovie)

	// user handlers
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	admin.Post("/users", userHandler.HandlePostUser)
	admin.Get("/users", userHandler.HandleGetUsers)
	admin.Delete("/users/:id", userHandler.HandleDeleteUser)

	//rent handlers
	admin.Get("/rents", rentHandler.HandleGetRents)

	app.Listen(os.Getenv("LISTEN_ADDR"))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
