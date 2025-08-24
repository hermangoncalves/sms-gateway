package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hermangoncalves/sms-gateway/internal/config"
	"github.com/hermangoncalves/sms-gateway/internal/polling"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load confi;g %v", err)
	}

	fmt.Println(cfg)

	poller := polling.NewPoller(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down... Sayonara!")
		cancel()
	}()

	poller.Start(ctx)
}
