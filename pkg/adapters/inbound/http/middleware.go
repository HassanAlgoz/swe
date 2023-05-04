package http

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Compose(middlewares ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

func WithLogging(logger zerolog.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("Started %s %s", r.Method, r.URL.Path)
			next(w, r)
			logger.Printf("Completed %s %s", r.Method, r.URL.Path)
		}
	}
}

func WithRequestDeduplication() Middleware {
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

func WithFeatureFlagsCheck(requiredFeatureFlags []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for i := range requiredFeatureFlags {
				if !viper.GetBool(fmt.Sprintf("app.features.%s", requiredFeatureFlags[i])) {
					w.WriteHeader(http.StatusForbidden)
				}
			}
			next(w, r)
		}
	}
}

func WithRequestMethodAndHeaderAssertion(allowedMethods []string, requiredHeaders []string) Middleware {
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
