package config

import (
	"cmp"
	"flag"
	"os"
)

type Config struct {
	Host        string
	DBDsn       string
	MigratePath string
	AuthAddr    string
	BooksAddr   string
	Debug       bool
}

const (
	defaultDbDSN       = "postgres://ru:2595@127.0.0.1:5432/ru_DB?sslmode=disable"
	defaultMigratePath = "migrations"
	defaultHost        = ":8080"
	defaultAuthAddr    = "localhost:8081"
	defaultBooksAddr   = "localhost:8082"
)

func ReadConfig() Config {
	var host, dbDsn, migratePath string
	flag.StringVar(&host, "host", "", "server host")
	flag.StringVar(&dbDsn, "db", "", "data base address")
	flag.StringVar(&migratePath, "m", "", "path to migrations")
	debug := flag.Bool("debug", false, "enable debug logging level")
	flag.Parse()

	hostEnv := os.Getenv("SERVER_HOST")
	dbDsnEnv := os.Getenv("DB_DSN")
	migratePathEnv := os.Getenv("MIGRATE_PATH")
	authAddrEnv := os.Getenv("AUTH_ADDR")
	booksAddrEnv := os.Getenv("BOOKS_ADDR")

	host = cmp.Or(host, hostEnv, defaultHost)
	dbDsn = cmp.Or(dbDsn, dbDsnEnv, defaultDbDSN)
	migratePath = cmp.Or(migratePath, migratePathEnv, defaultMigratePath)
	authAddr := cmp.Or(authAddrEnv, defaultAuthAddr)
	booksAddr := cmp.Or(booksAddrEnv, defaultBooksAddr)

	return Config{
		Host:        host,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
		AuthAddr:    authAddr,
		BooksAddr:   booksAddr,
		Debug:       *debug,
	}
}
