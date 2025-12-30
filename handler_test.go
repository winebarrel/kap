package kap_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/kap"
)

func TestHandlePing(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	kap.HandlePing(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("pong", w.Body.String())
}

func TestAuthHandlerWithHeader(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key:    "my-key",
			Secret: testMakeSecrets(t, []string{"my-secret"}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-key", "my-secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("proxied", w.Body.String())
}

func TestAuthHandlerWitQuery(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key:    "my-key",
			Secret: testMakeSecrets(t, []string{"my-secret"}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	q := r.URL.Query()
	q.Add("my-key", "my-secret")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("proxied", w.Body.String())
}

func TestAuthHandlerMultiKeyWithHeader(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key:    "my-key",
			Secret: testMakeSecrets(t, []string{"my-secret", "my-secret2"}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	tt := []struct {
		key    string
		status int
		body   string
	}{
		{
			key:    "my-secret",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret2",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret3",
			status: http.StatusForbidden,
			body:   "forbidden\n",
		},
	}
	for _, t := range tt {
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		r.Header.Add("my-key", t.key)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(t.status, w.Code)
		assert.Equal(t.body, w.Body.String())
	}
}

func TestAuthHandlerMultiKeyWithQuery(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key:    "my-key",
			Secret: testMakeSecrets(t, []string{"my-secret", "my-secret2"}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	tt := []struct {
		key    string
		status int
		body   string
	}{
		{
			key:    "my-secret",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret2",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret3",
			status: http.StatusForbidden,
			body:   "forbidden\n",
		},
	}
	for _, t := range tt {
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		q := r.URL.Query()
		q.Add("my-key", t.key)
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(t.status, w.Code)
		assert.Equal(t.body, w.Body.String())
	}
}

func TestAuthHandlerForbidden(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key:    "my-key",
			Secret: testMakeSecrets(t, []string{"my-secret"}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	rs := []*http.Request{
		httptest.NewRequest(http.MethodGet, "http://example.com", nil),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			r.Header.Add("my-key", "xmy-secret")
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			r.Header.Add("xmy-key", "my-secret")
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			q := r.URL.Query()
			q.Add("my-key", "xmy-secret")
			r.URL.RawQuery = q.Encode()
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			q := r.URL.Query()
			q.Add("xmy-key", "my-secret")
			r.URL.RawQuery = q.Encode()
			return r
		}(),
	}

	for _, r := range rs {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(http.StatusForbidden, w.Code)
		assert.Equal("forbidden\n", w.Body.String())
	}
}

func TestAuthHandlerMultiHashedKeyWithHeader(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key: "my-key",
			Secret: testMakeSecrets(t, []string{
				"$2y$05$NB57z9o7bUjGSAnaLACGreXbK8lZzyKZFvdKHIsiEMwTgDNzNmeQ.", // my-secret
				"$2y$05$zxOnauk.m1RI3PT/6tKImeIi7WinlJzEI6F.fLIe4Z1fXYU1rLYTm", // my-secret2
			}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	tt := []struct {
		key    string
		status int
		body   string
	}{
		{
			key:    "my-secret",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret2",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret3",
			status: http.StatusForbidden,
			body:   "forbidden\n",
		},
	}
	for _, t := range tt {
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		r.Header.Add("my-key", t.key)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(t.status, w.Code)
		assert.Equal(t.body, w.Body.String())
	}
}

func TestAuthHandlerMultiHashedKeyWithQuery(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key: "my-key",
			Secret: testMakeSecrets(t, []string{
				"$2y$05$WfTzW7zy2I3mba4zszwSm.TjURhBMnAnvzv3upArZ7iD58JaHOU0.", // my-secret
				"$2y$05$T.j4bBTUXb3DWSN1PAvVBu09y99Y0RrcH8uBRif45q1XzkZ4yvcrS", // my-secret2
			}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	tt := []struct {
		key    string
		status int
		body   string
	}{
		{
			key:    "my-secret",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret2",
			status: http.StatusOK,
			body:   "proxied",
		},
		{
			key:    "my-secret3",
			status: http.StatusForbidden,
			body:   "forbidden\n",
		},
	}
	for _, t := range tt {
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		q := r.URL.Query()
		q.Add("my-key", t.key)
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(t.status, w.Code)
		assert.Equal(t.body, w.Body.String())
	}
}

func TestAuthHandlerHashedKeyForbidden(t *testing.T) {
	assert := assert.New(t)

	handler := kap.AuthHandler{
		Options: &kap.Options{
			Key: "my-key",
			Secret: testMakeSecrets(t, []string{
				"$2y$05$00Q8iI4Z0/bOEGmNArXqo.pJ1WaY9umI0K1vdr5p0fxD0QijVWR8C", // my-secret
			}),
		},
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	{
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		r.Header.Add("my-key", "my-secret")
		handler.ServeHTTP(w, r)
		assert.Equal(http.StatusOK, w.Code)
		assert.Equal("proxied", w.Body.String())
	}

	rs := []*http.Request{
		httptest.NewRequest(http.MethodGet, "http://example.com", nil),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			r.Header.Add("my-key", "xmy-secret")
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			r.Header.Add("xmy-key", "my-secret")
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			q := r.URL.Query()
			q.Add("my-key", "xmy-secret")
			r.URL.RawQuery = q.Encode()
			return r
		}(),
		func() *http.Request {
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			q := r.URL.Query()
			q.Add("xmy-key", "my-secret")
			r.URL.RawQuery = q.Encode()
			return r
		}(),
	}

	for _, r := range rs {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		assert.Equal(http.StatusForbidden, w.Code)
		assert.Equal("forbidden\n", w.Body.String())
	}
}
