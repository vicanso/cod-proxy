# cod-proxy

[![Build Status](https://img.shields.io/travis/vicanso/cod-proxy.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-proxy)

Proxy middleware for cod, it can proxy http request to other host.

```go
package main

import (
	"net/url"

	"github.com/vicanso/cod"

	proxy "github.com/vicanso/cod-proxy"
)

func main() {
	d := cod.New()

	target, _ := url.Parse("https://www.baidu.com")

	d.GET("/*url", proxy.New(proxy.Config{
		Target: target,
		Host:   "www.baidu.com",
	}))

	d.ListenAndServe(":7001")
}
```