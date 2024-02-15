package db

const (
	MongoDBName = "mongodb"
)

type Pagination struct {
	Page  int
	Limit int
}

type Store struct {
	User  UserStore
	Movie MovieStore
	Rent  RentStore
}
