package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/auth"
	"github.com/kylegk/notes/db"
	"github.com/kylegk/notes/lib"
	"github.com/kylegk/notes/model"
	"net/http"
	"strconv"
)

// CreateNote handles the request to insert a note into the data store and create a relationship between a user and their note
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	userID, err := auth.ValidateUserToken(r)
	if err != nil {
		return
	}

	body := model.CreateNoteRequest{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	noteID := db.IncrementNoteID()
	err = lib.InsertNoteDB(noteID, body.Content)
	if err != nil {
		return
	}

	err = lib.InsertUserNoteDB(userID, noteID)
	if err != nil {
		return
	}

	sendResponse(model.CreateNoteResponse{NoteID: noteID}, http.StatusOK, w)
}

// UpdateNote handles the request to update a single note
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	userID, err := auth.ValidateUserToken(r)
	if err != nil {
		return
	}

	id := mux.Vars(r)["id"]
	noteID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	body := model.UpdateNoteRequest{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	err = lib.ValidateNoteOwnershipDB(userID, noteID)
	if err != nil {
		return
	}

	err = lib.UpdateNoteDB(noteID, body.Content)
	if err != nil {
		return
	}

	sendResponse(model.GenericResponse{Message: "Note updated"}, http.StatusOK, w)
}

// GetNote handles the request to retrieve a single note
func GetNote(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	userID, err := auth.ValidateUserToken(r)
	if err != nil {
		return
	}

	id := mux.Vars(r)["id"]
	noteID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	err = lib.ValidateNoteOwnershipDB(userID, noteID)
	if err != nil {
		return
	}

	note, err := lib.GetNoteDB(noteID)
	if err != nil {
		return
	}

	sendResponse(note, http.StatusOK, w)
}

// GetAllNotesForUser gets all the notes associated with a user
func GetAllNotesForUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	userID, err := auth.ValidateUserToken(r)
	if err != nil {
		return
	}

	noteIDs, err := lib.GetAllNotesForUserDB(userID)
	if err != nil {
		return
	}

	sendResponse(model.GetAllNotesForUserResponse{Notes: noteIDs}, http.StatusOK, w)
}

// DeleteNote handles the request to delete a note and the user's relationship to that note
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	id := mux.Vars(r)["id"]
	noteID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	userID, err := auth.ValidateUserToken(r)
	if err != nil {
		return
	}

	err = lib.ValidateNoteOwnershipDB(userID, noteID)
	if err != nil {
		return
	}

	_, err = lib.DeleteNoteDB(noteID)
	if err != nil {
		return
	}

	_, err = lib.DeleteUserNoteDB(noteID)
	if err != nil {
		return
	}

	sendResponse(model.GenericResponse{Message: "Note deleted"}, http.StatusOK, w)
}