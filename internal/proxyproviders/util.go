package proxyproviders

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"time"
)

var EnablePrintToken bool

func printToken(tok interface{}) {
	if !EnablePrintToken || tok == nil {
		return
	}
	elm := reflect.ValueOf(tok).Elem()
	if elm.Kind() != reflect.Struct {
		return
	}
	for _, name := range []string{"AccessToken", "RefreshToken"} {
		f := elm.FieldByName(name)
		if f.IsValid() && f.String() != "" {
			fmt.Println(name+":", f.String())
		}
	}
}

// PollAuthHeader starts a test server and prints auth headers by polling.
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
