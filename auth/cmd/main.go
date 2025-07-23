package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	grpcHandler "github.com/mesameen/micro-app/auth/internal/handler/grpc"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const serviceName = "auth"

func main() {
	port := 8094
	log.Printf("Starting auth service on port: %d", port)
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to create ceritficate. Error: %v", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen. Error: %v", err)
	}
	h := grpcHandler.New(func() []byte {
		return []byte("test-secret")
	})
	server := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(server)
	gen.RegisterAuthServiceServer(server, h)
	go func() {
		logger.Infof("%s service is up and running on :%d", serviceName, port)
		if err := server.Serve(lis); err != nil {
			logger.Panicf("Failed to serve gRPC reuests. Error: %v", err)
		}
	}()
	// wait till gets shutdown signal
	<-ctx.Done()
	// do graceful shutdown
	server.GracefulStop()
}
