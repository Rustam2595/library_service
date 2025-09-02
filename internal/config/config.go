package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Host        string
	DBDsn       string
	MigratePath string
	//AuthAddr    string
	Debug bool
}

const (
	defaultBbDSN       = "postgres://ru:2595@127.0.0.1:5432/ru_DB?sslmode=disable"
	defaultMigratePath = "migrations"
	defaultHost        = ":8080"
	//defaultAuthAddr    = "localhost:8081"
)

func ReadConfig() Config {
	var host, dbDsn, migratePath string
	flag.StringVar(&host, "host", defaultHost, "server host")
	flag.StringVar(&dbDsn, "db", defaultBbDSN, "data base address")
	flag.StringVar(&migratePath, "mPath", defaultMigratePath, "migrate path")
	debug := flag.Bool("debug", false, "enable debug logging lvl")
	flag.Parse()
	hostEnv := os.Getenv("SERVER_HOST") //хост взяли из переменной окружения (echo $SERVER_HOST)
	if hostEnv != "" && host == defaultHost {
		host = hostEnv
	}
	log.Println("host: ", host)
	dbDsnEnv := os.Getenv("DB_DSN") //хост взяли из переменной окружения (echo $DB_DSN)
	if dbDsnEnv != "" && dbDsn == defaultBbDSN {
		dbDsn = dbDsnEnv
	}
	log.Println("dbDsn: ", dbDsn)
	migratePathEnv := os.Getenv("MIGRATE_PATH") //хост взяли из переменной окружения (echo $MIGRATE_PATH)
	if migratePathEnv != "" && migratePath == defaultMigratePath {
		migratePath = migratePathEnv
	}
	log.Println("migratePath: ", migratePath)
	return Config{
		Host:        host,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
		Debug:       *debug,
	}
}
