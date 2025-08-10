package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AugustSerenity/order-service/internal/cache"
	"github.com/AugustSerenity/order-service/internal/config"
	"github.com/AugustSerenity/order-service/internal/handler"
	"github.com/AugustSerenity/order-service/internal/kafka/consumer"
	"github.com/AugustSerenity/order-service/internal/service"
	"github.com/AugustSerenity/order-service/internal/storage"
	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 5 * time.Second

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()

	cfg := config.ParseConfig(*configPath)

	db := storage.InitDB(cfg.DB)
	defer storage.CloseDB(db)

	storage := storage.New(db)
	cache := cache.NewCache()
	srv := service.NewOrderService(cache, storage)
	h := handler.New(srv)

	address := []string{"kafka:9092"}
	topic := "order"
	groupID := "order-service-group"

	cons, err := consumer.NewConsumer(srv, address, topic, groupID)
	if err != nil {
		logrus.Fatalf("failed to create consumer: %v", err)
	}
	go cons.Start()

	s := http.Server{
		Addr:         cfg.Server.Address,
		Handler:      h.Route(),
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("starting server on", cfg.Server.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutting down consumer and server...")

	if err := cons.Stop(); err != nil {
		logrus.Errorf("error stopping consumer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
}
