package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func viewRank(w http.ResponseWriter, r *http.Request) {
  vars:= mux.Vars(r)
  id := vars["id"]
	fmt.Fprintf(w, "got rank %s", id)
}

func viewLeaderboard(w http.ResponseWriter, r *http.Request) {
  vars:= mux.Vars(r)
  id := vars["id"]
	fmt.Fprintf(w, "got leaderboard %s", id)
}

func updateRank(w http.ResponseWriter, r *http.Request) {
  vars:= mux.Vars(r)
  id := vars["id"]
	fmt.Fprintf(w, "updating rank %s", id)
}

func main() {
	r := mux.NewRouter()
  r.HandleFunc("/leaderboard", viewLeaderboard).Methods("GET")
	r.HandleFunc("/rank/{id}", viewRank).Methods("GET")
	r.HandleFunc("/rank/{id}", updateRank).Methods("POST")

	http.ListenAndServe(":8080", r)
}
