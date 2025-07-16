package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/mesameen/micro-app/rating/internal/controller"
	grpcHandler "github.com/mesameen/micro-app/rating/internal/handler/grpc"
	"github.com/mesameen/micro-app/rating/internal/ingester/kafkautil"
	"github.com/mesameen/micro-app/rating/internal/repository/mysql"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/discovery"
	"github.com/mesameen/micro-app/src/pkg/discovery/consulimpl"
	"github.com/mesameen/micro-app/src/pkg/logger"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

const serviceName = "rating"

func main() {
	log.Println("Starting the rating service")
	f, err := os.Open("default.yaml")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Panic(err)
	}
	fmt.Printf("cfg : %+v\n", cfg)

	err = logger.Init()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	registry, err := consulimpl.NewRegistry(cfg.ServiceDiscovery.Consul.Address)
	if err != nil {
		logger.Panicf("Failed to connect to service registry. Error: %v", err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("rating:%d", cfg.API.Port)); err != nil {
		logger.Panicf("Failed to register instance %s of service %s to service registry", instanceID, serviceName)
	}
	go func() {
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
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo, err := mysql.New()
	if err != nil {
		logger.Panicf("%v", err)
	}
	ingester, err := kafkautil.NewIngester("localhost", "rating-service", "ratings")
	if err != nil {
		logger.Panicf("Failed to connect to kafka ingestor. Error:%v", err)
	}
	ctrl := controller.New(repo, ingester)
	go func() {
		if err := ctrl.StartIngestion(ctx); err != nil {
			logger.Panicf("Failed to start ingesting kafka events. Error: %v", err)
		}
	}()
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.API.Port))
	if err != nil {
		logger.Panicf("Failed to listen on 8092. Error: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	go func() {
		logger.Infof("%s service is up and running on :%d", serviceName, cfg.API.Port)
		if err := srv.Serve(lis); err != nil {
			logger.Panicf("Failed to sert grpc server. Error: %v", err)
		}
	}()
	<-ctx.Done()
	// do graceful shutdown
	srv.GracefulStop()
	// Commented out to start the service as HTTP server
	// router := gin.Default()
	// router.GET("/rating", h.GetRatings)
	// router.PUT("/rating", h.SaveRatings)
	// server := http.Server{
	// 	Addr:    fmt.Sprintf(":%d", port),
	// 	Handler: router,
	// }
	// go func() {
	// 	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Panicf("Failed to start the server. Error: %v", err)
	// 	}
	// }()
	// logger.Infof("Rating service is up and running on: %d", port)
	// <-ctx.Done()
	// timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	// defer timeoutCancel()
	// if err := server.Shutdown(timeoutCtx); err != nil {
	// 	logger.Errorf("Failed to shutdown server. Error: %v", err)
	// 	return
	// }

	logger.Infof("Server shutdown gracefully")
}
