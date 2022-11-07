package database

import (
	"scraping_service/pkg/models"
)

type Databaser interface {
	Ping() error
	GetBotStats(botName, querry string) ([]models.FileStat, error)
	GetFileStats(botName string, unixTime int64) (*models.FileStat, error)
	GetBotFileNames(botName string) ([]models.File, error)
}

// TODO(miha): Should we call functions from here? This way we only need to
// import database (func CreateDatabase(name string) Databaser { ), and can add
// some loggers/metrics in easy way along the way.

func (db Database) Ping() error {
	return db.dber.Ping()
}

func (db Database) CreateDatabase(dbName string) error {
	return db.dber.Ping()
}
func (db Database) CreateTable(tableName string) error {
	return db.dber.Ping()
}

func (db Database) GetBotStats(botName, querry string) ([]models.FileStat, error) {
	files, err := db.dber.GetBotStats(botName, querry)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (db Database) GetFileStats(botName string, unixTime int64) (*models.FileStat, error) {
	file, err := db.dber.GetFileStats(botName, unixTime)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (db Database) GetBotFileNames(botName string) ([]models.File, error) {
	files, err := db.dber.GetBotFileNames(botName)
	if err != nil {
		return nil, err
	}

	return files, nil
}
