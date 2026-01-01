package config

import (
	"cmp"
	"flag"
	"log"
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
	flag.StringVar(&host, "host", defaultHost, "server host")
	flag.StringVar(&dbDsn, "db", defaultDbDSN, "data base address")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")
	debug := flag.Bool("debug", false, "enable debug logging level")
	flag.Parse()

	hostEnv := os.Getenv("SERVER_HOS")
	dbDsnEnv := os.Getenv("DB_DSN")
	migratePathEnv := os.Getenv("MIGRATE_PATH")
	log.Println(hostEnv)
	host = cmp.Or(hostEnv, defaultHost)
	dbDsn = cmp.Or(dbDsnEnv, defaultDbDSN)
	migratePath = cmp.Or(migratePathEnv, defaultMigratePath)
	authAddr := cmp.Or(os.Getenv("AUTH_ADDR"), defaultAuthAddr)
	booksAddr := cmp.Or(os.Getenv("BOOKS_ADDR"), defaultBooksAddr)
	return Config{
		Host:        host,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
		AuthAddr:    authAddr,
		BooksAddr:   booksAddr,
		Debug:       *debug,
	}
}
