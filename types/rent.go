package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rent struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID  primitive.ObjectID `bson:"userID" json:"userID"`
	MovieID primitive.ObjectID `bson:"movieID" json:"movieID"`
	From    time.Time          `bson:"from" json:"from"`
	To      time.Time          `bson:"to" json:"to"`
}

type CheckRentParams struct {
	UserID  primitive.ObjectID `bson:"userID" json:"userID"`
	MovieID primitive.ObjectID `bson:"movieID" json:"movieID"`
	From    time.Time          `bson:"from" json:"from"`
	To      time.Time          `bson:"to" json:"to"`
}

type CreateRentParams struct {
	UserID  primitive.ObjectID `bson:"userID" json:"userID"`
	MovieID primitive.ObjectID `json:"movieID"`
}

func NewRentFromParams(params CreateRentParams) *Rent {
	return &Rent{
		UserID:  params.UserID,
		MovieID: params.MovieID,
		From:    time.Now(),
		To:      time.Now().Add(time.Hour * 24),
	}
}

/*
func (params CreateRentParams) Validate() map[string]string {
	errors := map[string]string{}

	if params.From.Before(time.Now()) {
		errors["from"] = fmt.Sprintf("rent can't be made in the past, %s", params.From)
	}
	return errors
}
*/
