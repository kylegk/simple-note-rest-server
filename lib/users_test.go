package lib

import (
	"fmt"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/db"
	"testing"
)

func TestInsertUser(t *testing.T) {
	app.Init()

	userName := "test.account1"
	want := 1
	have, err := InsertUserDB(userName)
	if err != nil {
		fmt.Errorf("failed to insert user: %s", err.Error())
	}

	if have != want {
		fmt.Errorf("unexpected user id, have %v, want: %v", have, want)
	}

	// Insert a couple more test users and verify the latest ID is the value expected
	moreUsers := []string{"test.account2","test.account3"}
	for _, user := range moreUsers {
		have, err = InsertUserDB(user)
		if err != nil {
			fmt.Errorf("failed to insert user: %s", err.Error())
		}
		want = db.GetCurrentUserID()
		if have != want {
			fmt.Errorf("unexpected user id, have %v, want: %v", have, want)
		}
	}
}