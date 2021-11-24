package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func UUIDContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uuidContext := context.WithValue(ctx, "uuid", uuid.New())

		r = r.WithContext(uuidContext)
		next.ServeHTTP(rw, r)
	})
}
