package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, r)
		msg := fmt.Sprintf(
			"method: %s, url: %s, duration: %d milliseconds",
			r.Method,
			r.RequestURI,
			time.Since(start)/1e6,
		)
		log.Println(msg)
	})
}
