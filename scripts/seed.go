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
}
