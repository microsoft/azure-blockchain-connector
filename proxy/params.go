package proxy

const (
	MethodBasicAuth       = ""
	MethodOAuthAuthCode   = "authcode"
	MethodOAuthDeviceFlow = "device"

	LogWhenOnError  = "onError"  // print log only for those who raise exceptions
	LogWhenOnNon200 = "onNon200" // print log for those who have a non-200 response, or those who raise exceptions
	LogWhenAlways   = "always"   // print log for every request

	LogWhatBasic    = "basic"    // print the request's method and URI and the response status code (and the exception message, if exception raised) in the log
	LogWhatDetailed = "detailed" // print the request's method, URI and body, and the response status code and body (and the exception message, if exception raised) in the log
	//LogAll          = "all"      // to be supported later. Compared to whatlog_detail, all Headers are printed in whatlog_all
)

type Params struct {
	Local  string
	Remote string
	Method string

	Whenlog string
	Whatlog string
}
