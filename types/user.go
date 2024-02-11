package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minUsernameLen  = 3
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
	emailRegex      = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username          string             `bson:"username" json:"username"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	Email             string             `bson:"email" json:"email"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

type CreateUserParams struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Username:          params.Username,
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		EncryptedPassword: string(encpwd),
		Email:             params.Email,
	}, nil
}

func (p CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.Username) < minUsernameLen {
		errors["username"] = fmt.Sprintf("username must be at least %d characters long", minUsernameLen)
	}
	if len(p.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFirstNameLen)
	}
	if len(p.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLastNameLen)
	}
	if len(p.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLen)
	}
	if isValidEmail(p.Email) == false {
		errors["email"] = fmt.Sprintf("invalid email address: %s", p.Email)
	}
	return errors
}

func isValidEmail(e string) bool {
	email := regexp.MustCompile(emailRegex)
	return email.MatchString(e)
}
