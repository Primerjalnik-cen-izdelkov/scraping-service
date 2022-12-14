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

type AuthDatabase struct {
    name string
    dber AuthDatabaser
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

func CreateAuthDatabase(dbName string) (*AuthDatabase, error) {
    switch dbName {
    case "AuthPostgresDB":
        {
            db := &repo.AuthPostgresDB{}
            if IsAuthDatabaser(db) {
                return &AuthDatabase{name: "AuthPostgresDB", dber: repo.CreateAuthPostgresDB()}, nil
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

func IsAuthDatabaser(db interface{}) bool {
	if _, ok := db.(AuthDatabaser); ok {
		return true
	} else {
		return false
	}
}
