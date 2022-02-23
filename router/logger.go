package router

import (
	"bytes"
	"log"
	"fmt"
	"time"
	"net/http"
	"strings"
)

// LoggingResponseWriter will encapsulate a standard ResponseWritter with a copy of its statusCode
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// ResponseWriterWrapper is supposed to capture statusCode from ResponseWriter
func ResponseWriterWrapper(w http.ResponseWriter) *LoggingResponseWriter {
	var buf bytes.Buffer
	return &LoggingResponseWriter{w, http.StatusOK, buf}
}

// WriteHeader is a surcharge of the ResponseWriter method
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(buf []byte) (int, error) {
	lrw.body.Write(buf)
	return (lrw.ResponseWriter).Write(buf)
}

// Logger is a gorilla/mux middleware to add log to the API
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := ResponseWriterWrapper(w)

		start := time.Now()
			//start.Format(time.RFC3339),
		// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 286.219µs
		fmt.Printf("\n")
		log.Printf("%s http://%s/%s %s Content-Length %d From %s\n",
			r.Method,
			r.Host,
			r.RequestURI,
			r.Proto, // string "HTTP/1.1"
			r.ContentLength,
			r.RemoteAddr)
		for k, v := range r.Header {
			if !strings.Contains(k, "Xpress-") {
				log.Printf("%s: %v\n", k, v)
			}
		}

		inner.ServeHTTP(wrapper, r)

		log.Printf("==== Response: %d %s\n", wrapper.statusCode, time.Since(start))
		respHeaders := w.Header()
		for k, v := range respHeaders {
			log.Printf("%s: %v\n", k, v)
		}
		log.Printf("%s\n", &wrapper.body)
	})
}
