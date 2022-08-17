package xmonitor

import "net/http"

// From https://github.com/zsais/go-gin-prometheus/blob/2199a42d96c1d40f249909ed2f27d42449c7fc94/middleware.go#L397
func calcRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
