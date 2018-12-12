package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// check checks if an error is not nil, then print error message and exit.
// Keep the same exit code 2 with the former implementation.
// (notice: flag package use the code 2 to exit, see FlagSet.Parse ExitOnError)
func check(err error, v ...interface{}) {
	if err != nil {
		v = append(v, err)
		log.Print(v...)
		os.Exit(2)
		//log.Fatal(v, err)
	}
}

func whenCancelling(fn func()) chan struct{} {
	c := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)
		<-s
		if fn != nil {
			fn()
		}
		close(c)
	}()
	return c
}

func main() {
	c := whenCancelling(nil)

	p := newProxyFromFlags()
	check(p.Provider.RequestAccess())

	go func() {
		check(http.ListenAndServe(p.Local, p))
	}()
	fmt.Println("Tunneling:", p.Local, "->", p.Remote)

	<-c
}
