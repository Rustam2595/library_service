package main

import (
	"context"
	"github.com/Rustam2595/library_service/internal/config"
	"github.com/Rustam2595/library_service/internal/logger"
	serv "github.com/Rustam2595/library_service/internal/server"
	store "github.com/Rustam2595/library_service/internal/storage"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cnf := config.ReadConfig()
	log := logger.Get(cnf.Debug) //add zerolog
	log.Debug().Msg("logger was inited")
	log.Debug().Any("config", cnf).Send()
	//<-- graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		cancel()
	}()
	//-->
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

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		log.Info().Msg("starting server")
		if err := server.Run(gCtx); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		return <-server.ErrChan
	})
	group.Go(func() error {
		<-gCtx.Done()
		return server.ShutdownServer(gCtx)
	})
	if err := group.Wait(); err != nil {
		log.Fatal().Err(err).Msg("fatal server stopped")
	}
	log.Info().Msg("server was stopped")
}
