package search

import (
	"context"
	"fmt"
	"sync"

	"github.com/hassanalgoz/swe/internal/common"
	pb "github.com/hassanalgoz/swe/internal/outbound/grpc/search/pb"
	"github.com/hassanalgoz/swe/internal/outbound/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Client struct {
	client pb.SearchServiceClient
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
		instance.client = pb.NewSearchServiceClient(conn)
	})

	if err != nil {
		log.Fatal().Msgf("failed to instantiate search service client: %v", err)
	}

	return instance
}

func (c *Client) GetSearchResults(ctx context.Context, req *pb.SearchRequest) ([]common.SearchResult, error) {
	resp, err := c.client.GetSearchResults(ctx, req)
	if err != nil {
		return nil, err
	}
	// map to domain entity
	results1 := resp.GetResults()
	results2 := make([]common.SearchResult, len(results1))
	for i := range results1 {
		results2[i] = common.SearchResult{
			Title: results1[i].GetTitle(),
			URL:   results1[i].GetUrl(),
		}
	}
	return results2, nil
}
