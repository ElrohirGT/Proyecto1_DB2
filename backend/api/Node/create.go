package node

import "net/http"

func NewCreateNodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
