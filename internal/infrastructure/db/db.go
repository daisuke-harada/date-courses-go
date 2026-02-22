package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Connect connects to Postgres using environment variables and returns *sql.DB.
// Environment variables (with sensible defaults):
//   - POSTGRES_HOST (default "localhost")
//   - POSTGRES_PORT (default "5432")
//   - POSTGRES_USER (default "postgres")
//   - POSTGRES_PASSWORD (default "")
//   - POSTGRES_DB (default "postgres")
//   - POSTGRES_SSLMODE (default "disable")
//   - DB_MAX_OPEN_CONNS, DB_MAX_IDLE_CONNS, DB_CONN_MAX_LIFETIME (optional)
func Connect(ctx context.Context) (*sql.DB, error) {
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5432")
	user := getenv("POSTGRES_USER", "postgres")
	pass := getenv("POSTGRES_PASSWORD", "")
	dbname := getenv("POSTGRES_DB", "postgres")
	sslmode := getenv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// pool settings
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			dbConn.SetMaxOpenConns(n)
		}
	} else {
		dbConn.SetMaxOpenConns(25)
	}
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			dbConn.SetMaxIdleConns(n)
		}
	} else {
		dbConn.SetMaxIdleConns(25)
	}
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			dbConn.SetConnMaxLifetime(d)
		}
	} else {
		dbConn.SetConnMaxLifetime(5 * time.Minute)
	}

	// quick ping with timeout to verify connectivity
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
		dbConn.Close()
		return nil, err
	}

	return dbConn, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
