// Package config contains all setup
package config

import (
	"net/http"

	"github.com/bookfire/controller"
)

// ChainMiddleware group the middlewares in each request
func ChainMiddleware(mw ...controller.Middleware) controller.Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}
