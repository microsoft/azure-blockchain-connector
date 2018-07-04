package main

import (
	"flag"
	"fmt"
	"io"
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
	//req.URL.Host = "api.github.com"
	req.URL.Host = proxyParameter.host
	req.Host = ""
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "identity")
	client := &http.Client{}
	response, _ := client.Do(req)
	defer response.Body.Close()
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
