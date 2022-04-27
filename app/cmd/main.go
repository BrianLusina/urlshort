package main

import (
	"flag"
	"net/http"
	"urlshort/app/api"
	"urlshort/app/internal/repo"
)

var (
	listenAddr = flag.String("http", ":3000", "http listen address")
	dataFile   = flag.String("file", "store.json", "data store file name")
	hostname   = flag.String("host", "localhost", "host name and port")
)

func main() {
	flag.Parse()
	store := repo.NewUrlStore(*dataFile)
	handler := api.NewHandler(store)

	http.HandleFunc("/", handler.Redirect)
	http.HandleFunc("/add", handler.Add)
	http.ListenAndServe(*listenAddr, nil)
}
