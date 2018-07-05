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

type httpProxyParameter struct {
	port string
	host string
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//req.URL.
	req.URL.Scheme = "https"
	// req.URL.Host = "api.github.com"
	// req.URL.Host = "127.0.0.1:1235"
	// req.URL.Host = "104.215.148.235:1234"
	req.URL.Host = proxyParameter.host
	req.Host = ""
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "identity")

	// ca cert setting
	pool := x509.NewCertPool()
	caCertPath := "ca.crt"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	// client setting
	tp := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		// 	// skip verify
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tp}

	// do request and get response
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	io.Copy(rw, response.Body)
}

func initParameter() {

	flag.StringVar(&proxyParameter.port, "port", "1234", "The port you want to listen")
	flag.StringVar(&proxyParameter.host, "host", "api.github.com", "The host you want to send to")
	flag.Parse()
	fmt.Print("The aim port is: ")
	fmt.Println(proxyParameter.port)
	fmt.Print("The aim host is: ")
	fmt.Println(proxyParameter.host)
}

func main() {
	/*proxyParameter := httpProxyParameter{
		//port: "1234",
		//host: "api.github.com",
	}*/

	initParameter()
	//fmt.Print(proxyParameter)
	http.HandleFunc("/", ServeHTTP)
	fmt.Println("Listen on 127.0.0.1:" + proxyParameter.port)
	fmt.Println("The request will be transport to: " + proxyParameter.host)
	http.ListenAndServe("127.0.0.1:"+proxyParameter.port, nil)
	fmt.Println("Listen on port " + proxyParameter.port)
}
