package request

import (
	"bytes"
	"fmt"
	"github.com/smf8/gurl/internal/pkg/validation"
	"io"
	"log"
	"net/http"
	"net/url"
)

// GURL represents a gurl command. it's initialized from cli arguments and options.
type GURL struct {
	URL           string
	Method        string
	Headers       map[string]string
	QueryParams   map[string][]string
	Data          string
	FilePath      string
	JSONMessage   string
	ClientTimeout int
}

func (g *GURL) ToHTTPRequest() (*http.Request, error) {
	requestURL, err := url.ParseRequestURI(g.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid url format: %w", err)
	}

	// here we only "append" to existing URL parameters.
	urlParams := requestURL.Query()
	for key, values := range g.QueryParams {
		for _, value := range values {
			urlParams.Add(key, value)
		}
	}

	requestURL.RawQuery = urlParams.Encode()

	// validate and set request body
	dataReader, err := g.parseData()
	if err != nil {
		return nil, fmt.Errorf("invalid request data: %w", err)
	}

	req, err := http.NewRequest(g.Method, requestURL.String(), dataReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// setting request headers.
	for key, value := range g.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (g *GURL) parseData() (io.Reader, error) {
	if g.JSONMessage != "" {
		if _, ok := g.Headers[HeaderContentType]; !ok {
			g.Headers[HeaderContentType] = ContentTypeJSON
		}

		if !validation.IsStringJSON(g.JSONMessage) {
			return nil, fmt.Errorf("invalid JSON data")
		}

		return bytes.NewReader([]byte(g.JSONMessage)), nil
	} else if g.FilePath != "" {
		if _, ok := g.Headers[HeaderContentType]; !ok {
			g.Headers[HeaderContentType] = ContentTypeOctetStream
		}

		filePayload := &FilePayload{
			FilePath: g.FilePath,
		}

		defer func() {
			if err := filePayload.CloseFile(); err != nil {
				log.Printf("failed to close file: %s", err.Error())
			}
		}()

		fileData, err := io.ReadAll(filePayload.Data())
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(fileData), nil
	} else {
		g.Headers[HeaderContentType] = ContentTypeFormURLEncoded

		return bytes.NewReader([]byte(g.Data)), nil
	}
}
