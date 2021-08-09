package handler

import (
	"encoding/json"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/model"
	"log"
	"net/http"
)

// SendGenericNotFoundResponse returns a generic 404 error
func SendGenericNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&model.GenericResponse{Error: "Not Found", Code: http.StatusNotFound, Message: "Resource not found"}, http.StatusNotFound, w)
}

// SendGenericNotAllowedResponse returns a generic 405 error
func SendGenericNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&model.GenericResponse{Error: "Not Allowed", Code: http.StatusMethodNotAllowed, Message: "You are not authorized to access this resource"}, http.StatusMethodNotAllowed, w)
}

// SendGenericInternalServerError returns a generic 500 error
func SendGenericInternalServerError(w http.ResponseWriter, r *http.Request) {
	sendResponse(&model.GenericResponse{Error: "Internal Server Error", Code: http.StatusInternalServerError, Message: "An error has occurred"}, http.StatusInternalServerError, w)
}

// SendGenericBadRequestResponse returns a generic 400 error
func SendGenericBadRequestResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&model.GenericResponse{Error: "Bad Request", Code: http.StatusBadRequest, Message: "Invalid Request"}, http.StatusMethodNotAllowed, w)
}

// SendGenericNotAuthorizedResponse returns a generic 401 error
func SendGenericNotAuthorizedResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&model.GenericResponse{Error: "Not Authorized", Code: http.StatusUnauthorized, Message: "User is not authorized"}, http.StatusMethodNotAllowed, w)
}

func sendResponse(payload interface{}, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)

	switch err.Error() {
	case app.InvalidTokenError:
		SendGenericNotAuthorizedResponse(w, r)
	case app.UserExistsError, app.InvalidUserError, app.InvalidRequestError:
		SendGenericBadRequestResponse(w, r)
	default:
		SendGenericInternalServerError(w, r)
	}

}