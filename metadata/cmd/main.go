package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/mesameen/micro-app/metadata/internal/controller/metadata"
	grpchandler "github.com/mesameen/micro-app/metadata/internal/handler/grpc"
	"github.com/mesameen/micro-app/metadata/internal/repository/inmemory"
	"github.com/mesameen/micro-app/pkg/discovery"
	"github.com/mesameen/micro-app/pkg/discovery/consulimpl"
	"github.com/mesameen/micro-app/pkg/logger"
	"github.com/mesameen/micro-app/src/api/gen"
	"google.golang.org/grpc"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8091, "API handler port")
	flag.Parse()
	log.Println("Starting the movie metadata service")
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	registry, err := consulimpl.NewRegistry("localhost:8500")
	if err != nil {
		logger.Panicf("unable to connect to service registry. Error: %v", err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		logger.Panicf("Failed to register instance %s of service %s to service registry", instanceID, serviceName)
	}
	go func() {
		// continusous interval to post health status
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := registry.ReportHealthyState(ctx, instanceID, serviceName); err != nil {
					logger.Panicf("Failed to report health status of instance %s of service %s to service Registry. Error: %v", instanceID, serviceName, err)
				}
			}
		}
	}()
	// deregistering instance of the metada service
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := inmemory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8091")
	if err != nil {
		logger.Panicf("Failed to listen on 8091. Error: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	go func() {
		logger.Infof("%s service is up and running on :8091", serviceName)
		if err := srv.Serve(lis); err != nil {
			logger.Panicf("Failed to sert grpc server. Error: %v", err)
		}
	}()
	<-ctx.Done()
	// do graceful shutdown
	srv.GracefulStop()
	// Commented out to start the service as HTTP server
	// router := gin.Default()
	// router.GET("/metadata", h.GetMetadata)
	// server := http.Server{
	// 	Addr:    fmt.Sprintf(":%d", port),
	// 	Handler: router,
	// }
	// go func() {
	// 	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Panicf("Failed to start the server. Error: %v", err)
	// 	}
	// }()
	// logger.Infof("Metadata service is up and running on: %d", port)
	// <-ctx.Done()
	// timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	// defer timeoutCancel()
	// if err := server.Shutdown(timeoutCtx); err != nil {
	// 	logger.Errorf("Failed to shutdown server. Error: %v", err)
	// 	return
	// }

	logger.Infof("Server shutdown gracefully")
}
