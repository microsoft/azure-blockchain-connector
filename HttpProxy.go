package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var proxyParameter httpProxyParameter
var pool *x509.CertPool
var tp *http.Transport
var client *http.Client

type httpProxyParameter struct {
	port string
	host string
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//change url
	req.URL.Scheme = "https"
	// req.URL.Host = "104.215.148.235:1234"
	req.URL.Host = proxyParameter.host
	req.Host = ""
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "identity")

	// do request and get response
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer response.Body.Close()
	io.Copy(rw, response.Body)
}

func initParameter() {

	flag.StringVar(&proxyParameter.port, "port", "1234", "The port you want to listen")
	flag.StringVar(&proxyParameter.host, "host", "104.215.148.235:1234", "The host you want to send to")
	flag.Parse()
	fmt.Print("The aim port is: ")
	fmt.Println(proxyParameter.port)
	fmt.Print("The aim host is: ")
	fmt.Println(proxyParameter.host)
}

func initClient() {
	// ca cert setting
	pool = x509.NewCertPool()
	caCertPath := "ca.crt"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	// client setting
	tp = &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tp}
}

func main() {
	/*proxyParameter := httpProxyParameter{
		//port: "1234",
		//host: "api.github.com",
	}*/

	initParameter()
	initClient()
	//fmt.Print(proxyParameter)
	http.HandleFunc("/", ServeHTTP)
	fmt.Println("Listen on 127.0.0.1:" + proxyParameter.port)
	fmt.Println("The request will be transport to: " + proxyParameter.host)
	http.ListenAndServe("127.0.0.1:"+proxyParameter.port, nil)
	fmt.Println("Listen on port " + proxyParameter.port)
}
