package cf

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
)

var (
	// Enables Logging
	Trace bool = false
	// Writer used for tracing (osdefault .Stderr)
	Out io.Writer = os.Stderr
)

func traceResp(resp *http.Response) {
	if Trace {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Fprintln(Out, sanitize(string(dump)), "\n")
	}
}
func traceReq(req *http.Request) {
	if Trace {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Fprintln(Out, sanitize(string(dump)), "\n")
	}
}

const hidden = "*"

var (
	reAuth *regexp.Regexp = regexp.MustCompile(`(?m)^Authorization: .*`)
	rePass *regexp.Regexp = regexp.MustCompile(`password=[^&]*&`)
	reAT   *regexp.Regexp = regexp.MustCompile(`"access_token":"[^"]*"`)
	reRT   *regexp.Regexp = regexp.MustCompile(`"refresh_token":"[^"]*"`)
)

func sanitize(input string) (sanitized string) {
	sanitized = reAuth.ReplaceAllString(input, "Authorization: *")
	sanitized = rePass.ReplaceAllString(sanitized, "password=*&")
	sanitized = reAT.ReplaceAllString(sanitized, `"access_token":*`)
	sanitized = reRT.ReplaceAllString(sanitized, `"refresh_token":*`)
	return
}
