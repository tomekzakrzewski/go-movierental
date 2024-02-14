package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minTitleLen  = 2
	maxTitleLen  = 200
	minLengthLen = 1
	minGenreLen  = 1
	minRating    = 0
	maxRating    = 10
	minYear      = 1888
)

type Movie struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title  string             `bson:"title" json:"title"`
	Genre  []string           `bson:"genre" json:"genre"`
	Length int                `bson:"length" json:"length"`
	Year   int                `bson:"year" json:"year"`
	Rating int                `bson:"rating" json:"rating"`
}

type CreateMovieParams struct {
	Title  string   `json:"title"`
	Genre  []string `json:"genre"`
	Length int      `json:"length"`
	Year   int      `json:"year"`
}

func NewMovieFromParams(params CreateMovieParams) *Movie {
	return &Movie{
		Title:  params.Title,
		Genre:  params.Genre,
		Length: params.Length,
		Year:   params.Year,
	}
}

type UpdateMovieRating struct {
	Rating int `json:"rating"`
}

type UpdateMovieParams struct {
	Title  string   `json:"title"`
	Genre  []string `json:"genre"`
	Length int      `json:"length"`
	Year   int      `json:"year"`
	Rating int      `json:"rating"`
}

func Validate(params CreateMovieParams) map[string]string {
	errors := map[string]string{}
	if params.Length < minLengthLen {
		errors["lenght"] = fmt.Sprintf("lenght should be at least %d minutes", minLengthLen)
	}
	if params.Year > time.Now().Year()+1 || params.Year < minYear {
		errors["year"] = fmt.Sprintf("year should be between 1888 and %d", time.Now().Year()+1)
	}
	if len(params.Title) < minTitleLen || len(params.Title) > maxTitleLen {
		errors["title"] = fmt.Sprintf("title should be at least %d and max %d characters", minTitleLen, maxTitleLen)
	}
	if len(params.Genre) < minGenreLen {
		errors["genre"] = fmt.Sprintf("movie should have at least %d genre", minGenreLen)
	}
	return errors
}

func (p UpdateMovieParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.Genre) >= minGenreLen {
		m["genre"] = p.Genre
	}
	if len(p.Title) > minTitleLen && len(p.Title) <= maxTitleLen {
		m["title"] = p.Title
	}
	if p.Length > minLengthLen {
		m["length"] = p.Length
	}
	if p.Year >= minYear && p.Year <= time.Now().Year()+1 {
		m["year"] = p.Year
	}
	if p.Rating > minRating && p.Rating <= maxRating {
		m["rating"] = p.Rating
	}
	return m
}
