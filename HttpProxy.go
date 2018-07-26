package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const defaultLocal = "127.0.0.1:3100"

type proxyParams struct {
	local, remote                string
	username, password, certPath string
	pool                         *x509.CertPool
	client                       *http.Client
}

type proxyHandler struct{
	params *proxyParams
}

func (handler proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var params = handler.params

	req.URL.Scheme = "https"
	req.URL.Host = params.remote
	req.Host = ""
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "identity")
	req.SetBasicAuth(params.username, params.password)

	// do request and get response
	response, err := params.client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer response.Body.Close()
	io.Copy(rw, response.Body)
}

func initParameter(params *proxyParams) {
	flag.StringVar(&params.local, "local", "", "Local address to bind to")
	flag.StringVar(&params.remote, "remote", "", "Remote endpoint address")
	flag.StringVar(&params.username, "username", "", "The username you want to login with.")
	flag.StringVar(&params.password, "password", "", "The password you want to login with.")
	flag.StringVar(&params.certPath, "cert", "", "(Optional) File path to root CA")
	flag.Parse()

	if params.local == "" {
		params.local = defaultLocal
	}

	if params.remote == "" || params.username == "" || params.password == "" {
		flag.Usage()
		os.Exit(-1)
	}
}

func initCACert(params *proxyParams)  {
	var caCertPath = params.certPath
	if caCertPath != "" {
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		params.pool.AppendCertsFromPEM(caCrt)
	}
}

func initHttpClient(params *proxyParams) {
	params.pool = x509.NewCertPool()
	initCACert(params)

	// client setting
	var tp = http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: params.pool},
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	params.client = &http.Client{Transport: &tp}
}

func main() {
	var params = proxyParams{}
	initParameter(&params)
	initHttpClient(&params)
	fmt.Println("The request will be transport to: " + params.remote)
	fmt.Println("Listen on " + params.local)
	http.ListenAndServe(params.remote, proxyHandler{params:&params})
	fmt.Println("Listen on local " + params.local)
}
