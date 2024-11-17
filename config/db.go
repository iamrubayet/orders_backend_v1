package config

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func ConnectDB(cfg *Config) (*sql.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    return db, nil
}
