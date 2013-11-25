package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	listing     Listing
	listingSale Listing

	lastNewName, lastNewType string

	Cache = struct {
		All  bytes.Buffer
		Sale bytes.Buffer
	}{}
)

func init() {
	http.HandleFunc("/", HandleOnSale)
	http.HandleFunc("/onsale", HandleOnSale)
	http.HandleFunc("/all", HandleAll)
}

func HandleAll(w http.ResponseWriter, r *http.Request) {
	if Cache.All.Len() == 0 {
		if err := writeList(&Cache.All, &listing); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Cached listing: All")
	}
	w.Header().Add("Last-Modified", listing.LastUpdate.Format(time.RFC1123))
	w.Write(Cache.All.Bytes())
}

func HandleOnSale(w http.ResponseWriter, r *http.Request) {
	if Cache.Sale.Len() == 0 {
		if err := writeList(&Cache.Sale, &listingSale); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Cached listing: Sale")
	}
	w.Header().Add("Last-Modified", listingSale.LastUpdate.Format(time.RFC1123))
	w.Write(Cache.Sale.Bytes())
}

func writeList(w io.Writer, l *Listing) (err error) {
	u := new(Uniques)
	tpl := template.New("index.html")
	tpl.Funcs(template.FuncMap{
		"isNewName": u.Name,
		"isNewType": u.Type,
		"countType": l.CountType,
	})
	if _, err = tpl.ParseFiles("assets/templates/index.html"); err != nil {
		return
	}
	return tpl.Execute(w, &l)
}

type Uniques struct {
	name, typ string
}

func (u *Uniques) Name(s string) bool {
	defer func() { u.name = s }()
	return s != u.name
}

func (u *Uniques) Type(s string) bool {
	defer func() { u.typ = s }()
	return s != u.typ
}
