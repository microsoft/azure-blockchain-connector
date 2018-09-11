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
	"strings"
)

const defaultLocal = "127.0.0.1:3100"

type proxyParams struct {
	local, remote                string
	username, password, certPath string
	insecure                     bool
	pool                         *x509.CertPool
	client                       *http.Client
	whenlog                      int
	whatlog                      int
}

type proxyHandler struct {
	params *proxyParams
}

const (
	whenlog_onError  int = iota //print log only for those who raise exceptions
	whenlog_onNon200            //print log for those who have a non-200 response, or those who raise exceptions
	whenlog_always              //print log for every request
)

const (
	whatlog_basic    int = iota //print the request's method and URI and the response status code (and the exception message, if exception raised) in the log
	whatlog_detailed            //print the request's method, URI and body, and the response status code and body (and the exception message, if exception raised) in the log
	//whatlog_all				//to be supported later. Compared to whatlog_detail, all Headers are printed in whatlog_all
)

func (handler proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var params = handler.params

	/**logFlag is initialized with value true.
	 * it will be set false if our program finally ensure it's not needed to print the log (depends on the running state and params.whenlog).
	 * when ServeHTTP finished (or crashed), if logFlag remains true, log will be printed
	 */
	logFlag := new(bool)
	*logFlag = true

	logStrBuilder := new(strings.Builder)

	/**Notice that here the func in defer is needed!
	 * By doing so, defer will register the pointer strBuilder and flag, and we can change what the pointers point to later.
	 * Without the func, what defer registers is not the pointers, and defer will know nothing about the later changes to stringbulider and flag.
	 */
	defer func(strBuilder *strings.Builder, flag *bool) {
		if *flag {
			fmt.Println(strBuilder.String())
		}
	}(logStrBuilder, logFlag)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	req.URL.Host = params.remote
	req.URL.Scheme = "https"

	logStrBuilder.WriteString(fmt.Sprintln("Requesting:", req.Method, req.URL))
	if params.whatlog >= whatlog_detailed {
		logStrBuilder.WriteString(buf.String() + "\n")
	}

	req1, err := http.NewRequest(req.Method, req.URL.String(), buf)
	if err != nil {
		logStrBuilder.WriteString(fmt.Sprintln("Error when make transport request:\n", err))
		return
	}
	req1.ContentLength = req.ContentLength
	req1.Header = req.Header
	req1.Method = req.Method
	req1.SetBasicAuth(params.username, params.password)

	// do request and get response
	response, err := params.client.Do(req1)
	if err != nil {
		logStrBuilder.WriteString(fmt.Sprintln("Error when send the transport request:\n", err))
		return
	}
	defer response.Body.Close()

	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	logStrBuilder.WriteString(fmt.Sprintln("Response Status Code:", response.StatusCode))
	if params.whatlog >= whatlog_detailed {
		logStrBuilder.WriteString(fmt.Sprintln(buf.String()))
	}

	rw.WriteHeader(response.StatusCode)
	rw.Write(buf.Bytes())

	if params.whenlog == whenlog_onError {
		*logFlag = false
	}
	if params.whenlog == whenlog_onNon200 {
		if response.StatusCode == 200 {
			*logFlag = false
		}
	}
}

func initParameter(params *proxyParams) {
	var whenlogstr string
	var whatlogstr string
	var debugmode bool
	flag.StringVar(&params.local, "local", "", "Local address to bind to")
	flag.StringVar(&params.remote, "remote", "", "Remote endpoint address")
	flag.StringVar(&params.username, "username", "", "The username you want to login with")
	flag.StringVar(&params.password, "password", "", "The password you want to login with")
	flag.StringVar(&params.certPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.insecure, "insecure", false, "Skip certificate verifications")

	flag.StringVar(&whenlogstr, "whenlog", "onError", "Configuration about in what cases logs should be prited. Alternatives: always, onNon200 and onError. Default: onError")
	flag.StringVar(&whatlogstr, "whatlog", "basic", "Configuration about what information should be included in logs. Alternatives: basic and detailed. Default: basic")
	flag.BoolVar(&debugmode, "debugmode", false, "Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.")

	flag.Parse()

	if params.local == "" {
		params.local = defaultLocal
	}

	if params.remote == "" || params.username == "" || params.password == "" {
		flag.Usage()
		os.Exit(-1)
	}

	switch whenlogstr {
	case "onError":
		params.whenlog = whenlog_onError
	case "onNon200":
		params.whenlog = whenlog_onNon200
	case "always":
		params.whenlog = whenlog_always
	default:
		fmt.Println("Unexpected whenlog value. Expected: always, onNon200 or onError")
		os.Exit(-1)
	}

	switch whatlogstr {
	case "basic":
		params.whatlog = whatlog_basic
	case "detailed":
		params.whatlog = whatlog_detailed
	default:
		fmt.Println("Unexpected whatlog value. Expected: basic or detailed")
		os.Exit(-1)
	}

	if debugmode {
		params.whenlog = whenlog_always
		params.whatlog = whatlog_detailed
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
