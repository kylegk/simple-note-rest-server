package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initUsersTests() *mux.Router {
	app.Init()
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	return router
}

func TestCreateUser(t *testing.T) {
	router := initUsersTests()

	// Attempt to create a valid user
	createUser := model.CreateUserRequest{User: "test.account"}
	j, _ := json.Marshal(createUser)

	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("failed to create user, status code: %v", response.Code)
	}

	// Attempt to create the same user again
	request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 405
	have = response.Code
	want = 405
	if have != want {
		t.Errorf("failed to create user, status code: %v", response.Code)
	}
}