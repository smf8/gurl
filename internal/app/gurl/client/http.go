package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type GURLClient struct {
	C *http.Client
}

func NewGURL() *GURLClient {
	httpClient := &http.Client{}

	return &GURLClient{
		C: httpClient,
	}
}

func (g *GURLClient) Send(request *http.Request, timeout time.Duration) (*http.Response, error) {
	var resp *http.Response

	// we use custom cancellation to detect if response is being filled or not
	c, cancel := context.WithCancel(context.Background())
	req := request.WithContext(c)

	go func() {
		<-time.After(timeout)
		if resp == nil {
			cancel()
		}
	}()

	resp, err := g.C.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GURLClient) DumpResponse(response *http.Response, shouldPrint bool) map[string]interface{} {
	result := make(map[string]interface{})

	result["headers"] = response.Header
	result["status"] = response.Status
	result["code"] = response.StatusCode
	result["Protocol"] = response.Proto

	// handle response body
	// TODO: do something else about files with non-text content

	if body, err := ioutil.ReadAll(response.Body); err == nil {
		result["Body"] = string(body)
	} else {
		log.Printf("failed to read response body: %s", err.Error())
	}

	if shouldPrint {
		for key, value := range result {
			fmt.Printf("[%s] : %v\n===========================\n", key, value)
		}
	}
	return result
}
