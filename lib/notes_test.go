package lib

import (
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/model"
	"testing"
)

func TestInsertNoteDB(t *testing.T) {
	app.Init()

	note := model.Note{
		NoteID: 1,
		Content: "this is a test",
	}

	// Insert a valid note
	err := InsertNoteDB(note.NoteID, note.Content)
	if err != nil {
		t.Errorf("failed to insert note: %s", err.Error())
	}

	// Attempt to insert a duplicate record
	err = InsertNoteDB(note.NoteID, note.Content)
	if err == nil {
		t.Errorf("insert should have failed")
	}
}

func TestUpdateNoteDB(t *testing.T) {
	app.Init()

	note := model.Note{
		NoteID: 1,
		Content: "this is a test",
	}

	// Insert a valid note
	err := InsertNoteDB(note.NoteID, note.Content)
	if err != nil {
		t.Errorf("failed to insert note: %s", err.Error())
	}

	// Update the content of the note
	note.Content = "this is an updated note"
	err = UpdateNoteDB(note.NoteID, note.Content)
	if err != nil {
		t.Errorf("update note failed: %s", err.Error())
	}

	// Verify the content matches
	updatedNote, err := GetNoteDB(note.NoteID)
	if updatedNote.Content != note.Content {
		t.Errorf("update failed, strings don't match")
	}

	// Attempt to update a note that doesn't exist
	newNote := model.Note{
		NoteID: 2,
		Content: "This should fail",
	}
	err = UpdateNoteDB(newNote.NoteID, newNote.Content)
	if err == nil {
		t.Errorf("updated note that doesn't exist")
	}
}

func TestDeleteNoteDB(t *testing.T) {
	app.Init()

	note := model.Note{
		NoteID: 1,
		Content: "this is a test",
	}

	// Insert a valid note
	err := InsertNoteDB(note.NoteID, note.Content)
	if err != nil {
		t.Errorf("failed to insert note: %s", err.Error())
	}

	// Verify that delete succeeds with the correct number of records (1)
	want := 1
	have, err := DeleteNoteDB(note.NoteID)
	if err != nil {
		t.Errorf("failed to delete")
	}
	if have != want {
		t.Errorf("deleted an incorrect number of rows, have: %v, want: %v", have, want)
	}

	// Attempt to delete a note that doesn't exist
	have, err = DeleteNoteDB(99999)
	want = 0
	if have != want {
		t.Errorf("attempt to delete an invalid record succeeded")
	}
}

func TestInsertUserNoteDB(t *testing.T) {
	app.Init()

	userNote := model.UserNote{
		UserID: 1,
		NoteID: 1,
	}

	// Insert a user note
	err := InsertUserNoteDB(userNote.UserID, userNote.NoteID)
	if err != nil {
		t.Errorf("insertion of user note failed: %s", err.Error())
	}

	// Attempt to insert a duplicate user note
	err = InsertUserNoteDB(userNote.UserID, userNote.NoteID)
	if err == nil {
		t.Errorf("inserted user note, when it should have failed")
	}
}

func TestDeleteUserNoteDB(t *testing.T) {
	app.Init()

	userNote := model.UserNote{
		UserID: 1,
		NoteID: 1,
	}

	// Insert a user note
	err := InsertUserNoteDB(userNote.UserID, userNote.NoteID)
	if err != nil {
		t.Errorf("insertion of user note failed: %s", err.Error())
	}

	// Delete a valid user note
	want := 1
	have, err := DeleteUserNoteDB(userNote.UserID)
	if err != nil {
		t.Errorf("failed to delete user note")
	}
	if have != want {
		t.Errorf("deleted an incorrect number of rows, have: %v, want: %v", have, want)
	}

	// Attempt to delete a user note that doesn't exist
	userNote.NoteID = 9999
	want = 0
	have, err = DeleteUserNoteDB(userNote.NoteID)
	if err != nil {
		t.Errorf("error %s", err.Error())
	}
	if have != want {
		t.Errorf("deleted an incorrect number of rows, have: %v, want: %v", have, want)
	}
}

func TestGetAllNotesForUserDB(t *testing.T) {
	app.Init()

	// Insert some user notes
	userID := 1
	noteIDs := []int{1,2,3,4}
	for _, noteID := range noteIDs {
		userNote := model.UserNote{UserID: userID, NoteID: noteID}
		err := InsertUserNoteDB(userNote.UserID, userNote.NoteID)
		if err != nil {
			t.Errorf("insertion of user note failed: %s", err.Error())
		}
	}

	ownedIDs, err := GetAllNotesForUserDB(userID)
	if err != nil {
		t.Errorf("failed to retrieve ownership: %s", err.Error())
	}
	have := len(ownedIDs)
	want := len(noteIDs)
	if have != want {
		t.Errorf("incorrect number of notes owned, have: %v, want: %v", have, want)
	}
}

func TestValidateNoteOwnershipDB(t *testing.T) {
	app.Init()

	userNote := model.UserNote{
		UserID: 1,
		NoteID: 10,
	}

	// Insert a user note
	err := InsertUserNoteDB(userNote.UserID, userNote.NoteID)
	if err != nil {
		t.Errorf("insertion of user note failed: %s", err.Error())
	}

	// Verify a valid note owner
	err = ValidateNoteOwnershipDB(userNote.UserID, userNote.NoteID)
	if err != nil {
		t.Errorf("failed to find valid ownership")
	}

	// Verify an invalid note owner
	userNote.UserID = 5
	err = ValidateNoteOwnershipDB(userNote.UserID, userNote.NoteID)
	if err == nil {
		t.Errorf("failed to reject ownership")
	}
}