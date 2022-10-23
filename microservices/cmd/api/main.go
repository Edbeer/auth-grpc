package main

import (
	"context"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/Edbeer/microservices/internal/config"
	service "github.com/Edbeer/microservices/internal/services"
	postgres "github.com/Edbeer/microservices/internal/storage/psql"
	redstorage "github.com/Edbeer/microservices/internal/storage/redis"
	"github.com/Edbeer/microservices/internal/transport/grpc"
	"github.com/Edbeer/microservices/internal/transport/grpc/handlers"
	"github.com/Edbeer/microservices/internal/transport/grpc/interceptor"
	"github.com/Edbeer/microservices/internal/transport/rest"
	"github.com/Edbeer/microservices/pkg/db/psql"
	"github.com/Edbeer/microservices/pkg/db/redis"
	"github.com/Edbeer/microservices/pkg/jwt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	// context
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

	// jaeger init
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: config.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.Jaeger.LogSpans,
			LocalAgentHostPort: config.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	log.Println("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	log.Println("Opentracing connected")

	// jwt token manager
	manager, err := jwt.NewManager(config.GrpsServer.JwtSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Init storage, service and handlers
	storagePsql := postgres.NewStorage(psql)
	redisStorage := redstorage.NewStorage(redis)
	service := service.NewService(service.Deps{
		Psql:    storagePsql,
		Redis:   redisStorage,
		Manager: manager,
	})
	handlers := handlers.NewHandlers(handlers.Deps{
		AccountService: service.Account,
		SessionService: service.Session,
		Config:         config,
	})
	// authorization interceptor
	interceptor := interceptor.NewAccountInterceptor(manager)
	interceptor.SetMinimumPermissionLevelForMethod("/example.v1.ExampleService/Hello", "admin")
	interceptor.SetMinimumPermissionLevelForMethod("/example.v1.ExampleService/World")
	interceptor.SetMinimumPermissionLevelForMethod("/example.v1.ExampleService/StreamWorld", "admin")
	// Init grpc server
	grpcServer := grpc.NewServer(grpc.Deps{
		Account:     handlers.Account,
		Example:     handlers.Example,
		Interceptor: interceptor,
		Config:      config,
	})
	
	go func() {
		log.Println("Listening grpc server on port:", config.GrpsServer.Port)
		if err := grpcServer.ListenAndServe(config.GrpsServer.Port); err != nil {
			log.Fatal(err)
		}
	}()
	
	// init rest server
	restServer := rest.NewRestServer(rest.Deps{
		Account: handlers.Account,
		Example: handlers.Example,
		Config: config,
	})
	go func() {
		log.Println("Listening rest server on port:", restServer.Deps.Config.RestServer.Port)
		if err := restServer.ListenAndServe(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case q := <-quit:
		log.Printf("signal.Notify: %v", q)
	case done := <-ctx.Done():
		log.Printf("ctx.Done: %v", done)
	}

	grpcServer.Stop()
	restServer.Stop(ctx)
}
