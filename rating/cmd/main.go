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
	"github.com/mesameen/micro-app/rating/internal/controller/rating"
	httpHandler "github.com/mesameen/micro-app/rating/internal/handler/http"
	"github.com/mesameen/micro-app/rating/internal/logger"
	"github.com/mesameen/micro-app/rating/internal/repository/inmemory"
)

func main() {
	log.Println("Starting the rating service")
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	repo := inmemory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)
	router := gin.Default()
	router.GET("/rating", h.GetRatings)
	router.PUT("/rating", h.SaveRatings)
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Failed to start the server. Error: %v", err)
		}
	}()
	logger.Infof("Server is up and running on: 8090")
	<-ctx.Done()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeoutCancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		logger.Errorf("Failed to shutdown server. Error: %v", err)
	}
	logger.Infof("Server shutdown gracefully")

}
