package util

import (
	"net/http"

	"github.com/kataras/muxie"
	"golang.org/x/exp/slices"
)

func NewCorsMiddleware(origin []string, allowCredentials bool, handler http.Handler) muxie.Wrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if origin != nil && slices.Contains(origin, req.Host) {
				res.Header().Set("Access-Control-Allow-Origin", req.Host)
			}
			if allowCredentials {
				res.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			next.ServeHTTP(res, req)
		})
	}
}
