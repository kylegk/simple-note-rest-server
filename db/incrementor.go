package db

var noteIDIncrementer int
var userIDIncrementer int

func GetCurrentNoteID() int {
	return noteIDIncrementer
}

func IncrementNoteID() int {
	noteIDIncrementer++
	return noteIDIncrementer
}

func GetCurrentUserID() int {
	return userIDIncrementer
}

func IncrementUserID() int {
	userIDIncrementer++
	return userIDIncrementer
}