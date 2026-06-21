package middleware

import (
	"io"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/trafficstats"
	"github.com/gin-gonic/gin"
)

// TrafficStats records app-observable HTTP bytes and estimates TLS/TCP/IP overhead.
func TrafficStats(cfg config.TrafficConfig) gin.HandlerFunc {
	statsCfg := trafficstats.NormalizeConfig(trafficstats.Config{
		Enabled:               cfg.Enabled,
		Source:                cfg.Source,
		TLSRecordPayloadBytes: cfg.TLSRecordPayloadBytes,
		TLSRecordOverhead:     cfg.TLSRecordOverheadBytes,
		TCPIPHeaderBytes:      cfg.TCPIPHeaderBytes,
		TCPPayloadBytes:       cfg.TCPPayloadBytes,
	})
	if !statsCfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		if c == nil || c.Request == nil {
			c.Next()
			return
		}
		counter := trafficstats.NewCounter(statsCfg, c.Request)
		if c.Request.Body != nil {
			c.Request.Body = &trafficCountingReadCloser{ReadCloser: c.Request.Body, counter: counter}
		}
		c.Writer = &trafficResponseWriter{
			ResponseWriter: c.Writer,
			counter:        counter,
			proto:          c.Request.Proto,
		}
		c.Request = c.Request.WithContext(trafficstats.IntoContext(c.Request.Context(), counter))
		c.Next()
		counter.FinalizeResponse(c.Writer.Status(), c.Writer.Header(), c.Request.Proto)
	}
}

type trafficCountingReadCloser struct {
	io.ReadCloser
	counter *trafficstats.Counter
}

func (r *trafficCountingReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if n > 0 {
		r.counter.AddRequestBody(int64(n))
	}
	return n, err
}

type trafficResponseWriter struct {
	gin.ResponseWriter
	counter *trafficstats.Counter
	proto   string
}

func (w *trafficResponseWriter) WriteHeader(statusCode int) {
	w.counter.FinalizeResponse(statusCode, w.Header(), w.proto)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *trafficResponseWriter) WriteHeaderNow() {
	status := w.Status()
	if status <= 0 {
		status = http.StatusOK
	}
	w.counter.FinalizeResponse(status, w.Header(), w.proto)
	w.ResponseWriter.WriteHeaderNow()
}

func (w *trafficResponseWriter) Write(data []byte) (int, error) {
	status := w.Status()
	if status <= 0 {
		status = http.StatusOK
	}
	w.counter.FinalizeResponse(status, w.Header(), w.proto)
	n, err := w.ResponseWriter.Write(data)
	if n > 0 {
		w.counter.AddResponseBody(int64(n))
	}
	return n, err
}

func (w *trafficResponseWriter) WriteString(s string) (int, error) {
	status := w.Status()
	if status <= 0 {
		status = http.StatusOK
	}
	w.counter.FinalizeResponse(status, w.Header(), w.proto)
	n, err := w.ResponseWriter.WriteString(s)
	if n > 0 {
		w.counter.AddResponseBody(int64(n))
	}
	return n, err
}
