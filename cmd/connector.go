package main

import (
	"fmt"
	"net/http"
)

func main() {
	c := whenCancelling(nil)

	p := newProxyFromFlags()

	check(http.ListenAndServe(p.Local, p))
	fmt.Println("Tunneling:", p.Local, "->", p.Remote)

	<-c
}
