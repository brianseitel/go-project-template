package middleware

import (
	"net/http"
	"time"
)

// TimingMiddleware is a set of middleware that helps track how long each request
// takes. First, it identifies the current timestamp when the request begins, then
// it executes the request, then it calculates the elapsed time. Finally, it takes
// the total elapsed time and adds it to the X-Timing header.
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Current time
		start := time.Now()
		// Continue executing the request
		next.ServeHTTP(w, r)
		// This gets current time again and subtracts the start time from it,
		// which should give us the total elapsed time of the above request.
		elapsed := time.Since(start)

		// Add this to our header.
		w.Header().Add("X-Timing", elapsed.String())
	})
}
