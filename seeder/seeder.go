package seeder

import (
	"fmt"

	"github.com/theHamdiz/gost/router"
	"github.com/theHamdiz/gost/runner"
)

func DbInit(appName string) error {
	script := GetSeedingScript()

	doneMessage := `Database setup complete. Framework database has been created and populated with initial data.
	`

	db, err := router.GetDbPath(appName)

	if err != nil {
		return err
	}

	err = runner.RunCommand("sqlite3", db, script)
	if err != nil {
		fmt.Println(err)
	}
	err = runner.RunCommand("echo", doneMessage)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func SeedDBData(appName string) error {
	return DbInit(appName)
}

func GetSeedingScript() string {
	return `
-- Create settings table
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY,
    key TEXT UNIQUE,
    value TEXT
);

-- Insert default settings
INSERT OR IGNORE INTO settings (key, value) VALUES ('gost_site_name', 'Gost Site');
INSERT OR IGNORE INTO settings (key, value) VALUES ('admin_email', 'admin@go.dev');

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT,
    email TEXT,
    password TEXT
);

-- Insert default user
INSERT OR IGNORE INTO users (name, email, password) VALUES ('admin', 'admin@go.dev', 'gost');

-- Create plugins table
CREATE TABLE IF NOT EXISTS plugins (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE,
    enabled BOOLEAN
);

-- Insert default plugins
INSERT OR IGNORE INTO plugins (name, enabled) VALUES ('auth', 1);
INSERT OR IGNORE INTO plugins (name, enabled) VALUES ('logger', 1);

-- Create models table
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE,
    schema TEXT
);

-- Insert example model
INSERT OR IGNORE INTO models (name, schema) VALUES ('User', '{\"id\": \"INTEGER PRIMARY KEY\", \"name\": \"TEXT\", \"email\": \"TEXT\"}');
`
}
