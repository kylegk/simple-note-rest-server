package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/auth"
	"github.com/kylegk/notes/lib"
	"github.com/kylegk/notes/model"
	"net/http"
)

// CreateUser adds a user to the data store and returns a token
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sendErrorResponse(w, r, err)
			return
		}
	}()

	// TODO: Implement authorization check
	// If this were a production ready application, this is where we'd verify account creation access

	body := model.CreateUserRequest{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		err = fmt.Errorf(app.InvalidRequestError)
		return
	}

	userID, err := lib.InsertUserDB(body.User)
	if err != nil {
		return
	}

	token, err := auth.GenerateUserToken(userID)
	if err != nil {
		return
	}

	sendResponse(model.CreateUserResponse{UserID: userID, Token: token}, http.StatusOK, w)
}