package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Result represents the result from sending a HTTP request and reading the
// response body. Holds the HTTP response and the raw data from reading the
// response body.
type Result struct {
	// The HTTP response from sending a HTTP request with response body read to
	// completion and closed. Attempting to read from the body will result in an
	// error.
	Response *http.Response
	// Raw data from reading all of the response body.
	RawData []byte
}

// WithResult allows for returning a Result after sending a HTTP request.
// Must not be used directly, but instead by calling a With* function on
// Request.
type WithResult struct {
	req       *Request
	unmarshal func(data []byte) error
}

// Do sends an HTTP request and returns a Result containing a HTTP response
// and its raw data from reading response body. Closes the response body.
func (wr *WithResult) Do(ctx context.Context, method, url string) (*Result, error) {
	if wr.req == nil {
		return nil, fmt.Errorf("request: missing request")
	}

	resp, err := wr.req.Do(ctx, method, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("request: read response body: %w", err)
	}

	if wr.unmarshal != nil {
		if err := wr.unmarshal(data); err != nil {
			return nil, err
		}
	}

	return &Result{
		Response: resp,
		RawData:  data,
	}, nil
}
