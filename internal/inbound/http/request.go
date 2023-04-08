package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func requireHeaders(w http.ResponseWriter, r *http.Request, headers []string) map[string]string {
	m := map[string]string{}
	missing := []string{}
	for i := range headers {
		k := headers[i]
		v := r.Header.Get(k)
		if v == "" {
			missing = append(missing, k)
			continue
		}
		m[k] = v
	}
	if len(missing) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Error: Error{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Missing Headers: %v", missing),
			},
		})
	}
	return m
}
