package database

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed *.sql
var schemaFS embed.FS

// TODO: Need to update this to store in a better location for distribution
func InitDB() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(homeDir, ".gorunning.db")
	db, err := sql.Open("sqlite", dbPath)

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetLocation(db *sql.DB) (string, error) {
	var address string
	err := db.QueryRow("SELECT address FROM location ORDER BY id DESC LIMIT 1").Scan(&address)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return address, err
}

func SaveLocation(db *sql.DB, address string) error {
	_, err := db.Exec("INSERT INTO location (address) VALUES (?)", address)
	return err
}
