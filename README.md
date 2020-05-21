# cm-server
A light weight and easy to use http server framework. 一个轻量且简单易用的HTTP服务器框架

## Installation

``` shell
go get github.com/CxZMoE/cm-server
```

## Usage

``` go
package main

import (
	"os"

	"github.com/CxZMoE/cm-server/cmserver"
	"github.com/CxZMoE/cm-server/middleware"
)

func main() {
	server := cmserver.NewServer()
	server.Get("/music", func(c *cmserver.Context) error {
		f, err := os.Open("Answer.mp3")
		if err != nil {
			return err
		}
		defer f.Close()
		return c.File(f, 200)
	}, middleware.Log)

	server.Get("/test", func(c *cmserver.Context) error {
		return c.String("test", 200)
	}, middleware.Log)

	server.ListenAndServe("127.0.0.1:3030")
}
```
