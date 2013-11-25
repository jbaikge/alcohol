package main

import (
	"log"
	"net/http"
)

func main() {
	if err := initListing("data", &listing); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":"+env("PORT", "8081"), nil))
}
