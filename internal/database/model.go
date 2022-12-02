package database

import (
	"scraping_service/pkg/models"
    "fmt"
)

type Databaser interface {
	Ping() error
	GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error)
	GetFileLogs(botName string, unixTime int64) (*models.FileLog, error)
	GetBotFileNames(botName string) ([]models.File, error)
    UpdateBot(botName string) error
    GetBot(botName string) (*models.Bot, error)
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

func (db Database) GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error) {
	logs, err := db.dber.GetBotLogs(botName, qm)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (db Database) GetFileLogs(botName string, unixTime int64) (*models.FileLog, error) {
	file, err := db.dber.GetFileLogs(botName, unixTime)
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

func (db Database) UpdateBot(botName string) error {
	err := db.dber.UpdateBot(botName)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetBot(botName string) (*models.Bot, error) {
	bot, err := db.dber.GetBot(botName)
	if err != nil {
		return nil, err
	}

    fmt.Println("model.go:", bot)

	return bot, nil
}
