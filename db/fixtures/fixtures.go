package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
)

func AddUser(store *db.Store, fName, lName string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Username:  fmt.Sprintf("%s123", fName),
		FirstName: fName,
		LastName:  lName,
		Password:  fmt.Sprintf("%s123", lName),
		Email:     fmt.Sprintf("%s@%s.com", fName, lName),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddMovie(store *db.Store, title string, genre []string, length int, year int) *types.Movie {
	movie := types.NewMovieFromParams(types.CreateMovieParams{
		Title:  title,
		Genre:  genre,
		Length: length,
		Year:   year,
	})
	insertedMovie, err := store.Movie.InsertMovie(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	return insertedMovie
}
