package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-redis/redis"
)

var a app

func TestGetRank(t *testing.T) {
	a = app{}
	a.initialize()

	req, _ := http.NewRequest("GET", "/user/5", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, `{"id":"5"}`, response.Body.String())

	clearLeaderboard()
}

func TestUpdateRank(t *testing.T) {
	a = app{}
	a.initialize()

	data := strings.NewReader(`{"score":123}`)
	req, _ := http.NewRequest("POST", "/user/5", data)
	response := executeRequest(req)

	req, _ = http.NewRequest("GET", "/user/5", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, `{"id":"5"}`, response.Body.String())

	data = strings.NewReader(`{"score":121}`)
	req, _ = http.NewRequest("POST", "/user/6", data)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, `""`, response.Body.String())

	req, _ = http.NewRequest("GET", "/user/5", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, `{"id":"5","rank":1}`, response.Body.String())

	clearLeaderboard()
}

func TestTopRanks(t *testing.T) {
	a = app{}
	a.initialize()

	insertUsers(10)

	req, _ := http.NewRequest("GET", "/topusers", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedBody := `[{"id":"9","score":10},{"id":"8","score":9},{"id":"7","score":8}]`
	checkResponseBody(t, expectedBody, response.Body.String())

	clearLeaderboard()
}

func TestViewLeaderboard(t *testing.T) {
	a = app{}
	a.initialize()

	insertUsers(10)

	req, _ := http.NewRequest("GET", "/leaderboard/5", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedBody := `[{"id":"4","score":5},{"id":"5","score":6},{"id":"6","score":7}]`
	checkResponseBody(t, expectedBody, response.Body.String())

	clearLeaderboard()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkResponseBody(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected response body %q. Got %q\n", expected, actual)
	}
}

func clearLeaderboard() {
	a.Redis.Del("leaderboard")
}

func insertUsers(number int) {

	for i := 0; i < number; i++ {
		a.Redis.ZAdd("leaderboard", redis.Z{float64(i + 1), i})
	}
}
