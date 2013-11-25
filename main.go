package main

import (
	"log"
	"net/http"
)

func main() {
	if err := initListing("data", &listing); err != nil {
		log.Fatal(err)
	}
	initListingSale(&listingSale, &listing)

	log.Fatal(http.ListenAndServe(":"+env("PORT", "8081"), nil))
}
