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

func ExtractRequestId(r *http.Request) (string, error) {
	s := strings.Split(H(r, HeaderRequestId), " ")
	if len(s) != 2 {
		return "", fmt.Errorf(`invalid header "%s"`, HeaderRequestId)
	}
	return s[1], nil
}
