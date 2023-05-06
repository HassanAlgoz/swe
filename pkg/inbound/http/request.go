package http

import (
	"fmt"
	"net/http"
	"strings"
)

type Header string

const (
	HeaderAuthorization Header = "Authorization"
	HeaderRequestId     Header = "X-Request-Id"
)

// H is just a type-safe replacement for r.Header.Get(key)
func H(r *http.Request, key Header) string {
	return r.Header.Get(string(key))
}

func GetAuthToken(r *http.Request) (string, error) {
	s := strings.Split(H(r, HeaderAuthorization), " ")
	if len(s) != 2 {
		return "", fmt.Errorf(`invalid header "%s"`, HeaderAuthorization)
	}
	return s[1], nil
}

func GetRequestId(r *http.Request) (string, bool) {
	reqId := H(r, HeaderRequestId)
	if len(reqId) == 0 {
		return "", false
	}
	return reqId, true
}
