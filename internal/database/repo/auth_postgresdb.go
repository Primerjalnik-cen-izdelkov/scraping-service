package repo

import (
    "fmt"
    "github.com/eaigner/jet"
    "github.com/lib/pq"
    "scraping_service/pkg/models"
    "os"
)

type AuthPostgresDB struct {
    Client *jet.Db
}

func CreateAuthPostgresDB() *AuthPostgresDB {
    uri := os.Getenv("PG_URI")
    pqUrl, err := pq.ParseURL(uri)
    if err != nil {
        fmt.Println("parse url err")
        return nil
    }

    client, err := jet.Open("postgres", pqUrl)
    if err != nil {
        fmt.Println("jet open err")
        return nil
    }

    return &AuthPostgresDB{Client: client}
}

func (db AuthPostgresDB) GetUser(name string) (models.User, error) {
    var users []models.User
    err := db.Client.Query("SELECT * FROM users WHERE users.name = ($1)", name).Rows(&users)
    if err != nil {
        fmt.Println("db.client err", err)
    }
    for _, user := range users {
        fmt.Printf("id: %d, name: %s, hash: %s\n", user.Id, user.Name, user.PasswordHash)
    }

    return users[0], nil
}
