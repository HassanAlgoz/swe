package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hassanalgoz/swe/internal/services/notify/consumer"
	"github.com/hassanalgoz/swe/internal/services/notify/controller"
	"github.com/hassanalgoz/swe/internal/services/notify/producer"
	"github.com/spf13/viper"
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

	p, err := producer.New(ctx,
		viper.GetString("kafka.bootstrap.servers"),
		viper.GetString("services.notify.producer.topic"),
	)
	if err != nil {
		panic(err)
	}

	ctrl := controller.New(p)

	// Initialize
	c := consumer.New(ctx, ctrl,
		viper.GetString("kafka.bootstrap.servers"),
		viper.GetString("services.notify.consumer.group.id"),
		viper.GetStringSlice("services.notify.consumer.topics"),
	)
	c.Start(done)
}
