package auth

import (
	"net/http"
)

const AuthContextKey = "auth"

func GetUser(req *http.Request) User {
	user, ok := req.Context().Value(AuthContextKey).(User)
	if !ok {
		panic("Could not get the User out of the context")
	}
	return user
}
