package main

import (
	"net/http"

	log "gopkg.in/inconshreveable/log15.v2"
)

type ResponseWriterWrapper struct {
	rw     http.ResponseWriter
	status int
}

func (r *ResponseWriterWrapper) Write(p []byte) (int, error) {
	return r.rw.Write(p)
}

func (r *ResponseWriterWrapper) Header() http.Header {
	return r.rw.Header()
}

func (r *ResponseWriterWrapper) WriteHeader(status int) {
	r.status = status
	r.rw.WriteHeader(status)
}

func accessLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Info("Received request",
			"method", req.Method,
			"path", req.URL.Path,
		)

		rww := ResponseWriterWrapper{rw, 200}
		handler.ServeHTTP(&rww, req)

		log.Info("Processed request",
			"method", req.Method,
			"path", req.URL.Path,
			"status", rww.status,
		)
	})
}
