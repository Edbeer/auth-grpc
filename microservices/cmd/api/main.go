package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"guthub.com/Edbeer/microservices/internal/config"
	service "guthub.com/Edbeer/microservices/internal/services"
	postgres "guthub.com/Edbeer/microservices/internal/storage/psql"
	redstorage "guthub.com/Edbeer/microservices/internal/storage/redis"
	"guthub.com/Edbeer/microservices/internal/transport/grpc"
	"guthub.com/Edbeer/microservices/internal/transport/grpc/handlers"
	"guthub.com/Edbeer/microservices/internal/transport/grpc/interceptor"
	"guthub.com/Edbeer/microservices/pkg/db/psql"
	"guthub.com/Edbeer/microservices/pkg/db/redis"
	"guthub.com/Edbeer/microservices/pkg/jwt"
)

func accessRole() map[string][]string {
	const examplePath = "/example.v1.ExampleService/"
	return map[string][]string{
		examplePath + "Hello": {"admin", "user"},
		examplePath + "World": {"admin"},
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//config
	config := config.GetConfig()

	// postgres
	psql, err := psql.NewPsqlDB(config)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v\n", psql.Stats())
	}
	defer psql.Close()

	// redis
	redis := redis.NewRedisClient(config)
	defer redis.Close()
	log.Println("Redis connected")

	// jwt token manager
	manager, err := jwt.NewManager(config.GrpsServer.JwtSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Init storage, service and handlers
	storagePsql := postgres.NewStorage(psql)
	redisStorage := redstorage.NewStorage(redis)
	service := service.NewService(service.Deps{
		Psql: storagePsql,
		Redis: redisStorage,
		Manager: manager,
	})
	handlers := handlers.NewHandlers(handlers.Deps{
		AccountService: service.Account,
		SessionService: service.Session,
	})
	interceptor := interceptor.NewAccountInterceptor(manager, accessRole())
	// Init grpc server
	grpcServer := grpc.NewServer(grpc.Deps{
		Account: handlers.Account,
		Example: handlers.Example,
		Interceptor: interceptor,
	})

	go func() {
		if err := grpcServer.ListenAndServe(config.GrpsServer.Port); err != nil {
			log.Fatal(err)
		}
	}()
	defer grpcServer.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case q := <-quit:
		log.Printf("signal.Notify: %v", q)
	case done := <-ctx.Done():
		log.Printf("ctx.Done: %v", done)
	}
}