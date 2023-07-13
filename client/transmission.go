package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Method    string `json:"method"`
	Arguments struct {
		Fields map[string]string `json:"fields"`
	} `json:"arguments"`
}

func Version() {

	request := Request{
		Method: "session-get",
	}

	request.Arguments.Fields = map[string]string{
		"version": "",
	}

	body, _ := json.Marshal(request)

	requestHttp, _ := http.NewRequest(http.MethodPost, "http://localhost:9091/transmission/rpc", bytes.NewReader(body))
	response, _ := http.DefaultClient.Do(requestHttp)

	if response.StatusCode == 409 {
		println("here")
		body, _ := json.Marshal(request)
		requestHttp, _ = http.NewRequest(http.MethodPost, "http://localhost:9091/transmission/rpc", bytes.NewReader(body))
		requestHttp.Header.Set("X-Transmission-Session-Id", response.Header.Get("X-Transmission-Session-Id"))
	}

	response, err := http.DefaultClient.Do(requestHttp)
	if err != nil {
		panic(err)
	}
	responseBody, _ := io.ReadAll(response.Body)

	var buffer bytes.Buffer
	json.Indent(&buffer, responseBody, " ", " ")
	fmt.Printf("%s", buffer.Bytes())
}
