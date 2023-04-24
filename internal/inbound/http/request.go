package http

import "net/http"

type FeatureFlag string

const (
	FeatureFlagMoneyTransfer FeatureFlag = "money_transfer"
	FeatureFlagAnalytics     FeatureFlag = "analytics"
)

type Header string

const (
	HeaderAuthorization Header = "Authorization"
	HeaderRequestId     Header = "X-Request-Id"
)

// H is a type-safe replacement for r.Header.Get(key)
func H(r *http.Request, key Header) string {
	return r.Header.Get(string(key))
}
