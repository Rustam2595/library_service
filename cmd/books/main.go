package main

import (
	"context"
	"github.com/Rustam2595/library_service/internal/config"
	"github.com/Rustam2595/library_service/internal/logger"
	serv "github.com/Rustam2595/library_service/internal/server"
	store "github.com/Rustam2595/library_service/internal/storage"
)

func main() {
	cnf := config.ReadConfig()
	//log.Println(cnf)
	log := logger.Get(cnf.Debug) //add zerolog
	log.Debug().Msg("logger was inited")
	log.Debug().Any("config", cnf).Send()
	var str serv.Storage
	str, err := store.NewRepo(context.Background(), cnf.DBDsn)
	if err != nil {
		str = store.New()
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	if err := store.Migrations(cnf.DBDsn, cnf.MigratePath); err != nil {
		log.Fatal().Err(err).Msg("Migrations failed")
	}

	server := serv.New(cnf.Host, str)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
