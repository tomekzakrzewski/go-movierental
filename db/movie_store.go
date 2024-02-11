package db

import (
	"context"

	"github.com/tomekzakrzewski/go-movierental/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const MongoMovieName = "mongodb"
const movieColl = "movies"

type MovieStore interface {
	InsertMovie(context.Context, *types.Movie) (*types.Movie, error)
	GetMovies(context.Context) ([]*types.Movie, error)
	GetMovieByID(context.Context, string) (*types.Movie, error)
	PutMovie(context.Context, string, types.UpdateMovieParams) error
	DeleteMovie(context.Context, string) error
	UpdateRating(context.Context, string, int) error
}

type MongoMovieStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMovieStore(client *mongo.Client) *MongoMovieStore {
	dbname := MongoMovieName
	return &MongoMovieStore{
		client: client,
		coll:   client.Database(dbname).Collection(movieColl),
	}
}

func (s *MongoMovieStore) InsertMovie(ctx context.Context, movie *types.Movie) (*types.Movie, error) {
	res, err := s.coll.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}
	movie.ID = res.InsertedID.(primitive.ObjectID)
	return movie, nil
}

func (s *MongoMovieStore) GetMovies(ctx context.Context) ([]*types.Movie, error) {
	res, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var movies []*types.Movie
	err = res.All(ctx, &movies)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MongoMovieStore) PutMovie(ctx context.Context, id string, params types.UpdateMovieParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoMovieStore) DeleteMovie(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	_, err = s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoMovieStore) GetMovieByID(ctx context.Context, id string) (*types.Movie, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var movie types.Movie
	err = s.coll.FindOne(ctx, filter).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *MongoMovieStore) UpdateRating(ctx context.Context, id string, rating int) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"rating": rating}}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
