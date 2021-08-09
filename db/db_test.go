package db

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/kylegk/notes/model"
	"testing"
)

func initTestDB(schema *memdb.DBSchema) (DB, error) {
	var db DB
	db, err := InitDB(schema)
	if err != nil {
		return db, fmt.Errorf("failed to initialize database: %s", err.Error())
	}
	if db.Conn == nil {
		return db, fmt.Errorf("failed to initialize database")
	}

	return db, nil
}

func TestDB_Query(t *testing.T) {
	db, err := initTestDB(Schema)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test invalid table name
	_, err = db.Query("invalid_table", "id_idx")
	if err == nil {
		t.Errorf("invalid table query should have failed")
	}

	// Test valid table, invalid index
	_, err = db.Query(NotesTable, "invalid_index")
	if err == nil {
		t.Errorf("invalid index query should have failed")
	}

	// Test valid table, valid index
	_, err = db.Query(NotesTable, IDIdx)
	if err != nil {
		t.Errorf("valid table and index")
	}
}

func TestDB_Upsert(t *testing.T) {
	db, err := initTestDB(Schema)
	if err != nil {
		t.Errorf(err.Error())
	}

	note := model.Note{NoteID: 1, Content: "test note"}

	// Test invalid table name
	err = db.Upsert("invalid_table", note)
	if err == nil {
		t.Errorf("upsert should have failed with an invalid table name")
	}

	// Test invalid insert interface type
	err = db.Upsert(UserNotesTable, note)
	if err == nil {
		t.Errorf("upsert should have failed with an invalid interface type")
	}

	// Insert valid data
	err = db.Upsert(NotesTable, note)
	if err != nil {
		t.Errorf("upsert should have succeeded")
	}
}

func TestDB_Delete(t *testing.T) {
	db, err := initTestDB(Schema)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Insert data
	notes := []int{1,2,3,4,5}
	for _, noteID := range notes {
		note := model.Note{NoteID: noteID, Content: "test note"}
		err = db.Upsert(NotesTable, note)
		if err != nil {
			t.Errorf("failed to insert data")
		}
	}

	// Try to delete invalid table name
	_, err = db.Delete("invalid_table", "invalid_idx")
	if err == nil {
		t.Errorf("delete should have failed on invalid table")
	}

	// Delete record
	have, err := db.Delete(NotesTable, IDIdx, 1)
	if err != nil {
		t.Errorf("failed to delete note: %s", err.Error())
	}
	want := 1
	if have != want {
		t.Errorf("deleted the wrong number of notes, have: %v, want: %v", have, want)
	}
}