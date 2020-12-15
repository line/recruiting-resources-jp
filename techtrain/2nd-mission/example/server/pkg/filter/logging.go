package filter

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"todo-example/server/pkg/common"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(body []byte) (int, error) {
	lrw.body = body
	return lrw.ResponseWriter.Write(body)
}

// Logger records client access requests
func Logger() Decorator {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := r.Context().Value("logger").(*common.Log)
			start := time.Now()
			var host string
			// if api is behind a proxy
			if r.Header.Get("X-Forwarded-For") != "" {
				host = r.Header.Get("X-Forwarded-For")
			} else {
				host, _, _ = net.SplitHostPort(r.RemoteAddr)
			}
			path := r.Method + " " + r.URL.Path
			if r.URL.RawQuery != "" {
				path = path + "?" + r.URL.RawQuery
			}
			route := mux.CurrentRoute(r)
			routePath := ""
			if route != nil {
				routePath, _ = route.GetPathTemplate()
				routePath = common.GetPathWithoutVersion(routePath)
			}
			defer r.Body.Close()
			reader := common.GetRequestBody(r)
			requestBody, _ := ioutil.ReadAll(reader)
			logFields := log.Fields{
				"type":        "AccessLog",
				"requestBody": string(requestBody),
			}
			var statusCode int
			lrw := &loggingResponseWriter{w, http.StatusOK, []byte{}}
			h.ServeHTTP(lrw, r)
			logFields["responseBody"] = string(lrw.body)
			statusCode = lrw.statusCode
			logFields["duration"] = time.Since(start)
			logger.Info(
				logFields,
				"%v %v %v %v %v %v",
				host,
				path,
				r.Proto,
				statusCode,
				r.ContentLength,
				r.Header.Get("User-Agent"),
			)
		})
	}
}
