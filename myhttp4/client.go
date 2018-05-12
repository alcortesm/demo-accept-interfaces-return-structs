package myhttp4

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func NewClient(inner HTTPClient, contentType string) *Client {
	return &Client{
		inner:       inner,
		contentType: contentType,
	}
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
	inner       HTTPClient
	contentType string
}

func (c *Client) Get(url string) (string, error) {
	resp, err := c.inner.Get(url)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

func (c *Client) Post(url string, body io.Reader) (string, error) {
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
