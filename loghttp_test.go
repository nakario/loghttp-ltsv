package loghttpltsv

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type test struct {
	got  string
	want string
}

func splitlv(lv string) (string, string) {
	split := strings.SplitN(lv, ":", 2)
	return split[0], split[1]
}

func parseLog(log string) map[string]string {
	lvs := strings.Split(log, "\t")
	result := make(map[string]string)
	for _, lv := range lvs {
		l, v := splitlv(lv)
		result[l] = v
	}
	return result
}

func TestRoundTrip(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		s := req.RequestURI
		switch {
		case strings.Contains(s, "?"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(req.URL.Query().Encode()))
		case strings.HasPrefix(s, "/e"):
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	buf := new(bytes.Buffer)
	cli := http.Client{
		Transport: NewTransport(buf),
	}

	cli.Get(ts.URL + "/foo?def=5&abc=3")
	cli.Post(ts.URL+"/e", "text/plain", bytes.NewBufferString("a"))
	cli.Get(ts.URL + "/bar")

	log := buf.String()
	logs := strings.Split(log, "\n")

	log0 := parseLog(logs[0])
	log1 := parseLog(logs[1])
	log2 := parseLog(logs[2])
	tests := []test{
		{log0["host"], "127.0.0.1"},
		{log0["forwardedfor"], "-"},
		{log0["req"], "GET /foo?def=5&abc=3 HTTP/1.1"},
		{log0["method"], "GET"},
		{log0["uri"], "/foo?def=5&abc=3"},
		{log0["status"], "200"},
		{log0["size"], "11"},
		{log0["referer"], "-"},
		{log0["ua"], "-"},
		{log0["vhost"], ts.Listener.Addr().String()},
		{log1["req"], "POST /e HTTP/1.1"},
		{log1["status"], "500"},
		{log1["size"], "5"},
		{log2["req"], "GET /bar HTTP/1.1"},
		{log2["status"], "404"},
		{log2["size"], "0"},
	}

	for i, v := range tests {
		t.Run(fmt.Sprint("case", i), func(t *testing.T) {
			if v.got != v.want {
				t.Errorf("got %q, want %q", v.got, v.want)
			}
		})
	}

	if _, ok := log0["response_time"]; !ok {
		t.Errorf("log0 doesn't contain 'response_time'")
	}
	if _, ok := log0["apptime"]; !ok {
		t.Errorf("log0 doesn't contain 'apptime")
	}
	rt := log1["response_time"]
	rtf, err := strconv.ParseFloat(rt, 64)
	if err != nil {
		t.Errorf("failed to parse 'request_time'")
	}
	if rtf < 1.0 {
		t.Errorf("'response_time' isn't correct")
	}
}
