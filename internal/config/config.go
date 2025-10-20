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
	defaultDbDSN       = "postgres://ru:2595@127.0.0.1:5432/ru_DB?sslmode=disable"
	defaultMigratePath = "migrations"
	defaultHost        = ":8080"
	//defaultAuthAddr    = "localhost:8081"
)

func ReadConfig() Config {
	var host, dbDsn, migratePath string
	flag.StringVar(&host, "host", defaultHost, "server host")
	flag.StringVar(&dbDsn, "db", defaultDbDSN, "data base addres")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")
	debug := flag.Bool("debug", false, "enable debug logging level")
	flag.Parse()

	hostEnv := os.Getenv("SERVER_HOS")
	dbDsnEnv := os.Getenv("DB_DSN")
	migratePathEnv := os.Getenv("MIGRATE_PATH")
	log.Println(hostEnv)
	if hostEnv != "" && host == defaultHost {
		host = hostEnv
	}
	if dbDsnEnv != "" && dbDsn == defaultDbDSN {
		dbDsn = dbDsnEnv
	}
	if migratePathEnv != "" && migratePath == defaultMigratePath {
		migratePath = migratePathEnv
	}
	//authAddr := cmp.Or(os.Getenv("AUTH_ADDR"), defaultAuthAddr)
	return Config{
		Host:        host,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
		//AuthAddr:    authAddr,
		Debug: *debug,
	}
}

//func getEnvDefault(key, defaultVal, currVal string) *string {
//	if os.Getenv(key) != "" && currVal == defaultVal {
//		log.Printf("defaultValue of %s is %s", key, defaultVal)
//		return &defaultVal
//	}
//	log.Printf("currentValue of %s is %s", key, currVal)
//	return &currVal
//}
