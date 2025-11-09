package grapher

import (
	"encoding/json"
	"net/http"

	"github.com/graph-gophers/graphql-go"
)

func NewHandler(schema *graphql.Schema, opts ...HandlerOpt) *Handler {
	h := &Handler{schema: schema}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

var _ http.Handler = (*Handler)(nil)

type Handler struct {
	schema   *graphql.Schema
	explorer http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.explorer == nil {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		} else {
			h.explorer.ServeHTTP(w, r)
		}
	case http.MethodPost:
		var params struct {
			Query         string         `json:"query"`
			OperationName string         `json:"operationName"`
			Variables     map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := h.schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type HandlerOpt func(*Handler)

func WithExplorer(explorer http.Handler) HandlerOpt {
	return func(h *Handler) {
		h.explorer = explorer
	}
}
