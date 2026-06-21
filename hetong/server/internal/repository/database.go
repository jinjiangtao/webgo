package repository

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dataDir, "dashboard.db")
	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}

	if err = createTables(DB); err != nil {
		return err
	}

	tables := []string{"dim_time", "dim_region", "dim_business", "agg_data", "raw_data", "snapshots"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			log.Printf("Warning: failed to check table %s: %v", table, err)
		} else {
			log.Printf("Table %s exists: %v", table, name == table)
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS dim_time (
			date TEXT PRIMARY KEY,
			year INTEGER NOT NULL,
			quarter INTEGER NOT NULL,
			month INTEGER NOT NULL,
			week INTEGER NOT NULL,
			month_name TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS dim_region (
			region_code TEXT PRIMARY KEY,
			region_name TEXT NOT NULL,
			parent_code TEXT,
			level INTEGER NOT NULL,
			path TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS dim_business (
			business_code TEXT PRIMARY KEY,
			business_name TEXT NOT NULL,
			parent_code TEXT,
			level INTEGER NOT NULL,
			category TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS agg_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			region_code TEXT NOT NULL,
			business_code TEXT NOT NULL,
			sales REAL NOT NULL DEFAULT 0,
			orders INTEGER NOT NULL DEFAULT 0,
			users INTEGER NOT NULL DEFAULT 0,
			amount REAL NOT NULL DEFAULT 0,
			agg_level INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS raw_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trace_no TEXT NOT NULL UNIQUE,
			date TEXT NOT NULL,
			region_code TEXT NOT NULL,
			business_code TEXT NOT NULL,
			order_no TEXT NOT NULL,
			amount REAL NOT NULL,
			user_id TEXT NOT NULL,
			product_name TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			trade_time DATETIME NOT NULL,
			extra TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			state TEXT NOT NULL,
			drill_path TEXT NOT NULL,
			filters TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT DEFAULT 'system'
		)`,
	}

	for _, tableSQL := range tables {
		_, err := db.Exec(tableSQL)
		if err != nil {
			return err
		}
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_agg_date ON agg_data(date)`,
		`CREATE INDEX IF NOT EXISTS idx_agg_region ON agg_data(region_code)`,
		`CREATE INDEX IF NOT EXISTS idx_agg_business ON agg_data(business_code)`,
		`CREATE INDEX IF NOT EXISTS idx_agg_level ON agg_data(agg_level)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_date ON raw_data(date)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_region ON raw_data(region_code)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_business ON raw_data(business_code)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_trace ON raw_data(trace_no)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshot_created ON snapshots(created_at)`,
	}

	for _, idxSQL := range indexes {
		_, err := db.Exec(idxSQL)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDB() *sql.DB {
	return DB
}
