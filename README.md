# request

Go package that provides syntactic sugar and sane defaults for sending HTTP requests.

All of the documentation can be found on the [go.dev](https://pkg.go.dev/github.com/zorcal/request?tab=doc) website.

## Usage

Performing a simple POST request (omitting error handling for brevity):

```go
type NewMessage struct {
	Text string `json:"text"`
}

type Message struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

req := request.New().
	WithTimeout(time.Second*10).
	WithBasicAuth("username", "password").
	WithJSONBody(&NewMessage{"Hello world!"}).
	WithAccept("application/json"). // This is not necessary when using WithJSONBody
	WithHeader("x-trace-id", "477cd6fa-758c-4f85-97b9-6f180a703039")

_ := req.Do(context.Background(), http.MethodPost, "http://localhost/api/v1/messages")
defer resp.Body.Close()

b, _ := io.ReadAll(resp.Body)

var msg Message
_ := json.Unmarshal(b, &msg)
```

Let `request` handle the unmarshalling:

```go
var msg Message
req := request.New().
	// ... options in previous example omitted for brevity
	WithJSONResult(&msg)

result, _ := Do(context.Background(), http.MethodPost, "http://localhost/api/v1/messages")

// Access the raw data from from reading the response body at result.RawData.
// Acces the *http.Respnse at result.Response. Reading from result.Response.Body
// results in an error as the body has already been read. No need to close the 
// response body.
```

Overriding the underlying HTTP client via context:



```go
// We can set the underlying *http.Client via context.Context passed to
// the Do function:
//
// 	ctx := request.AttachClientToContext(context.Background(), &http.Client{})
//
// This is great for testing as we can override the `Transport` field on the 
// client and thus mock the external system without needing to create an 
// interface.

// RoundTripperFunc simplifies creating a http.RoundTripper, which is used to
// intercept the transport of custom *http.Client in tests.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// Rountrip implements http.RoundTripper.
func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

mockedClient := http.Client{
	Timeout: time.Second * 5,
	Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		body := strings.NewReader(`{"id": 123, "text": "Test message."}`)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body)}, nil
	}),
}

// We attach our custom client to ctx for the request builder to use.
ctx := request.AttachClientToContext(context.Background(), &mockedClient)

_ := request.New().Do(ctx, http.MethodGry, "http://localhost/api/v1/messages/123")
```