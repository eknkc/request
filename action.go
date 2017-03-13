package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"strings"

	"github.com/google/go-querystring/query"
)

type action struct {
	client  *client
	ctx     context.Context
	url     *url.URL
	method  string
	errors  []error
	body    Body
	headers http.Header
	timeout time.Duration
}

func (a *action) addError(err error) {
	if err != nil {
		a.errors = append(a.errors, err)
	}
}

func (s *action) End() (Response, []byte, error) {
	if len(s.errors) > 0 {
		return nil, nil, Errors(s.errors)
	}

	var bodyReader io.Reader

	if s.body != nil {
		var err error

		bodyReader, err = s.body.Reader()

		if err != nil {
			return nil, nil, err
		}

		mime := s.body.Mime()

		if mime != "" && (s.headers == nil || s.headers.Get("Content-Type") == "") {
			s.Type(mime)
		}
	}

	req, err := http.NewRequest(s.method, s.url.String(), bodyReader)

	if err != nil {
		return nil, nil, err
	}

	if s.timeout > 0 {
		ctx := req.Context()
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	req.Header = s.headers

	resp, err := s.client.httplclient.Do(req)

	if err != nil {
		return resp, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp, body, err
}

func (s *action) EndStruct(target interface{}) (Response, []byte, error) {
	resp, body, err := s.End()

	if err == nil {
		err = json.Unmarshal(body, target)
	}

	return resp, body, err
}

func (s *action) EndString() (Response, string, error) {
	resp, body, err := s.End()
	return resp, string(body), err
}

func (a *action) WithContext(ctx context.Context) Session {
	a.ctx = ctx
	return a
}

func (a *action) Type(mimeType string) Session {
	return a.Header("Content-Type", mimeType)
}

func (a *action) Header(name, value string) Session {
	if a.headers == nil {
		a.headers = http.Header{}
	}

	a.headers.Set(name, value)
	return a
}

func (a *action) Timeout(duration time.Duration) Session {
	a.timeout = duration
	return a
}

func (a *action) Query(key string, value string) Session {
	oldquery := a.url.Query()
	oldquery.Add(key, value)
	a.url.RawQuery = oldquery.Encode()
	return a
}

func (a *action) QueryStruct(values interface{}) Session {
	vals, err := query.Values(values)
	oldquery := a.url.Query()

	if err != nil {
		a.addError(err)
	} else {
		for key, val := range vals {
			for _, v := range val {
				oldquery.Add(key, v)
			}
		}
	}

	a.url.RawQuery = oldquery.Encode()
	return a
}

func (a *action) Body(r io.Reader) Session {
	a.body = &readerBody{r: r}
	return a
}

func (a *action) BodyBytes(b []byte) Session {
	return a.Body(bytes.NewReader(b))
}

func (a *action) JSON(data interface{}) Session {
	a.body = &jsonBody{value: data}
	return a
}

func (a *action) Form(key, value string) Session {
	if fb, ok := a.body.(*formBody); ok {
		fb.data.Set(key, value)
	}

	fb := &formBody{
		data: url.Values{},
	}

	fb.data.Set(key, value)
	a.body = fb

	return a
}

func (a *action) FormStruct(value interface{}) Session {
	vals, err := query.Values(value)

	if err != nil {
		a.addError(err)
	} else {
		a.body = &formBody{
			data: vals,
		}
	}

	return a
}

type Body interface {
	Mime() string
	Reader() (io.Reader, error)
}

type jsonBody struct {
	value interface{}
}

func (j *jsonBody) Mime() string {
	return "application/json"
}

func (j *jsonBody) Reader() (io.Reader, error) {
	by, err := json.Marshal(j.value)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(by), nil
}

type formBody struct {
	data url.Values
}

func (f *formBody) Mime() string {
	return "application/x-www-form-urlencoded"
}

func (j *formBody) Reader() (io.Reader, error) {
	return strings.NewReader(j.data.Encode()), nil
}

type readerBody struct {
	r io.Reader
}

func (b *readerBody) Mime() string {
	return ""
}

func (b *readerBody) Reader() (io.Reader, error) {
	return b.r, nil
}
