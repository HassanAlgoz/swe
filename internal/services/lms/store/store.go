package store

import (
	"sync"

	port "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/infra/database"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance port.Querier
)

type Adapter struct {
	port.Querier
}

// New dbname must not exceeed 31 characters
func New(dbname string) port.Querier {
	if len(dbname) > 31 {
		panic("dbname must not exceeed 31 characters")
	}
	switch viper.GetString("env") {
	default:
		once.Do(func() {
			instance = &Adapter{
				port.New(database.Get(dbname)),
			}
		})

	case "test":
		return &Adapter{
			port.New(database.Get(dbname)),
		}
	}
	return instance
}
