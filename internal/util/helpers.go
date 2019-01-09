package util

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

// PollAuthHeader starts a test server and prints auth headers by polling.
// e.g p.Params.Remote = proxyproviders.PollAuthHeader(p.Params.Local, time.Minute)
func PollAuthHeader(local string, interval time.Duration) string {
	cnt := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(cnt, r.Header.Get("Authorization"))
		cnt++
	}))

	go func() {
		for {
			_, _ = http.Get("http://" + local)
			time.Sleep(interval)
		}
	}()
	u, _ := url.Parse(s.URL)
	return u.Host
}
