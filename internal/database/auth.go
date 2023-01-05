package database

import (
    "scraping_service/pkg/models"
)

type AuthDatabaser interface {
    GetUser(name string) (models.User, error)
}


func (db AuthDatabase) GetUser(name string) (models.User, error) {
    return db.dber.GetUser(name)
}
