package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type GURLClient struct {
	C *http.Client
}

func NewGURL(timeout int) *GURLClient {
	// TODO: investivate if this can be achieved from request side with RequestWithContext
	transport := http.Transport{
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
	}
	httpClient := &http.Client{
		Transport: &transport,
	}

	return &GURLClient{
		C: httpClient,
	}
}

func (g *GURLClient) Send(request *http.Request) (*http.Response, error) {
	return g.C.Do(request)
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
