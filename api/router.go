package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Router(r *mux.Router, jsonFile []byte) {
	r.HandleFunc("/products", fncProduct(jsonFile))
}

func fncProduct(jsonFile []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonFile))
	}
}
