// Package loghttpltsv enables a http.Client to log information of
// requests and responses.
// The log format is LTSV, which can be easily profiled with
// tools like alp(https://github.com/tkuchiki/alp).
package loghttpltsv

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Transport writes a ltsv-formatted log for each request and
// response pair.
type Transport struct {
	Transport http.RoundTripper
	w io.Writer
}

// NewTransport returns a new Transport which writes ltsv-formatted
// logs in the given io.Writer.
func NewTransport(w io.Writer) *Transport {
	return &Transport{w: w}
}

// RoundTrip summarizes information of a request and a response and
// writes a line of ltsv-formatted log.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	duration := time.Now().Sub(start).Seconds()
	durStr := fmt.Sprintf("%.3f", duration)

	lvs := []string{
		lv("host", "127.0.0.1"),
		lv("forwardedfor", ""),
		lv("req", fmt.Sprintf("%s %s %s", req.Method, req.URL.RequestURI(), req.Proto)),
		lv("method", req.Method),
		lv("uri", req.URL.RequestURI()),
		lv("status", resp.StatusCode),
		lv("size", resp.ContentLength),
		lv("referer", req.Referer()),
		lv("ua", req.UserAgent()),
		lv("response_time", durStr),
		lv("apptime", durStr),
		lv("vhost", req.Host),
	}

	if t.w != nil {
		_, err := t.w.Write([]byte(ltsv(lvs) + "\n"))
		if err != nil {
			return resp, err
		}
	}
	
	return resp, nil
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
