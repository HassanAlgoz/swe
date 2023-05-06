package port

import (
	"net/http"

	inbound "github.com/hassanalgoz/swe/pkg/inbound/http"
)

type endpointOptions struct {
	RequiredFeatureFlags []string // see example: "money_transfer" (according to ./etc/config.yaml:app.features)
	RequiredHeaders      []inbound.Header
}

func (s *service) registerEndpoint(methods []string, pattern string, handlerFunc http.HandlerFunc, opts *endpointOptions) {
	headers := make([]string, len(opts.RequiredHeaders))
	for i := range opts.RequiredHeaders {
		headers[i] = string(opts.RequiredHeaders[i])
	}
	flags := make([]string, len(opts.RequiredFeatureFlags))
	for i := range opts.RequiredFeatureFlags {
		flags[i] = string(opts.RequiredFeatureFlags[i])
	}
	mws := []inbound.Middleware{
		inbound.WithMetrics(),
		inbound.WithLogging(),
		inbound.WithRequestDeduplication(),
		inbound.WithRequestMethodAndHeaderAssertion(
			methods,
			headers,
		),
		inbound.WithFeatureFlagsCheck(flags),
	}
	s.mux.HandleFunc(pattern, inbound.Compose(mws...)(handlerFunc))
}
