package main

import (
	"fmt"
	"net/http"
)

func init() {
	for _, dir := range []string{"js"} {
		http.Handle(fmt.Sprintf("/%s/", dir), http.FileServer(http.Dir("assets")))
	}
}
