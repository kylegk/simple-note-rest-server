package lib

import (
	"fmt"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/db"
	"github.com/kylegk/notes/model"
)

// InsertUserDB inserts a user into the data store
func InsertUserDB(user string) (int, error) {
	userID := db.IncrementUserID()

	account := model.UserAccount{
		UserID: userID,
		User: user,
	}

	if user == "" {
		return 0, fmt.Errorf(app.InvalidUserError)
	}

	res, err := app.Context.DB.Query(db.UsersTable, db.UserIdx, user)
	if len(res) > 0 {
		err = fmt.Errorf(app.UserExistsError)
	}
	if err != nil {
		return 0, err
	}

	err = app.Context.DB.Upsert(db.UsersTable, account)
	if err != nil {
		return 0, err
	}

	return userID, nil
}