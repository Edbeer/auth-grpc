package rest

import (
	"context"
	"log"
	"net"

	"net/http"
	"time"

	"github.com/Edbeer/microservices/internal/config"
	accountpb "github.com/Edbeer/microservices/proto/api/account/v1"
	examplepb "github.com/Edbeer/microservices/proto/api/example/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type Deps struct {
	Account accountpb.AccountServiceServer
	Example examplepb.ExampleServiceServer
	Config  *config.Config
}

type RestServer struct {
	srv  *http.Server
	Deps Deps
}

func NewRestServer(deps Deps) *RestServer {
	return &RestServer{
		srv: &http.Server{
			Addr:           deps.Config.RestServer.Port,
			ReadTimeout:    time.Duration(deps.Config.RestServer.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(deps.Config.RestServer.WriteTimeout) * time.Second,
			IdleTimeout:    time.Duration(deps.Config.RestServer.IdleTimeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		Deps: deps,
	}
}

func (s *RestServer) ListenAndServe(ctx context.Context) error {

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	
	// account handler
	err := accountpb.RegisterAccountServiceHandlerServer(ctx, grpcMux, s.Deps.Account)
	if err != nil {
		log.Fatal(err)
	}
	// example handler
	err = examplepb.RegisterExampleServiceHandlerServer(ctx, grpcMux, s.Deps.Example)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	s.srv.Handler = mux

	lis, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		log.Fatal(err)
	}


	if s.Deps.Config.RestServer.TLS {
		if err := s.srv.ServeTLS(lis, "cert/server-cert.pem", "cert/server-key.pem"); err != nil {
			log.Fatal(err)
		}
	}

	if err := s.srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *RestServer) Stop(ctx context.Context) {
	s.srv.Shutdown(ctx)
}
