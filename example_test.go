package request_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zorcal/request"
)

// RoundTripperFunc simplifies creating a http.RoundTripper used for
// intercepting the transport of a custom *http.Client in tests.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// Rountrip implements the http.RoundTripper interface.
func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func Example() {
	// HTTP client repeating whatever is in the request body. We override the
	// transport for the purpose of not sending real HTTP requests in this
	// example.
	echolaliaClient := http.Client{
		Timeout: time.Second * 5,
		Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r.Body,
			}, nil
		}),
	}

	// We attach our custom client to ctx for the builder to use.
	ctx := request.AttachClientToContext(context.Background(), &echolaliaClient)

	type payload struct {
		Message string `json:"message"`
	}

	resp, err := request.New().
		WithTimeout(time.Second*10). // Override whatever is set on the HTTP client we passed via ctx.
		WithBasicAuth("username", "password").
		WithJSONBody(&payload{"This is an example."}).
		Do(ctx, http.MethodPost, "http://localhost")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(data))
	// Output:
	// Status: 200
	// Body: {"message":"This is an example."}
}

func Example_withJSONResult() {
	// HTTP client repeating whatever is in the request body. We override the
	// transport for the purpose of not sending real HTTP requests in this
	// example.
	echolaliaClient := http.Client{
		Timeout: time.Second * 5,
		Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r.Body,
			}, nil
		}),
	}

	// We attach our custom client to ctx for the builder to use.
	ctx := request.AttachClientToContext(context.Background(), &echolaliaClient)

	type payload struct {
		Message string `json:"message"`
	}

	var respData payload
	res, err := request.New().
		WithJSONBody(&payload{"This is an example."}).
		WithJSONResult(&respData).
		Do(ctx, http.MethodPost, "http://localhost")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Status: %d\n", res.Response.StatusCode)
	fmt.Printf("Body: %s\n", strings.TrimSpace(string(res.RawData)))
	fmt.Printf("Message: %s\n", respData.Message)
	// Output:
	// Status: 200
	// Body: {"message":"This is an example."}
	// Message: This is an example.
}
