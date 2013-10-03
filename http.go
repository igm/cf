package cf

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// HttpClient used for HTTP Requests. Replace with different http client if needed.
var HttpClient = http.DefaultClient

type closeReader struct{ *bytes.Reader }

func (c closeReader) Close() error { return nil }

func (target *Target) sendRequest(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("authorization", "bearer "+target.AccessToken)
	req.Header.Set("accept", "application/json")
	req.Header.Set("User-Agent", "cf90")
	if r, ok := req.Body.(closeReader); ok {
		req.ContentLength = int64(r.Len())
	}

	traceReq(req)
	if resp, err = HttpClient.Do(req); err != nil {
		return
	}
	traceResp(resp)

	if resp.StatusCode == http.StatusUnauthorized {
		// WARNING!!! io.Seeker.Seek does not work if tracing is enabled, body is replaced in trace
		if !Trace || req.Method == "GET" {
			if seaker, ok := req.Body.(io.Seeker); ok {
				seaker.Seek(0, 0)
			}
			if err = target.refreshToken(); err == nil {
				return target.sendRequest(req)
			}
			return
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		e := new(Error)
		if json.NewDecoder(resp.Body).Decode(&e) == nil {
			err = e
		}
	}
	return
}
