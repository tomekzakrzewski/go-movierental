package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomekzakrzewski/go-movierental/api"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		User:  db.NewUserStore(client),
		Movie: db.NewMovieStore(client),
		Rent:  db.NewRentStore(client),
	}

	user := fixtures.AddUser(store, "tomek", "zak", false)
	fmt.Println("tomek ->", api.CreateTokenFromUser(user))
	user = fixtures.AddUser(store, "zuzia", "poz", false)
	fmt.Println("zuzia ->", api.CreateTokenFromUser(user))
	user = fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(user))

	fixtures.AddMovie(store, "The Matrix", []string{"Action", "Sci-Fi"}, 120, 1999)
	fixtures.AddMovie(store, "Titanic", []string{"Drama", "Romance"}, 194, 1999)
	fixtures.AddMovie(store, "Star Wars: The Force Awakens", []string{"Action", "Sci-Fi"}, 136, 2015)
	fixtures.AddMovie(store, "The Godfather", []string{"Drama"}, 175, 1972)
	fixtures.AddMovie(store, "The Shawshank Redemption", []string{"Drama"}, 142, 1994)
	fixtures.AddMovie(store, "Schindler's List", []string{"Biography", "Drama", "History"}, 195, 1993)

}
