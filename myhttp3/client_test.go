package myhttp3_test

import (
	"fmt"
	"io/ioutil"
	"local/demo-accept-interfaces-return-structs/myhttp3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetOK(t *testing.T) {
	t.Parallel()

	// configure a server to answer 200 and the body fixture for a GET
	// request
	const body = "body fixture"
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("wrong method at server, want GET, got %q", r.Method)
			}
			fmt.Fprint(w, body)
		}))
	defer ts.Close()

	inner := &http.Client{Timeout: time.Second}
	sut := myhttp3.NewClient(inner, "irrelevant")
	got, err := sut.Get(ts.URL)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if got != body {
		t.Errorf("want %q, got %q", body, got)
	}
}

func TestGetResponseNot200(t *testing.T) {
	t.Parallel()

	// configure a server to answer 404
	const body = "body fixture"
	const status = http.StatusNotFound
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			fmt.Fprint(w, body)
		}))
	defer ts.Close()

	inner := &http.Client{Timeout: time.Second}
	sut := myhttp3.NewClient(inner, "irrelevant")
	got, err := sut.Get(ts.URL)
	if err == nil {
		t.Errorf("unexpected success, want 404")
	} else if !strings.Contains(err.Error(), http.StatusText(status)) {
		t.Errorf("wrong error, want 404, got: %v", err)
	}
	if got != body {
		t.Errorf("want %q, got %q", body, got)
	}
}

func TestPostOK(t *testing.T) {
	t.Parallel()

	// configure a server to answer 201 and the body fixture for a POST
	// request
	const requestBody = "request body fixture"
	const responseBody = "response body fixture"
	const contentType = "contentType fixture"
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// check request method
			if r.Method != "POST" {
				t.Errorf("wrong method at server, want POST, got %q", r.Method)
			}
			// check request content type
			{
				contentTypes := map[string][]string(r.Header)["Content-Type"]
				var found bool
				for _, v := range contentTypes {
					if v == contentType {
						found = true
					}
				}
				if !found {
					t.Errorf("didn't find content-type %q inside %q",
						contentType, contentTypes)
				}
			}
			// check request body
			{
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal("cannot read request body")
				}
				got := string(b)
				if got != requestBody {
					t.Errorf("worng request body at server, want %q, got %q",
						requestBody, got)
				}
			}
			// send response
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, responseBody)
		}))
	defer ts.Close()

	inner := &http.Client{Timeout: time.Second}
	sut := myhttp3.NewClient(inner, contentType)
	got, err := sut.Post(ts.URL, strings.NewReader(requestBody))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if got != responseBody {
		t.Errorf("want %q, got %q", responseBody, got)
	}
}

func TestPostResponseNot200(t *testing.T) {
	t.Parallel()

	// configure a server to answer 404
	const body = "body fixture"
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, body)
		}))
	defer ts.Close()

	inner := &http.Client{Timeout: time.Second}
	sut := myhttp3.NewClient(inner, "irrelevant")
	got, err := sut.Post(ts.URL, strings.NewReader("irrelevant"))
	if err == nil {
		t.Errorf("unexpected success, want 404")
	} else if !strings.Contains(err.Error(),
		http.StatusText(http.StatusNotFound)) {
		t.Errorf("wrong error, want 404, got: %v", err)
	}
	if got != body {
		t.Errorf("want %q, got %q", body, got)
	}
}

func TestInnerError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc string
		// call calls a method of c with the given url
		// and returns the error it returns
		call func(c *myhttp3.Client, url string) error
	}{
		{
			"GET",
			func(c *myhttp3.Client, url string) error {
				_, err := c.Get(url)
				return err
			},
		}, {
			"POST",
			func(c *myhttp3.Client, url string) error {
				body := strings.NewReader("irrelevant")
				_, err := c.Post(url, body)
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			const want = "unsupported protocol scheme"
			inner := &http.Client{Timeout: time.Second}
			sut := myhttp3.NewClient(inner, "irrelevant")
			_, err := sut.Get("badscheme://example.com")
			if err == nil {
				t.Errorf("unexpected success, want error %q", want)
			} else if !strings.Contains(err.Error(), want) {
				t.Errorf("wrong error, want %q, got: %v", want, err)
			}
		})
	}
}
