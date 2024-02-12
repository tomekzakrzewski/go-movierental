package db

const (
	MongoDBName = "mongodb"
)

type Store struct {
	User  UserStore
	Movie MovieStore
	Rent  RentStore
}
