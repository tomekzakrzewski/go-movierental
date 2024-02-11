package db

import (
	"context"

	"github.com/tomekzakrzewski/go-movierental/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "mongodb"
const userColl = "users"

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewUserStore(client *mongo.Client) *MongoUserStore {
	dbname := dbName
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	res, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	err = res.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
