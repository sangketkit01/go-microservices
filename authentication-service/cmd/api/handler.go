package main

import (
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	if err := app.readJson(w, r, &requestPayload) ; err != nil{
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	// validate the user againt the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil{
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid{
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	_ = app.writeJson(w, http.StatusAccepted, payload)
}
