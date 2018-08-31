package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
)

const defaultLocal = "127.0.0.1:3100"

type proxyParams struct {
	local, remote                string
	username, password, certPath string
	insecure                     bool
	pool                         *x509.CertPool
	client                       *http.Client
}

type proxyHandler struct {
	params *proxyParams
}

func (handler proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var params = handler.params
	//req.URL.Scheme = "https"
	//req.URL.Host = params.remote
	//req.Host = ""
	//req.RequestURI = ""
	//req.Header.Set("Accept-Encoding", "identity")
	//req.SetBasicAuth(params.username, params.password)

	//fmt.Println("Requesting:", req.Method, req.URL)
	//logStr := "Requesting: " + req.Method + " " + req.URL.String() + "\n"
	logStr := fmt.Sprintln("Requesting:", req.Method, req.URL)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	//fmt.Println(buf.String())
	logStr += buf.String() + "\n"

	req1, err := http.NewRequest(req.Method, "https://"+params.remote, buf)
	if err != nil {
		fmt.Println("Error when make transport request:")
		fmt.Println(err)
	}
	req1.ContentLength = req.ContentLength
	req1.Header = req.Header
	req1.Method = req.Method
	req1.SetBasicAuth(params.username, params.password)

	// do request and get response
	response, err := params.client.Do(req1)
	if err != nil {
		fmt.Println("error:", err)
		//logStr += "Error:" + err.Error() + "\n\n"
		logStr = fmt.Sprintln(logStr, "Error:\n", err, "\n")
		fmt.Println(logStr)
		return
	}
	defer response.Body.Close()

	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	logStr = fmt.Sprintln(logStr, "Response Status Code:", response.StatusCode)
	logStr = fmt.Sprintln(logStr, buf.String(), "\n")
	fmt.Println(logStr)
	rw.WriteHeader(response.StatusCode)
	//io.Copy(rw, response.Body)
	rw.Write(buf.Bytes())
}

func initParameter(params *proxyParams) {
	flag.StringVar(&params.local, "local", "", "Local address to bind to")
	flag.StringVar(&params.remote, "remote", "", "Remote endpoint address")
	flag.StringVar(&params.username, "username", "", "The username you want to login with")
	flag.StringVar(&params.password, "password", "", "The password you want to login with")
	flag.StringVar(&params.certPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.insecure, "insecure", false, "Skip certificate verifications")
	flag.Parse()

	if params.local == "" {
		params.local = defaultLocal
	}

	if params.remote == "" || params.username == "" || params.password == "" {
		flag.Usage()
		os.Exit(-1)
	}
}

func initCACert(params *proxyParams) {
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
		TLSClientConfig: &tls.Config{
			RootCAs:            params.pool,
			InsecureSkipVerify: params.insecure,
		},
		MaxIdleConnsPerHost: 1024,
	}
	params.client = &http.Client{Transport: &tp}
}

func testConnection(params *proxyParams) {
	req, err := http.NewRequest(http.MethodGet, "https://"+params.remote, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.SetBasicAuth(params.username, params.password)
	res, err := params.client.Do(req)
	if err != nil {
		fmt.Println("Error occurred when sending test request to the remote host:")
		fmt.Println(err)
		fmt.Println("Please check your Internet connection and the remote host address.")
		os.Exit(-2)
	}
	if res.StatusCode == 401 {
		fmt.Println("Unable to pass the authentication on the remote server. Please Check your username and password.")
		os.Exit(-2)
	}

}

func main() {
	var cancellation = make(chan os.Signal)
	var params = proxyParams{}
	initParameter(&params)
	initHttpClient(&params)
	testConnection(&params)
	fmt.Println("The request will be transport to: " + params.remote)
	fmt.Println("Listen on " + params.local)
	if err := http.ListenAndServe(params.local, proxyHandler{params: &params}); err != nil {
		fmt.Println("Error on listening: ", err)
		os.Exit(-2)
	} else {
		fmt.Println("Connector started")
	}

	signal.Notify(cancellation, os.Interrupt)
	<-cancellation
	fmt.Println("Cancelling...")
}
