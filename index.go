package main

import (
	"html/template"
	"net/http"
)

var (
	listing Listing
)

func init() {
	http.HandleFunc("/", HandleIndex)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tpl.Execute(w, &listing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
