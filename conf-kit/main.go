package main

import (
	"fmt"
	"net/http"

	"github.com/ironzhang/matrix/conf-kit/server"
)

func main() {
	fmt.Printf("conf-kit start...\n")
	http.ListenAndServe(":7200", &server.Handler{})
}
