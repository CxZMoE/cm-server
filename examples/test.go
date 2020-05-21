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
		return c.String("test", 404)
	}, middleware.Log)

	server.ListenAndServe("127.0.0.1:3030")
}
