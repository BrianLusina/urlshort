package main

import (
	"net/http"
	"urlshort/app/api"
	"urlshort/app/internal/repo"
)

func main() {
	store := repo.NewUrlStore()
	handler := api.NewHandler(store)

	http.HandleFunc("/", handler.Redirect)
	http.HandleFunc("/add", handler.Add)
	http.ListenAndServe(":3000", nil)
}
