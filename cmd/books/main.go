package main

import (
	"github.com/Rustam2595/library_service/internal/config"
	serv "github.com/Rustam2595/library_service/internal/server"
	store "github.com/Rustam2595/library_service/internal/storage"
	"log"
)

func main() {
	cnf := config.ReadConfig()
	log.Println(cnf)
	storage := store.New()
	server := serv.New(cnf.Host, storage)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
