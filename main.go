package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func (a *app) viewRank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	zRank := a.Redis.ZRank("leaderboard", id)
	user := user{ID: id, Rank: zRank.Val()}

	respondWithJSON(w, http.StatusOK, user)
}

func (a *app) viewLeaderboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	zRank := a.Redis.ZRank("leaderboard", id)

	lower := zRank.Val() - 5
	upper := zRank.Val() + 4

	zRangeWithScores := a.Redis.ZRangeWithScores("leaderboard", lower, upper)

	users := []user{}
	for _, data := range zRangeWithScores.Val() {
		member, _ := data.Member.(string)

		user := user{ID: data.Member, Score: data.Score}
		users = append(users, user)
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *app) updateRank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	zAdd := a.Redis.ZAdd("leaderboard", redis.Z{100, id})
	user := user{ID: id, Rank: zAdd.Val()}

	respondWithJSON(w, http.StatusOK, user)
}

func (a *app) topRanks(w http.ResponseWriter, r *http.Request) {
	zRevRangeWithScores := a.Redis.ZRevRangeWithScores("leaderboard", 0, 9)

	users := []user{}
	for _, data := range zRevRangeWithScores.Val() {
		member, _ := data.Member.(string)

		user := user{ID: member, Score: data.Score}
		users = append(users, user)
	}

	respondWithJSON(w, http.StatusOK, users)
}

func respondWithJSON(w http.ResponseWriter, code int, user interface{}) {
	response, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type user struct {
	ID    string  `json:"id"`
	Rank  int64   `json:"rank"`
	Score float64 `json:"score"`
}

type app struct {
	Router *mux.Router
	Redis  *redis.Client
}

func (a *app) initialize() {
	a.Router = mux.NewRouter()
	a.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	a.initializeRoutes()
}

// "/product/{id:[0-9]+}"
// vars := mux.Vars(r)
// id, err := strconv.Atoi(vars["id"])
// if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
// 		return
// }
func (a *app) initializeRoutes() {
	a.Router.HandleFunc("/leaderboard/{id}", a.viewLeaderboard).Methods("GET")
	a.Router.HandleFunc("/rank/{id}", a.viewRank).Methods("GET")
	a.Router.HandleFunc("/rank/{id}", a.updateRank).Methods("POST")
	a.Router.HandleFunc("/topranks", a.topRanks).Methods("GET")
}

func (a *app) run() {
	http.ListenAndServe(":8080", a.Router)
}

func main() {
	a := app{}
	a.initialize()
	a.run()
}
