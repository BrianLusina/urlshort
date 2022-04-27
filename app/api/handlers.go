package api

import (
	"fmt"
	"net/http"
	"urlshort/app/internal/repo"
)

const addForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
</html></body>
`

type Handler struct {
	store *repo.UrlStore
}

func NewHandler(store *repo.UrlStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (hdl *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url := hdl.store.Get(key)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (hdl *Handler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, addForm)
		return
	}
	key := hdl.store.Put(url)

	fmt.Fprintf(w, "%s", key)
}
