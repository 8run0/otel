package middleware

import (
	"context"
	"net/http"

	"github.com/8run0/otel/backend/pkg/otel"
)

// HTTP middleware setting a value on the request context
func WithOTELTools(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, tools := otel.NewTools(r.Context(), "demo-otels")
		ctx := context.WithValue(r.Context(), otel.ToolsCtxKey, tools)
		defer tools.Cleanup()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SpanHttp(spanName string, hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tools := ctx.Value(otel.ToolsCtxKey).(*otel.Tools)
		ctx, span := tools.Tracer.Start(ctx, spanName)
		defer span.End()
		r = r.WithContext(ctx)
		hf(w, r)
	}
}
