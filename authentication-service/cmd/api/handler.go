package main

import (
	"bytes"
	"encoding/json"
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

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil{
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	_ = app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error{
	var entry struct{
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "  ")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil{
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil{
		return err
	}

	return nil
}
