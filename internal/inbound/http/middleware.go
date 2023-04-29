package http

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

type middlewareOptions struct {
	RequiredFeatureFlags []FeatureFlag
	RequiredHeaders      []Header
}

func (s *Server) registerEndpoint(methods []string, pattern string, handlerFunc http.HandlerFunc, opts *middlewareOptions) {
	headers := make([]string, len(opts.RequiredHeaders))
	for i := range opts.RequiredHeaders {
		headers[i] = string(opts.RequiredHeaders[i])
	}
	flags := make([]string, len(opts.RequiredFeatureFlags))
	for i := range opts.RequiredFeatureFlags {
		flags[i] = string(opts.RequiredFeatureFlags[i])
	}
	mws := []middleware{
		withMetrics(),
		s.withLogging(),
		withRequestDeduplication(),
		withRequestMethodAndHeaderAssertion(
			methods,
			headers,
		),
		withFeatureFlagsCheck(flags),
	}
	s.mux.HandleFunc(pattern, compose(mws...)(handlerFunc))
}

func compose(middlewares ...middleware) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

func (s *Server) withLogging() middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s.logger.Printf("Started %s %s", r.Method, r.URL.Path)
			next(w, r)
			s.logger.Printf("Completed %s %s", r.Method, r.URL.Path)
		}
	}
}

func withRequestDeduplication() middleware {
	requests := make(map[string]bool)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			requestID := H(r, HeaderRequestId)
			if requests[requestID] {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			requests[requestID] = true

			next(w, r)
		}
	}
}

func withFeatureFlagsCheck(requiredFeatureFlags []string) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for i := range requiredFeatureFlags {
				if !viper.GetBool(requiredFeatureFlags[i]) {
					w.WriteHeader(http.StatusForbidden)
				}
			}
			next(w, r)
		}
	}
}

func withRequestMethodAndHeaderAssertion(allowedMethods []string, requiredHeaders []string) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, method := range allowedMethods {
				if r.Method == method {
					for _, header := range requiredHeaders {
						if r.Header.Get(header) == "" {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprintf(w, "Missing required header: %s", header)
							return
						}
					}

					next(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Invalid request method. Allowed methods: %v", allowedMethods)
		}
	}
}
