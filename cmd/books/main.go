package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rustam2595/library_service/internal/config"
	authservicev1 "github.com/Rustam2595/library_service/internal/gen/go"
	books_servicev1 "github.com/Rustam2595/library_service/internal/genBooks/go"
	"github.com/Rustam2595/library_service/internal/logger"
	serv "github.com/Rustam2595/library_service/internal/server"
	store "github.com/Rustam2595/library_service/internal/storage"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if err = store.Migrations(cnf.DBDsn, cnf.MigratePath); err != nil {
		log.Fatal().Err(err).Msg("Migrations failed")
	}
	log.Debug().Msg("go to client")
	//auth_service:
	// Подключаемся к серверу
	connAuth, err := grpc.NewClient(cnf.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to grpc auth server")
	}
	defer func() {
		if err = connAuth.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to stop users gRPC server")
		}
	}()
	// Создаём клиента
	clientAuth := authservicev1.NewAuthServiceClient(connAuth)

	//books_service:
	// Подключаемся к серверу
	connBooks, err := grpc.NewClient(cnf.BooksAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to grpc books server")
	}
	defer func() {
		if err = connBooks.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to stop books gRPC server")
		}
	}()
	// Создаём клиента
	clientBooks := books_servicev1.NewBooksServiceClient(connBooks)

	server := serv.New(cnf.Host, str, clientAuth, clientBooks)

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		log.Info().Msg("starting server")
		if err = server.Run(gCtx); err != nil {
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
	if err = group.Wait(); err != nil {
		log.Fatal().Err(err).Msg("fatal server stopped")
	}
	log.Info().Msg("server was stopped")
}
