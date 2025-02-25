package pkg

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/plugin/environ"
	"github.com/drone/drone-go/plugin/logger"
)

type (
	handler struct {
		internal http.Handler
	}
)

func Handler(secret string, plugin environ.Plugin, logs logger.Logger) http.Handler {
	return &handler{
		environ.Handler(secret, plugin, logs),
	}
}

func (p *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "request", r)
	newRequest := r.WithContext(ctx)
	p.internal.ServeHTTP(w, newRequest)
}
