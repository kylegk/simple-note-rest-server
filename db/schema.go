package db

import "github.com/hashicorp/go-memdb"

const (
	NotesTable = "notes"
	UsersTable = "users"
	UserNotesTable = "user_notes"

	IDIdx = "id"
	ContentIdx = "content_idx"
	ModifiedIdx = "modified_idx"
	UserIdx = "user_idx"

	NoteIDFld = "NoteID"
	ContentFld = "Content"
	ModifiedFld = "Modified"
	UserIDFld = "UserID"
	UserFld = "User"
)

// Schema defines the schema used for the go-memdb database
var Schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		NotesTable: {
			Name: NotesTable,
			Indexes: map[string]*memdb.IndexSchema{
				IDIdx: {
					Name:    IDIdx,
					Unique:  true,
					Indexer: &memdb.IntFieldIndex{Field: NoteIDFld},
				},
				ContentIdx:{
					Name:    ContentIdx,
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: ContentFld},
				},
				ModifiedIdx: {
					Name:         ModifiedIdx,
					Unique:       false,
					Indexer:      &memdb.StringFieldIndex{Field: ModifiedFld},
					AllowMissing: true,
				},
			},
		},
		UsersTable: {
			Name: UsersTable,
			Indexes: map[string]*memdb.IndexSchema{
				IDIdx: {
					Name:    IDIdx,
					Unique:  true,
					Indexer: &memdb.IntFieldIndex{Field: UserIDFld},
				},
				UserIdx: {
					Name:    UserIdx,
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: UserFld},
				},
			},
		},
		UserNotesTable: {
			Name: UserNotesTable,
			Indexes: map[string]*memdb.IndexSchema{
				IDIdx: {
					Name:    IDIdx,
					Unique:  true,
					Indexer: &memdb.IntFieldIndex{Field: NoteIDFld},
				},
				UserIdx: {
					Name:    UserIdx,
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: UserIDFld},
				},

			},
		},
	},
}