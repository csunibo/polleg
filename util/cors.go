package util

import (
	"net/http"

	"github.com/kataras/muxie"
)

func NewCorsMiddleware(origin []string, allowCredentials bool, handler http.Handler) muxie.Wrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			res.Header().Set("Access-Control-Allow-Credentials", "true")
			// if origin != nil && slices.Contains(origin, req.Host) {
			// }
			// if allowCredentials {
			// }

			next.ServeHTTP(res, req)
		})
	}
}
