package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/metadata/internal/controller/metadata"
	httphandler "github.com/mesameen/micro-app/metadata/internal/handler/http"
	"github.com/mesameen/micro-app/metadata/internal/repository/inmemory"
	"github.com/mesameen/micro-app/pkg/discovery"
	"github.com/mesameen/micro-app/pkg/discovery/consulimpl"
	"github.com/mesameen/micro-app/pkg/logger"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8091, "API handler port")
	flag.Parse()
	log.Println("Starting the metadata service")
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
	h := httphandler.New(ctrl)
	router := gin.Default()
	router.GET("/metadata", h.GetMetadata1)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Panicf("Failed to start the server. Error: %v", err)
		}
	}()
	logger.Infof("Metadata service is up and running on: %d", port)
	<-ctx.Done()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeoutCancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		logger.Errorf("Failed to shutdown server. Error: %v", err)
		return
	}
	logger.Infof("Server shutdown gracefully")
}
