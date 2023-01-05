package database

import (
	"errors"
	"scraping_service/internal/database/repo"
    "github.com/rs/zerolog"
)

var ErrNotDatabaser = errors.New("Struct is not a Databaser")

type Database struct {
    name string
	dber Databaser
    logger *zerolog.Logger
}

func CreateDatabase(dbName string, logger *zerolog.Logger) (*Database, error) {
	switch dbName {
	case "MongoDB":
		{
			db := &repo.MongoDB{}
			if IsDatabaser(db) {
                return &Database{name: "MongoDB", dber: repo.CreateMongoDB(), logger: logger}, nil
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
