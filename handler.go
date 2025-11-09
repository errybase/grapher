package grapher

import (
	"context"
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
	ctxFn    ctxFn
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

		ctx := r.Context()
		if h.ctxFn != nil {
			ctx = h.ctxFn(r)
		}

		res := h.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

func WithContext(fn ctxFn) HandlerOpt {
	return func(h *Handler) {
		h.ctxFn = fn
	}
}

type ctxFn func(*http.Request) context.Context
