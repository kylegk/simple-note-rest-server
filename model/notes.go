package model

type Note struct {
	NoteID int
	Content string
	Modified string
}

type UserNote struct {
	UserID int
	NoteID int
}

// CreateNoteRequest defines the shape of the request used for creating notes
type CreateNoteRequest struct {
	Content string `json:"content"`
}

type CreateNoteResponse struct {
	NoteID int `json:"noteid"`
}

type UpdateNoteRequest struct {
	Content string
}

type GetAllNotesForUserResponse struct {
	Notes []int `json:"notes"`
}