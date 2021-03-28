package client

import (
	"context"
	"fmt"
	"github.com/smf8/gurl/internal/app/gurl/request"
	"github.com/smf8/gurl/internal/app/gurl/response"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func (g *GURLClient) DumpResponse(resp *http.Response, shouldPrint bool) map[string]interface{} {
	result := make(map[string]interface{})

	result["headers"] = resp.Header
	result["status"] = resp.Status
	result["code"] = resp.StatusCode
	result["Protocol"] = resp.Proto

	// handle resp body

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read resp body: %s", err.Error())
	}

	contentType := resp.Header.Get(request.HeaderContentType)

	if !strings.HasPrefix(contentType, "text") {
		err = response.SaveFile(contentType, body, resp.Request.URL.Path)
		if err != nil {
			log.Printf("failed to save file: %s", err.Error())
		}
	} else {
		result["Body"] = string(body)
	}

	if shouldPrint {
		for key, value := range result {
			fmt.Printf("[%s] : %v\n===========================\n", key, value)
		}
	}

	return result
}
