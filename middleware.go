package main

import "net/http"

type Middleware []http.Handler

func (m *Middleware) Add(handler http.Handler) {
	*m = append(*m, handler)
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//wrap the supplied responseWriter
	mw := NewMiddlewareReponseWriter(w)

	//loop through all of the registered handlers
	for _, handler := range m {
		//call handler with our MiddlewareReponseWriter
		handler.ServeHTTP(mw, r)
		if mw.written {
			return
		}
	}
	http.NotFound(w, r)
}

type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

func NewMiddlewareReponseWriter(w http.ResponseWriter) *MiddlewareResponseWriter {
	return &MiddlewareResponseWriter{
		ResponseWriter: w,
	}
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(bytes)
}

func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}
