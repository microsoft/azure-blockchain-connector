package main

import (
	"azure-blockchain-connector/aad"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

// needWebviewForAuthCode checks if the os.Args include the 'authcode-webview' flag,
// which enable this tool to open a window to help the user doing authenticates.
func needWebviewForAuthCode() bool {
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, flagAuthCodeWebview) {
			return true
		}
	}
	return false
}

func main() {

	// This block will parse essential args and invoke a webview to print auth code to stdout.
	// Notice that this part should not be invoked by user-entered args, but an implementation internal agreement.
	if needWebviewForAuthCode() {
		var authURL string
		flag.StringVar(&authURL, flagAuthCodeWebview, "", "")
		flag.Parse()

		aad.AuthCodeWebview(authURL, os.Stdout)
		return
	}

	// According proxy is generated from args and then authentication will be requested.
	// Once the authentication is granted, this tool will listen to a port to proxy traffic.
	c := whenCancelling(nil)

	p := newProxyFromFlags()
	check(p.Provider.RequestAccess())
	p.ConfigureClient()

	go func() {
		check(http.ListenAndServe(p.Local, p))
	}()
	fmt.Println("Tunnel:", p.Local, "->", p.Remote)

	<-c
}
