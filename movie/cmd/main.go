package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/movie/internal/controller/movie"
	metadataService "github.com/mesameen/micro-app/movie/internal/gateway/metadata/http"
	ratingService "github.com/mesameen/micro-app/movie/internal/gateway/rating/http"
	httpHandler "github.com/mesameen/micro-app/movie/internal/handler/http"
	"github.com/mesameen/micro-app/movie/internal/logger"
)

func main() {
	log.Println("Starting the movie service")
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	metadataGateway := metadataService.New("localhost:8091")
	ratingGateway := ratingService.New("localhost:8092")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	router := gin.Default()
	router.GET("/moview", h.GetMovieDetails)

	server := http.Server{
		Addr:    ":8093",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Failed to start the server. Error: %v", err)
		}
	}()
	logger.Infof("Server is up and running on: 8093")
	<-ctx.Done()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeoutCancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		logger.Errorf("Failed to shutdown server. Error: %v", err)
	}
	logger.Infof("Server shutdown gracefully")
}
