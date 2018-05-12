package myhttp4_test

import (
	"errors"
	"io"
	"io/ioutil"
	"local/demo-accept-interfaces-return-structs/myhttp4"
	"net/http"
	"strings"
	"testing"
)

type mock struct {
	t    *testing.T
	get  func(url string) (*http.Response, error)
	post func(url, contentType string, body io.Reader) (*http.Response, error)
}

func (m *mock) Get(url string) (*http.Response, error) {
	if m.get == nil {
		m.t.Fatal("unexpected call to Get")
	}
	return m.get(url)
}

func (m *mock) Post(a, b string, c io.Reader) (*http.Response, error) {
	if m.post == nil {
		m.t.Fatal("unexpected call to Post")
	}
	return m.post(a, b, c)
}

func TestGet(t *testing.T) {
	t.Parallel()

	urlFix := "url fixture"
	bodyFix := "body fixture"
	errFix := errors.New("error fixture")

	for _, tt := range []struct {
		desc          string
		want          string         // expected return value
		wantError     error          // expected return error
		innerResponse *http.Response // to be returned by the inner client
		innerError    error          // to be returned by the inner client
	}{
		{
			desc:      "200",
			want:      bodyFix,
			wantError: nil,
			innerResponse: &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       ioutil.NopCloser(strings.NewReader(bodyFix)),
			},
			innerError: nil,
		},
		{
			desc:      "404",
			want:      bodyFix,
			wantError: errors.New(http.StatusText(http.StatusNotFound)),
			innerResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     http.StatusText(http.StatusNotFound),
				Body:       ioutil.NopCloser(strings.NewReader(bodyFix)),
			},
			innerError: nil,
		},
		{
			desc:          "error",
			want:          "",
			wantError:     errFix,
			innerResponse: nil,
			innerError:    errFix,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			// mock an GetPoster to return innerResponse and innerError
			// to a GET request for the urlFix
			inner := &mock{
				t: t,
				get: func(url string) (*http.Response, error) {
					if url != urlFix {
						t.Errorf("wrong url, want %q, got %q", urlFix, url)
					}
					return tt.innerResponse, tt.innerError
				},
			}

			sut := myhttp4.NewClient(inner, "irrelevant")
			got, err := sut.Get(urlFix)
			if !equalError(err, tt.wantError) {
				t.Errorf("wrong error: want %q, got %q", tt.wantError, err)
			}
			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func equalError(a, b error) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil && b != nil:
		return false
	case a != nil && b == nil:
		return false
	default:
		return strings.Contains(a.Error(), b.Error())
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	const urlFix = "url fixture"
	const contentTypeFix = "contentType fixture"
	const reqBodyFix = "request body fixture"
	const respBodyFix = "response body fixture"
	errFix := errors.New("error fixture")

	for _, tt := range []struct {
		desc          string
		want          string         // expected return value
		wantError     error          // expected returned error
		innerResponse *http.Response // to be returned by the inner client
		innerError    error          // to be returned by the inner client
	}{
		{
			desc:      "200",
			want:      respBodyFix,
			wantError: nil,
			innerResponse: &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       ioutil.NopCloser(strings.NewReader(respBodyFix)),
			},
			innerError: nil,
		},
		{
			desc:      "404",
			want:      respBodyFix,
			wantError: errors.New(http.StatusText(http.StatusNotFound)),
			innerResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     http.StatusText(http.StatusNotFound),
				Body:       ioutil.NopCloser(strings.NewReader(respBodyFix)),
			},
			innerError: nil,
		},
		{
			desc:          "error",
			want:          "",
			wantError:     errFix,
			innerResponse: nil,
			innerError:    errFix,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			// mock an GetPoster to return innerResponse and innerError
			// to a POST request for the urlFix
			inner := &mock{
				t: t,
				post: func(url, contentType string, body io.Reader) (
					*http.Response, error) {
					if url != urlFix {
						t.Errorf("wrong url, want %q, got %q", urlFix, url)
					}
					if contentType != contentTypeFix {
						t.Errorf("wrong contentType, want %q, got %q",
							contentTypeFix, contentType)
					}
					b, err := ioutil.ReadAll(body)
					if err != nil {
						t.Fatal(err)
					}
					if str := string(b); str != reqBodyFix {
						t.Errorf("wrong body, want %q, got %q",
							reqBodyFix, str)
					}
					return tt.innerResponse, tt.innerError
				},
			}

			sut := myhttp4.NewClient(inner, contentTypeFix)
			got, err := sut.Post(urlFix, strings.NewReader(reqBodyFix))
			if !equalError(err, tt.wantError) {
				t.Errorf("wrong error: want %q, got %q", tt.wantError, err)
			}
			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
