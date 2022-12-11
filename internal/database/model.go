package database

import (
	"scraping_service/pkg/models"
	"scraping_service/internal/database/repo"
    "fmt"
    "net/url"
)

type Databaser interface {
	Ping() error
	GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error)
	GetFileLogs(botName string, unixTime int64) (*models.FileLog, error)
    UpdateBot(botName string) error
    GetBot(botName string, qp url.Values) (*models.Bot, error)
    GetFiles(qp url.Values) ([]models.File, error)
    GetLogs(qp url.Values) ([]models.FileLog, error)
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

// Updade bot with name 'botName' with new information - we update logs_count
// and last_run fields, this information can be computed and does't need any
// additionl parameters.
func (db Database) UpdateBot(botName string) error {
	err := db.dber.UpdateBot(botName)
	if err != nil {
		return err
	}

	return nil
}

// Get information from bot named 'botName' or return error.
func (db Database) GetBot(botName string, qp url.Values) (*models.Bot, error) {
	bot, err := db.dber.GetBot(botName, qp)
	if err != nil {
		return nil, db.handleDBErrors(err)
	}

	return bot, nil
}

// Handle errors based on which database we used.
//
// Why?: Each database has its own errors and we just want to give this errors
// some higher abstraction so they can be used in services (eg. if we can't
// connect to the DB maybe try again before returning an error).
//
// Params:
//  - err: Database error to handle
func (db Database) handleDBErrors(err error) error {
    switch db.name {
    case "MongoDB":
        switch err {
        case repo.ErrMongoCursor:
            return ErrCursorFailure
        case repo.ErrMongoDecoding:
            return ErrCantDecode
        case repo.ErrQSUnmarshal:
            return ErrCantUnmarshal
        case repo.ErrObjectIDConversion:
            return ErrCantConvert
        default:
            return ErrFailure
        }
    default:
        fmt.Println("Database ", db.name, " is currently not supported.")
        return ErrDBNotSupported
    }
}

// Try to retrive files from the 'Databaser'. 
//
// Params:
//  - qp: query parameters parsed from URL for building database query.
//
// Errors:
//  - ErrCursorFailure: Database cursor error
//  - ErrCantDecode: Can't decode results into given type (struct,...)
//  - ErrFailure: Other database error
//  - ErrDBNotSupported: Database is not supported
func (db Database) GetFiles(qp url.Values) ([]models.File, error) {
    files, err := db.dber.GetFiles(qp)
    if err != nil {
        return nil, db.handleDBErrors(err)
    }

    return files, nil
}

// Try to retrive files from the 'Databaser'. 
//
// Params:
//  - qp: query parameters parsed from URL for building database query.
//
// Errors:
//  - ErrCursorFailure: Database cursor error
//  - ErrCantDecode: Can't decode results into given type (struct,...)
//  - ErrFailure: Other database error
//  - ErrDBNotSupported: Database is not supported
func (db Database) GetLogs(qp url.Values) ([]models.FileLog, error) {
    logs, err := db.dber.GetLogs(qp)
    if err != nil {
        return nil, db.handleDBErrors(err)
    }

    return logs, nil
}
