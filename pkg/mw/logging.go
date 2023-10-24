package mw

import (
	"blackbox-v2/pkg/logging"
	"fmt"
	"net/http"
)

func LogIt(msg string) {
	logging.BlackboxLogger.Println(msg)
	logging.BlackboxCLILogger.Println(msg)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStr := fmt.Sprintf(
			"Request received: %s %s %s",
			r.Method,
			r.URL.Path,
			r.Proto,
		)
		LogIt(requestStr)
		next.ServeHTTP(w, r)
	})
}
