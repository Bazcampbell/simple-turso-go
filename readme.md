# simple-turso

Simple wrapper for [Turso/libSQL](https://turso.tech) with very straightforward CRUD operations.

Focused on dead-simple logging + common CRUD operations.

## Features

- Minimal external dependencies (net/http + standard libs)
- Simple structured logging to a `logs` table (id, timestamp, level, message, application)
- Generic `Select[T]` helper with custom row scanning to read
- Straightforward `Execute` for INSERT/UPDATE/DELETE queries
- Global singleton client (initialise once with url + key)

## Installation

```bash
go get github.com/Bazcampbell/simple-turso-go/pkg

```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    simpleturso "github.com/bazcampbell/simple-turso-go/pkg"
)

// Simple struct representing a DB table
type User struct {
    ID    int64
    Name  string
    Email string
}

func main() {
    godotenv.Load()

    config := turso.Config{
        DbUrl: os.Getenv("TURSO_URL"),  // libsql or https 
        DbKey: os.Getenv("TURSO_KEY"),
    }

    if err := turso.Init(config); err != nil {
        // handle
    }

    // ------------------------------
    // Logging
    // ------------------------------
    err := turso.LogToTurso(
        turso.LogLevelWarning,
        "application started",
        "my-app",
        "Australia/Sydney",
    )
    if err != nil {
        // handle
    }

    // ------------------------------
    // Query with generic Select
    // ------------------------------

    // Define the scan function first
    scanUser := func(row map[string]any) (User, error) {
        // Type checks optional
        id, ok := row["id"].(int64)
        if !ok {
            return User{}, fmt.Errorf("id is not int64")
        }

        name, ok := row["name"].(string)
        if !ok {
            return User{}, fmt.Errorf("name is not string")
        }

        email, ok := row["email"].(string)
        if !ok {
            return User{}, fmt.Errorf("email is not string")
        }
        return User{ID: id, Name: name, Email: email}, nil
    }
    
    users, err := turso.Select[User](
        "SELECT id, name, email FROM users WHERE email LIKE ?",
        []interface{"%@example.com"},
        scanUser,
    )
    if err != nil {
        // handle
    }


    // ------------------------------
    // Execute (INSERT/UPDATE/DELETE)
    // ------------------------------
    err = turso.Execute(
        "UPDATE users SET name = ? WHERE id = ?",
        []interface{"Updated Name", 42},
    )

    if err != nil {
        // handle
    }
}
