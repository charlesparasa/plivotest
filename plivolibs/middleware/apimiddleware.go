package middleware

import (
	"fmt"
	"github.com/charlesparasa/plivotest/plivolibs/logger"
	"net/http"
	"time"
	"github.com/felixge/httpsnoop"
)

//constructHTTPLogMessage
func constructHTTPLog(r *http.Request, m httpsnoop.Metrics, duration time.Duration) string {
	ctx := r.Context().Value(APICtx)
	if ctx != nil {
		tCtx := ctx.(APIContext)
		return fmt.Sprintf("|%s|%s|%s|%s|%s|%d|%d|%s|%s|",
			// Cannot modify original request/obtain apiContext through gorilla context, hence we won't get the apiContext data from the request object.
			tCtx.Name,
			"correlationId="+tCtx.CorrelationID+":requestId="+tCtx.RequestID,
			r.RemoteAddr,
			r.Method,
			r.URL,
			m.Code,
			m.Written,
			r.UserAgent(),
		)
	}
	return fmt.Sprintf("|%s|%s|%s|%d|%d|%s|%s|",
		// Cannot modify original request/obtain apiContext through gorilla context, hence we won't get the apiContext data from the request object.
		r.RemoteAddr,
		r.Method,
		r.URL,
		m.Code,
		m.Written,
		r.UserAgent(),
		duration,
	)

}

//logRequest logs each HTTP incoming Requests
func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m := httpsnoop.CaptureMetrics(handler, w, r)
		logger.HTTPLog(constructHTTPLog(r, m, time.Since(start)))
	})
}
