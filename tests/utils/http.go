package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ParseResponse(v interface{}, resp *http.Response) error {
	body, _ := ioutil.ReadAll(resp.Body)
	return json.Unmarshal(body, v)
}

// PostJsonRequest creates a new POST request with JSON body.
//
// The value can be nil, which means the body is empty.
func PostJsonRequest(uri string, v interface{}) *http.Request {
	var body io.Reader
	if v != nil {
		b, _ := json.Marshal(v)
		body = bytes.NewReader(b)
	}
	req := httptest.NewRequest(http.MethodPost, uri, body)
	req.Header.Add(`Content-Type`, `application/json`)
	return req
}
