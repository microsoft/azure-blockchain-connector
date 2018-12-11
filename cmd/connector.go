package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

const (
	defaultLocalAddr   = "127.0.0.1:3100"
	defaultAuthSvcAddr = "127.0.0.1:3101"
)

func main() {
	var cancellation = make(chan os.Signal)

	var p = NewProxyFromFlags()

	p.Init()
	p.TestConn()

	fmt.Println("Requests will be transport to: " + p.Remote)
	fmt.Println("Listen on " + p.Local)

	err := http.ListenAndServe(p.Local, p.Handler())
	if err != nil {
		fmt.Println("Error on listening: ", err)
		os.Exit(-2)
	}

	fmt.Println("Connector started")

	signal.Notify(cancellation, os.Interrupt)
	<-cancellation
	fmt.Println("Cancelling...")
}
