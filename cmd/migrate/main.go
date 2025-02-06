// cmd/migrate/main.go
package main

import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    var direction string
    flag.StringVar(&direction, "direction", "up", "migration direction (up or down)")
    flag.Parse()

    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }

    switch direction {
    case "up":
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    case "down":
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    default:
        log.Fatal("invalid direction")
    }

    fmt.Printf("Migration %s completed successfully\n", direction)
}