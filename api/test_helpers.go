package api

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tomekzakrzewski/go-movierental/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDb) teardown(t *testing.T) {
	dbname := "mongodb"
	//dbname := os.Getenv(db.MongoDBName)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URL")))
	if err != nil {
		t.Fatal(err)
	}
	return &testDb{
		client: client,
		Store: &db.Store{
			User:  db.NewUserStore(client),
			Movie: db.NewMovieStore(client),
			Rent:  db.NewRentStore(client),
		},
	}
}
