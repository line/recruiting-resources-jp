package filter

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"todo-example/server/pkg/common"
	"todo-example/server/pkg/db"
)

// Init initializes request context objects
// which will be used through filters and handler.
// request contexts:'dbClient', 'logger'
func Init() Decorator {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New().String()
			logger := &common.Log{RequestId: requestId}
			logger.Configure()
			logger.Info(log.Fields{}, "Init db client")
			dbClient := db.New()
			ctx := context.WithValue(r.Context(), "dbClient", dbClient)
			ctx = context.WithValue(ctx, "logger", logger)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}
