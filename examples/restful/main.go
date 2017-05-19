package main

import (
	"fmt"
	"net/http"
)

func main() {
	h, err := NewHTTPHandler()
	if err != nil {
		fmt.Printf("new http handler: %v\n", err)
		return
	}
	fmt.Printf("serve on :8080\n")
	http.ListenAndServe(":8080", h)
}
