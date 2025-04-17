package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Establecer conexión con CockroachDB
func Connect(connectionString string) (*sql.DB, error) {
	baseConnectionString := connectionString
	if dbNameIndex := strings.LastIndex(connectionString, "/"); dbNameIndex != -1 {
		if questionMarkIndex := strings.Index(connectionString[dbNameIndex:], "?"); questionMarkIndex != -1 {
			baseConnectionString = connectionString[:dbNameIndex+1] + connectionString[dbNameIndex+questionMarkIndex:]
		} else {
			baseConnectionString = connectionString[:dbNameIndex+1]
		}
	}

	tempDB, err := sql.Open("postgres", baseConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	defer tempDB.Close()

	// Creamos la base de datos si no existe
	_, err = tempDB.Exec("CREATE DATABASE IF NOT EXISTS stockdb")
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	// Ahora nos conectamos a la base de datos
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to CockroachDB")
	return db, nil
}
