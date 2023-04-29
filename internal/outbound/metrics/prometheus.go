package metrics

import "github.com/prometheus/client_golang/prometheus"

var MyCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "my_counter",
	Help: "This is my counter",
})

func init() {
	prometheus.MustRegister(MyCounter)
}
