package cf

import (
	"encoding/json"
	"io"
	"net/http"
)

// HttpClient used for HTTP Requests. Replace with different http client if needed.
var HttpClient = http.DefaultClient

func sendRequest(req *http.Request, target *Target) (resp *http.Response, err error) {
	req.Header.Set("authorization", "bearer "+target.AccessToken)
	req.Header.Set("accept", "application/json")
	req.Header.Set("User-Agent", "cf90")

	traceReq(req)
	if resp, err = HttpClient.Do(req); err != nil {
		return
	}
	traceResp(resp)

	if resp.StatusCode == http.StatusUnauthorized {
		if seaker, ok := req.Body.(io.Seeker); ok {
			seaker.Seek(0, 0)
		}
		if err = refreshToken(target); err == nil {
			return sendRequest(req, target)
		}
		return
	}

	if resp.StatusCode >= http.StatusBadRequest {
		e := new(Error)
		if json.NewDecoder(resp.Body).Decode(&e) == nil {
			err = e
		}
	}
	return
}
