package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/hassanalgoz/swe/internal/app"
	"github.com/hassanalgoz/swe/internal/inbound/http"
	"github.com/hassanalgoz/swe/internal/inbound/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt)
	doneChan := make(chan bool, 1)
	go func() {
		<-cancelChan
		cancel()
		doneChan <- true
	}()

	// Application Layer
	act := app.New(ctx)

	// Inbound
	kc := kafka.NewConsumer(ctx, act, "localhost:9001", "mygroup", []string{"topic1"})
	go kc.Start(doneChan)

	httpServer := http.NewServer(ctx, act)
	if err := httpServer.Listen(":8080"); err != nil {
		panic(err)
	}
}
