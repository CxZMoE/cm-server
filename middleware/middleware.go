package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("[%s] Remote: %s Request: %s", r.Method, r.RemoteAddr, r.URL.String()))
}

func Calc(a, b int) {
	log.Println(a + b)
}
