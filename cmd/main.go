package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hassanalgoz/swe/internal/actions"
	"github.com/hassanalgoz/swe/internal/contexts/transfer"
	"github.com/hassanalgoz/swe/internal/contexts/user"
	"github.com/hassanalgoz/swe/internal/inbound/http"
	"github.com/hassanalgoz/swe/internal/inbound/kafka"
	"github.com/hassanalgoz/swe/internal/outbound/search"
	"github.com/hassanalgoz/swe/pkg/config"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

	// Outbound
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/bank")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conn, err := grpc.Dial(viper.GetString("ports.search.address"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	svcSearch := search.NewSearchServiceClient(conn)

	// Application Layer
	transferContext := transfer.NewContext(db)

	userContext := user.NewContext(db, svcSearch)

	act := actions.New(
		ctx,
		transferContext,
		userContext,
	)

	// Inbound
	kc := kafka.NewConsumer(ctx, act, "localhost:9001", "mygroup", []string{"topic1"})
	go kc.Start(doneChan)

	httpServer := http.NewServer(ctx, act)
	if err = httpServer.Listen(":8080"); err != nil {
		panic(err)
	}
}
