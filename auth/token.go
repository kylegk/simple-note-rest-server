package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/db"
	"net/http"
	"strings"
	"time"
)

// TokenSecret is the string used to sign and validate the JWT for the application
// NOTE: If this were production code, the secret would be set in the environment or through an .env file
const TokenSecret = "THIS_IS_MY_SECRET"

// GenerateUserToken generates a JWT
func GenerateUserToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires": time.Now().Add(time.Minute * 15).Unix(),
		"userid": userID,
	})

	tokenString, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateUserToken validates a JWT
func ValidateUserToken(r *http.Request) (int, error) {
	tokenString := ExtractToken(r)

	if tokenString == "" {
		return 0, fmt.Errorf(app.InvalidTokenError)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(TokenSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf(app.InvalidTokenError)
	}

	var userID int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: Validate that the 'expires' claim is now or in the future and return an error if it has expired
		tmp := claims["userid"].(float64)
		userID = int(tmp)
	} else {
		return 0, fmt.Errorf(app.InvalidTokenError)
	}

	// Verify the user exists in the data store
	res, err := app.Context.DB.Query(db.UsersTable, db.IDIdx, userID)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, fmt.Errorf(app.InvalidTokenError)
	}

	return userID, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}