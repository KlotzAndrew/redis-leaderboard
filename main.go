package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func viewRank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client := redisClient()
	zRank := client.ZRank("leaderboard", id)

	fmt.Fprintf(w, "got rank %d", zRank.Val())
}

func viewLeaderboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client := redisClient()
	zRank := client.ZRank("leaderboard", id)

	lower := zRank.Val() - 5
	upper := zRank.Val() + 4

	zRangeWithScores := client.ZRangeWithScores("leaderboard", lower, upper)

	fmt.Fprintf(w, "got leaderboard %v", zRangeWithScores.Val())
}

func updateRank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client := redisClient()
	zAdd := client.ZAdd("leaderboard", redis.Z{100, id})

	fmt.Fprintf(w, "updating rank %d", zAdd.Val())
}

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/leaderboard/{id}", viewLeaderboard).Methods("GET")
	r.HandleFunc("/rank/{id}", viewRank).Methods("GET")
	r.HandleFunc("/rank/{id}", updateRank).Methods("POST")

	http.ListenAndServe(":8080", r)
}
