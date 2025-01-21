package database

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

//go:embed *.sql
var schemaFS embed.FS

func getDataDir() (string, error) {
	var dataDir string

	switch goos := runtime.GOOS; goos {
	case "windows":
		dataDir = os.Getenv("APPDATA")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(home, "Library", "Application Support")
	default: // linux and other unix systems
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(home, ".local", "share")
	}

	appDataDir := filepath.Join(dataDir, "gorunning")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return "", err
	}

	return appDataDir, nil
}

func InitDB() (*sql.DB, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "gorunning.db")
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
