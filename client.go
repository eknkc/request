package request

import (
	"context"
	"net/http"
	"net/url"
)

type client struct {
	httplclient *http.Client
}

func (c *client) Do(method, requestUrl string) Session {
	u, err := url.Parse(requestUrl)

	act := &action{
		client: c,
		ctx:    context.Background(),
		url:    u,
		method: method,
	}

	act.addError(err)

	return act
}

func (c *client) Get(url string) Session {
	return c.Do("GET", url)
}

func (c *client) Post(url string) Session {
	return c.Do("POST", url)
}

func (c *client) Put(url string) Session {
	return c.Do("PUT", url)
}

func (c *client) Head(url string) Session {
	return c.Do("HEAD", url)
}

func (c *client) Patch(url string) Session {
	return c.Do("PATCH", url)
}

func (c *client) Delete(url string) Session {
	return c.Do("DELETE", url)
}
