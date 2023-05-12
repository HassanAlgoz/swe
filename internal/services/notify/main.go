package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hassanalgoz/swe/internal/services/notify/consumer"
	"github.com/hassanalgoz/swe/internal/services/notify/controller"
	"github.com/hassanalgoz/swe/internal/services/notify/producer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		cancel()
	}()

	p, err := producer.New(ctx)
	if err != nil {
		panic(err)
	}

	ctrl := controller.New(p)

	// Initialize
	c := consumer.New(ctx, ctrl)
	c.Start(done)
}
