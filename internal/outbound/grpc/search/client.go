package search

import (
	"context"
	"fmt"
	"sync"

	"github.com/hassanalgoz/swe/internal/ent"
	"github.com/hassanalgoz/swe/internal/outbound/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Client struct {
	client SearchServiceClient
}

var (
	once     sync.Once
	instance = &Client{}
)

func Get() *Client {
	var err error
	log := logger.Get()

	// Create the singleton instance of DB
	once.Do(func() {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", viper.GetString("grpc.search.host"), viper.GetString("grpc.search.port")), grpc.WithInsecure())
		if err != nil {
			log.Fatal().Msgf("failed to connect to server: %v", err)
		}
		instance.client = NewSearchServiceClient(conn)
	})

	if err != nil {
		log.Fatal().Msgf("failed to instantiate search service client: %v", err)
	}
	return instance
	// Learning Notes
	//
	// this code implements the **Singleton** design pattern
	//
	// The `Get()` function creates a single instance of the `Client` struct, which is returned whenever the function is called.
	// This ensures that only one instance of the `Client` struct is created throughout the lifetime of the program, which can help to reduce resource consumption and improve performance.
	//
	// The `sync.Once` type is used to ensure that the initialization of the `Client` struct is performed only once, even if multiple concurrent calls to the `Get()` function are made.
	// This is achieved by using a `sync.Once` instance to perform the initialization code within a function passed to the `Do()` method. The `Do()` method ensures that the function is
	// only executed once, regardless of how many times the `Get()` function is called.
}

func (c *Client) GetSearchResults(ctx context.Context, req *SearchRequest) ([]ent.SearchResult, error) {
	resp, err := c.client.GetSearchResults(ctx, req)
	if err != nil {
		return nil, err
	}
	// map to domain entity
	results1 := resp.GetResults()
	results2 := make([]ent.SearchResult, len(results1))
	for i := range results1 {
		results2[i] = ent.SearchResult{
			Title: results1[i].GetTitle(),
			URL:   results1[i].GetUrl(),
		}
	}
	return results2, nil
	// Learning Notes:
	//
	// this code implements the **Adapter** design pattern
	//
	// The `GetSearchResults` method adapts the `resp` object returned from `c.client.GetSearchResults` method to a new object of type `[]ent.SearchResult`.
	// It achieves this by mapping the properties of the objects returned by the external service to a new object of the required type `ent.SearchResult`.
	//
	// This pattern is often used in situations where you have an existing interface that is not compatible with the interface required by the client,
	// and you need to create an adapter that acts as a bridge between the two interfaces. In this case, the external service is returning a response that is
	// not directly compatible with the `ent.SearchResult` type used in the client code, and the adapter method maps the response to the desired type.
}
