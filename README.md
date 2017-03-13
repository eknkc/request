# request
--
    import "github.com/eknkc/request"

Package request provides easy to use HTTP request methods

## Usage

#### type Body

```go
type Body interface {
	Mime() string
	Reader() (io.Reader, error)
}
```


#### type Client

```go
type Client interface {
	// Create a request session with custom HTTP method and URL
	Do(method, url string) Session
	// Create a GET request session with specified url
	Get(url string) Session
	// Create a POST request session with specified url
	Post(url string) Session
	// Create a PUT request session with specified url
	Put(url string) Session
	// Create a HEAD request session with specified url
	Head(url string) Session
	// Create a PATCH request session with specified url
	Patch(url string) Session
	// Create a DELETE request session with specified url
	Delete(url string) Session
}
```

Client is the constructor of request Sessions.

#### func  New

```go
func New() Client
```
New creates a new client. Each Client instance has an underlying http.Client
instance so try to create one and reuse it.

#### type Errors

```go
type Errors []error
```

Errors is a collection of errors occured during a request session A session
might generate more than one error during execution, they will be merged into a
single error

#### func (Errors) Error

```go
func (e Errors) Error() string
```

#### type Response

```go
type Response *http.Response
```

Response in an alias of http.Response

#### type Session

```go
type Session interface {
	// End request and return the resulting body as a byte array.
	End() (Response, []byte, error)
	// End request and deserialize JSON response to provided struct.
	EndStruct(target interface{}) (Response, []byte, error)
	// End request and return the result body as string.
	EndString() (Response, string, error)
	// Set request Content-Type
	Type(mime string) Session
	// Set request header
	Header(name, value string) Session
	// Set request querystring parameter
	Query(key string, value string) Session
	// Set request querystring parameters from a supplied struct. Uses `url` tags for field customization.
	QueryStruct(values interface{}) Session
	// Set request body
	Body(r io.Reader) Session
	// Set request body
	BodyBytes(b []byte) Session
	// Serialize and send JSON data
	JSON(data interface{}) Session
	// Set request form parameter
	Form(key string, value string) Session
	// Serialize struct and send as an url encoded form
	FormStruct(values interface{}) Session
	// Set request timeout
	Timeout(duration time.Duration) Session
	// Set request Context
	WithContext(ctx context.Context) Session
}
```

Session is an active HTTP request instance Almost all methods are chainable and
upon constructing the request you need to call EndXXX variants to actually
initiate the request Any error occured before End (such as JSON method producing
an error) would be returned from EndXXX methods
