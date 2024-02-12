package db

import (
	"context"
	"fmt"
	"time"

	"github.com/tomekzakrzewski/go-movierental/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	rentColl = "rent"
)

type RentStore interface {
	InsertRent(context.Context, *types.Rent) (*types.Rent, error)
	GetRents(context.Context) ([]*types.Rent, error)
	CheckRent(context.Context, types.CheckRentParams) error
}

type MongoRentStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewRentStore(client *mongo.Client) *MongoRentStore {
	return &MongoRentStore{
		client: client,
		coll:   client.Database(MongoDBName).Collection(rentColl),
	}
}

func (s *MongoRentStore) InsertRent(ctx context.Context, rent *types.Rent) (*types.Rent, error) {
	res, err := s.coll.InsertOne(ctx, rent)
	if err != nil {
		return nil, err
	}

	rent.ID = res.InsertedID.(primitive.ObjectID)

	return rent, err
}
func (s *MongoRentStore) GetRents(ctx context.Context) ([]*types.Rent, error) {
	res, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var rents []*types.Rent
	err = res.All(ctx, &rents)
	if err != nil {
		return nil, err
	}
	return rents, nil
}

func (s *MongoRentStore) CheckRent(ctx context.Context, params types.CheckRentParams) error {
	filter := bson.D{
		{"movieID", params.MovieID},
		{"userID", params.UserID},
		{"$or", bson.A{
			bson.D{
				{Key: "to", Value: bson.D{
					{Key: "$gt", Value: params.From},
					{Key: "$lt", Value: params.To},
				}},
			},
			bson.D{
				{Key: "to", Value: bson.D{
					{Key: "$gt", Value: params.From.Add(-24 * time.Hour)},
					{Key: "$lt", Value: params.To},
				}},
			},
		}},
	}

	res, err := s.coll.CountDocuments(ctx, filter)
	fmt.Println(res)
	if err != nil {
		return err
	}
	if res > 0 {
		return fmt.Errorf("already rented")
	}
	return nil
}
