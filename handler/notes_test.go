package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/lib"
	"github.com/kylegk/notes/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initNotesTest() *mux.Router {
	app.Init()
	router := mux.NewRouter()
	router.HandleFunc("/notes", GetAllNotesForUser).Methods("GET")
	router.HandleFunc("/notes", CreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}", GetNote).Methods("GET")
	router.HandleFunc("/notes/{id}", UpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", DeleteNote).Methods("DELETE")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	return router
}

// Create a user to test with
func createTestUser(router *mux.Router, userName string) (model.CreateUserResponse, error) {
	var user model.CreateUserResponse
	createUser := model.CreateUserRequest{User: userName}
	j, _ := json.Marshal(createUser)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	err := json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("failed to parse response returned when creating note")
	}

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		return user, fmt.Errorf("failed to create user")
	}

	return user, nil
}

func createValidTestNote(router *mux.Router, note model.CreateNoteRequest, token string) (int, error) {
	var createdNote model.CreateNoteResponse
	j, _ := json.Marshal(note)
	request, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have := response.Code
	want := 200
	if have != want {
		return 0, fmt.Errorf("create should have succeeded, have: %v, want: %v", have, want)
	}

	err := json.NewDecoder(response.Body).Decode(&createdNote)
	if err != nil {
		return 0, fmt.Errorf("failed to parse create note response")
	}

	return createdNote.NoteID, nil
}

func TestCreateNote(t *testing.T) {
	router := initNotesTest()
	user, err := createTestUser(router, "test.account")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Attempt to create a note without a token attached
	note := model.CreateNoteRequest{Content: "This is a test note"}
	j, _ := json.Marshal(note)
	request, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have := response.Code
	want := 405
	if have != want {
		t.Errorf("create should have failed with an error code, have: %v, want: %v", have, want)
	}

	_, err = createValidTestNote(router, note, user.Token)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateNote(t *testing.T) {
	router := initNotesTest()

	// Create a user
	user, err := createTestUser(router, "test.account")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a note for the user
	note := model.CreateNoteRequest{Content: "This is a test note"}
	newNoteID, err := createValidTestNote(router, note, user.Token)
	if err != nil {
		t.Errorf(err.Error())
	}

	url := "/notes/" + fmt.Sprintf("%d", newNoteID)

	// Try to update the note without a token
	note = model.CreateNoteRequest{Content: "This is a test note"}
	j, _ := json.Marshal(note)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have := response.Code
	want := 405
	if have != want {
		t.Errorf("update should have failed with an error code, have: %v, want: %v", have, want)
	}

	// Update the note with the correct token
	updateNote := model.UpdateNoteRequest{Content: "This is an updated note"}
	j, _ = json.Marshal(updateNote)
	request, _ = http.NewRequest("PUT", url, bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 200
	if have != want {
		t.Errorf("update should have failed with an error code, have: %v, want: %v", have, want)
	}

	// Try to update a note that doesn't exist
	updateNote.Content = "This is an updated note"
	j, _ = json.Marshal(updateNote)
	request, _ = http.NewRequest("PUT", "/notes/999", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 405
	if have != want {
		t.Errorf("update should have failed with an error code, have: %v, want: %v", have, want)
	}

	// Insert a user note that doesn't belong to our user
	_ = lib.InsertUserNoteDB(12345, 777)

	// Try to update the note that doesn't belong to our user
	updateNote.Content = "This is an updated note"
	j, _ = json.Marshal(note)
	request, _ = http.NewRequest("PUT", "/notes/777", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 405
	if have != want {
		t.Errorf("update should have failed with an error code, have: %v, want: %v", have, want)
	}
}

func TestGetNote(t *testing.T) {
	router := initNotesTest()

	// Create a user
	user, err := createTestUser(router, "test.account")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a note for the user
	note := model.CreateNoteRequest{Content: "This is a test note"}
	newNoteID, err := createValidTestNote(router, note, user.Token)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Get the note we just created
	url := "/notes/" + fmt.Sprintf("%d", newNoteID)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	have := response.Code
	want := 200
	if have != want {
		t.Errorf("unable to get note, have: %v, want: %v", have, want)
	}

	// Try to get a note that doesn't exist
	request, _ = http.NewRequest("GET", "/notes/1111", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 405
	if have != want {
		t.Errorf("get note should have failed with an error code, have: %v, want: %v", have, want)
	}

	// Insert a user note that doesn't belong to our user
	_ = lib.InsertNoteDB(999, "this is a note")
	_ = lib.InsertUserNoteDB(111, 999)

	// Try to get the note that doesn't belong to our user
	request, _ = http.NewRequest("GET", "/notes/999", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 405
	if have != want {
		t.Errorf("get note should have failed with an error code, have: %v, want: %v", have, want)
	}
}

func TestGetAllNotesForUser(t *testing.T) {
	router := initNotesTest()

	// Create a user
	user, err := createTestUser(router, "test.account")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Insert some rows into the user notes table
	noteIDs := []int{1,2,3,4}
	for _, noteID := range noteIDs {
		userNote := model.UserNote{UserID: user.UserID, NoteID: noteID}
		err := lib.InsertUserNoteDB(userNote.UserID, userNote.NoteID)
		if err != nil {
			t.Errorf("insertion of user note failed: %s", err.Error())
		}
	}

	// Get all notes that belong to the user
	request, _ := http.NewRequest("GET", "/notes", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("get all notes should have succeeded, have: %v, want: %v", have, want)
	}

	// Verify the number of notes that are owned by the user is correct
	r := model.GetAllNotesForUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		t.Errorf("unable to parse response")
	}
	have = len(r.Notes)
	want = len(noteIDs)

	if have != want {
		t.Errorf("number of notes returned is mismatched, have: %v, want: %v", have, want)
	}
}

func TestDeleteNote(t *testing.T) {
	router := initNotesTest()

	// Create a user
	user, err := createTestUser(router, "test.account")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create some notes for the user
	var note model.CreateNoteRequest
	for i := 0; i < 10; i++ {
		note = model.CreateNoteRequest{Content: "This is a test note"}
		_, err = createValidTestNote(router, note, user.Token)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	// Get all notes that belong to the user
	request, _ := http.NewRequest("GET", "/notes", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("get all notes should have succeeded, have: %v, want: %v", have, want)
	}
	r := model.GetAllNotesForUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		t.Errorf("unable to parse response")
	}
	createdNoteCount := len(r.Notes)

	// Delete one of the notes that was just created
	request, _ = http.NewRequest("DELETE", "/notes/5", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 200
	if have != want {
		t.Errorf("get all notes should have succeeded, have: %v, want: %v", have, want)
	}

	// Get all notes that belong to the user again
	request, _ = http.NewRequest("GET", "/notes", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + user.Token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	have = response.Code
	want = 200
	if have != want {
		t.Errorf("get all notes should have succeeded, have: %v, want: %v", have, want)
	}
	r = model.GetAllNotesForUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		t.Errorf("unable to parse response")
	}
	updatedNoteCount := len(r.Notes)

	if createdNoteCount == updatedNoteCount {
		t.Errorf("counts match, when they should be different")
	}
}