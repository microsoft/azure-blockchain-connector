package main

import (
	"fmt"
	"net/http"
)

// See flags.go for detail
const (
	defaultLocalAddr   = "127.0.0.1:3100"
	defaultAuthSvcAddr = "127.0.0.1:3101"
)

func main() {
	c := whenCancelling(nil)

	p := newProxyFromFlags()
	p.Init()
	check(p.TestConn())

	check(http.ListenAndServe(p.Local, p))
	fmt.Println("Tunneling:", p.Local, "->", p.Remote)

	<-c
}
