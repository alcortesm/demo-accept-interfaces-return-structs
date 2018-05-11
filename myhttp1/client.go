package myhttp1

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Get(url string) (string, error)
	Post(url string, body io.Reader) (string, error)
}

func NewClient(timeout time.Duration, contentType string) Client {
	return &client{
		inner:       &http.Client{Timeout: timeout},
		contentType: contentType,
	}
}

type client struct {
	inner       *http.Client
	contentType string
}

func (c *client) Get(url string) (string, error) {
	resp, err := c.inner.Get(url)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

func (c *client) Post(url string, body io.Reader) (string, error) {
	resp, err := c.inner.Post(url, c.contentType, body)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

func readBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New(resp.Status)
	}
	return string(b), err
}
