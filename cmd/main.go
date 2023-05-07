package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/hassanalgoz/swe/pkg/config"
	"github.com/hassanalgoz/swe/pkg/inbound/http"
	"github.com/hassanalgoz/swe/pkg/inbound/kafka"
	"github.com/spf13/viper"
)

func init() {
	config.SetupConfig()
}

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
	kc := kafka.NewConsumer(
		ctx,
		act,
		strings.Join(viper.GetStringSlice("kafka.brokers"), ","),
		viper.GetString("kafka.group"),
		viper.GetStringSlice("kafka.topics"),
	)
	go kc.Start(doneChan)

	httpServer := http.NewServer(ctx, act)
	if err := httpServer.Listen(fmt.Sprintf("%s:%s", viper.GetString("http.host"), viper.GetString("http.port"))); err != nil {
		panic(err)
	}
}
