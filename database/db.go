package database

import (
	"database/sql"

	"github.com/Kleinish/dascr-board/config"
)

// ModifySetupDB to include player stats table
func SetupDB(cfg *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Filename)
	if err != nil {
		return nil, err
	}

	// Table Generation statements
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS player (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            nickname TEXT,
            image TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS player_stats (
            player_id INTEGER PRIMARY KEY,
            games_played INTEGER DEFAULT 0,
            games_won INTEGER DEFAULT 0,
            games_lost INTEGER DEFAULT 0,
            highest_score INTEGER DEFAULT 0,
            best_finish INTEGER DEFAULT 0,
            checkout_rate REAL DEFAULT 0.0,
            three_dart_avg REAL DEFAULT 0.0,
            first_nine_ppr REAL DEFAULT 0.0,
            one_eighties INTEGER DEFAULT 0,
            one_forties INTEGER DEFAULT 0,
            ton_plus INTEGER DEFAULT 0,
            FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE
        );`,
	}

	// Create tables
	for _, create := range createTables {
		_, err = db.Exec(create)
		if err != nil {
			return nil, err
		}
	}

	// Return DB
	return db, nil
}
