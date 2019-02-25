package proxy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/vicanso/cod"
)

func TestProxy(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		target, _ := url.Parse("https://www.baidu.com")
		config := Config{
			Target:    target,
			Host:      "www.baidu.com",
			Transport: &http.Transport{},
			Rewrites: []string{
				"/api/*:/$1",
			},
		}
		fn := New(config)
		req := httptest.NewRequest("GET", "http://127.0.0.1/api/", nil)
		originalPath := req.URL.Path
		originalHost := req.Host
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn(c)
		if c.Request.URL.Path != originalPath {
			t.Fatalf("request path should be reverted")
		}
		if req.Host != originalHost {
			t.Fatalf("request host should be reverted")
		}

		if !done || c.StatusCode != http.StatusOK {
			t.Fatalf("http proxy fail")
		}
	})

	t.Run("target picker", func(t *testing.T) {
		target, _ := url.Parse("https://www.baidu.com")
		config := Config{
			TargetPicker: func(c *cod.Context) (*url.URL, error) {
				return target, nil
			},
			Host:      "www.baidu.com",
			Transport: &http.Transport{},
		}
		fn := New(config)
		req := httptest.NewRequest("GET", "http://127.0.0.1/", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn(c)
		if !done || c.StatusCode != http.StatusOK {
			t.Fatalf("http proxy fail")
		}
	})

	t.Run("target picker error", func(t *testing.T) {
		config := Config{
			TargetPicker: func(c *cod.Context) (*url.URL, error) {
				return nil, errors.New("abcd")
			},
			Host:      "www.baidu.com",
			Transport: &http.Transport{},
		}
		fn := New(config)
		req := httptest.NewRequest("GET", "http://127.0.0.1/", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		err := fn(c)
		if err.Error() != "abcd" {
			t.Fatalf("proxy should return error")
		}
	})

	t.Run("no target", func(t *testing.T) {
		config := Config{
			TargetPicker: func(c *cod.Context) (*url.URL, error) {
				return nil, nil
			},
			Host:      "www.baidu.com",
			Transport: &http.Transport{},
		}
		fn := New(config)
		req := httptest.NewRequest("GET", "http://127.0.0.1/", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		err := fn(c)
		if err.Error() != "category=cod-proxy, message=target can not be nil" {
			t.Fatalf("nil proxy should return error")
		}
	})
}
