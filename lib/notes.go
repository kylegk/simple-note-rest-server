package lib

import (
	"fmt"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/db"
	"github.com/kylegk/notes/model"
	"time"
)

// InsertNoteDB inserts the note into the data store
func InsertNoteDB(noteID int, body string) error {
	note := model.Note{
		NoteID:   noteID,
		Content:  body,
		Modified: time.Now().String(),
	}

	exists, err := GetNoteDB(noteID)
	if err != nil {
		return err
	}

	if exists.NoteID == noteID {
		return fmt.Errorf("cannot insert note; note already exists")
	}

	err = app.Context.DB.Upsert(db.NotesTable, note)
	if err != nil {
		return err
	}

	return nil
}

// UpdateNoteDB updates a note
func UpdateNoteDB(noteID int, body string) error {
	note := model.Note{
		NoteID:   noteID,
		Content:  body,
		Modified: time.Now().String(),
	}

	// Verify the row exists before attempting to modify
	n, err := app.Context.DB.Query(db.NotesTable, db.IDIdx, noteID)
	if len(n) == 0 {
		return fmt.Errorf("cannot update; row doesn't exist")
	}

	err = app.Context.DB.Upsert(db.NotesTable, note)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNoteDB deletes a note
func DeleteNoteDB(noteID int) (int, error) {
	count, err := app.Context.DB.Delete(db.NotesTable, db.IDIdx, noteID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetNoteDB retrieves a single note
func GetNoteDB(noteID int) (model.Note, error) {
	var note model.Note

	res, err := app.Context.DB.Query(db.NotesTable, db.IDIdx, noteID)
	if err != nil {
		return note, err
	}

	if len(res) == 0 {
		return note, nil
	}

	if len(res) > 1 {
		return note, fmt.Errorf("something went wrong, more than one note was found")
	}
	note = res[0].(model.Note)

	return note, nil
}

// InsertUserNoteDB creates the relationship between the user and a note
func InsertUserNoteDB(userID int, noteID int) error {
	userNote := model.UserNote{
		UserID: userID,
		NoteID: noteID,
	}

	// Verify that the note doesn't exist before attempting to insert
	res, err := app.Context.DB.Query(db.UserNotesTable, db.IDIdx, noteID)
	if err != nil {
		return err
	}
	if len(res) > 0 {
		return fmt.Errorf("cannot insert user note, note already exists")
	}

	err = app.Context.DB.Upsert(db.UserNotesTable, userNote)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserNoteDB deletes the relationship between the user and the note
func DeleteUserNoteDB(userID int) (int, error) {
	count, err := app.Context.DB.Delete(db.UserNotesTable, db.IDIdx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetAllNotesForUserDB retrieves all the note ids associated with the specified user
func GetAllNotesForUserDB(userID int) ([]int, error) {
	var noteIDs []int

	notes, err := app.Context.DB.Query(db.UserNotesTable, db.UserIdx, userID)
	if err != nil {
		return noteIDs, err
	}

	for _, userNote := range notes {
		noteIDs = append(noteIDs, userNote.(model.UserNote).NoteID)
	}

	return noteIDs, nil
}

// ValidateNoteOwnershipDB verifies the user attempting an action owns the note they're trying to act on
func ValidateNoteOwnershipDB(userID int, noteID int) error {
	res, err := app.Context.DB.Query(db.UserNotesTable, db.IDIdx, noteID)
	if err != nil {
		return err
	}

	if len(res) == 0  {
		return fmt.Errorf(app.InvalidRequestError)
	}

	if res[0].(model.UserNote).UserID != userID {
		return fmt.Errorf(app.InvalidTokenError)
	}

	return nil
}