package database

import (
	"errors"
	"scraping_service/internal/database/repo"
)

var ErrNotDatabaser = errors.New("Struct is not a Databaser")

type Database struct {
	dber Databaser
}

func CreateDatabase(dbName string) (*Database, error) {
	switch dbName {
	case "MongoDB":
		{
			db := &repo.MongoDB{}
			if IsDatabaser(db) {
				return &Database{dber: repo.CreateMongoDB()}, nil
			} else {
				return nil, ErrNotDatabaser
			}
		}
	default:
		return nil, ErrNotDatabaser
	}

}

func IsDatabaser(db interface{}) bool {
	if _, ok := db.(Databaser); ok {
		return true
	} else {
		return false
	}
}
