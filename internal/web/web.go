package web

import (
	"net/http"
	"os"
	"path/filepath"

	"ecksbee.com/kushim/internal/cache"
	"github.com/gorilla/mux"
)

func Catalog() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		data, err := cache.MarshalCatalog()
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func ConceptCard() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		href := r.URL.Query().Get("href")

		if len(href) <= 0 {
			http.Error(w, "Error: invalid href", http.StatusBadRequest)
			return
		}
		data, err := cache.MarshalConceptCard(href)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func Renderable() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		hash := vars["hash"]
		if len(hash) <= 0 {
			http.Error(w, "Error: invalid hash", http.StatusBadRequest)
			return
		}

		data, err := cache.MarshalRenderable(hash)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.Path("/concepts").HandlerFunc(ConceptCard()).Methods("GET")
	r.Path("/packages/{hash}").HandlerFunc(Renderable()).Methods("GET")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	taxonomypackagebrowser := http.FileServer(http.Dir((filepath.Join(wd, "heroicking-atrahasis"))))
	r.PathPrefix("/").Handler(http.StripPrefix("/", taxonomypackagebrowser))
	return r
}
