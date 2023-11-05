package web

import (
	"net/http"
	neturl "net/url"

	"ecksbee.com/kushim/internal/cache"
	"github.com/gorilla/mux"
)

func Namespaces() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		parsedquery, err := neturl.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ns, err := neturl.QueryUnescape(parsedquery.Get("ns"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := cache.MarshalNamespace(ns)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func RelationshipSets() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		parsedquery, err := neturl.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		roleuri, err := neturl.QueryUnescape(parsedquery.Get("roleuri"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := cache.MarshalRelationshipSet(roleuri)
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
	r.Path("/namespaces").HandlerFunc(Namespaces()).Methods("GET")
	r.Path("/relationshipsets").HandlerFunc(RelationshipSets()).Methods("GET")
	return r
}
