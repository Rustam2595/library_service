package main

import (
	"context"
	"github.com/Rustam2595/library_service/internal/config"
	serv "github.com/Rustam2595/library_service/internal/server"
	store "github.com/Rustam2595/library_service/internal/storage"
	"log"
)

func main() {
	cnf := config.ReadConfig()
	log.Println(cnf)
	var str serv.Storage
	str, err := store.NewRepo(context.Background(), cnf.DBDsn)
	if err != nil {
		str = store.New()
		log.Fatal(err.Error())
	}
	if err := store.Migrations(cnf.DBDsn, cnf.MigratePath); err != nil {
		log.Fatal(err.Error())
	}

	server := serv.New(cnf.Host, str)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
