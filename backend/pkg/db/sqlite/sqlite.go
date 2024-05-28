package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// Init initializes the database migrations
func Init() {
	// Set up SQLite database connection
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	// Set up migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://backend/pkg/db/migrations/sqlite", // Path to your migration files
		"sqlite", driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully!")
}

// LoadMigrations loads migrations from the specified path and returns sorted filenames
func LoadMigrations(path string) ([]string, error) {
	var migrations []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".up.sql") {
			return nil
		}
		migrations = append(migrations, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(migrations)
	return migrations, nil
}

func InitMigration() {
	// Load and apply migrations
	migrationFiles, err := LoadMigrations("backend/pkg/db/migrations/sqlite")
	if err != nil {
		fmt.Printf("Failed to load migrations: %v\n", err)
		return
	}
	for _, file := range migrationFiles {
		fmt.Printf("Applying migration from file: %s\n", file)
	}
	fmt.Println()
}
