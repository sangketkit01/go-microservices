package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct{
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct{
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	
	if err := app.readJson(w, r, &requestPayload) ; err != nil{
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action{
	case "auth" :
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload){
	jsonData, _ := json.MarshalIndent(entry, "", "  ")

	logServerUrl := "http://logger-service/log"

	request, err := http.NewRequest(http.MethodPost, logServerUrl, bytes.NewBuffer(jsonData))
	if err != nil{
		app.errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "application.json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted{
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth service
	jsonData, _ := json.MarshalIndent(a, "", "  ")

	// call the service
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil{
		app.errorJson(w, err)
		return
	}

	client := &http.Client{
		
	}
	response, err := client.Do(request)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized{
		app.errorJson(w, errors.New("invalid credentials"))
		return
	}else if response.StatusCode != http.StatusAccepted{
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse
	
	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	if jsonFromService.Error{
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "Authenticated!",
		Data: jsonFromService.Data,
	}

	_ = app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload){
	jsonData, _ := json.MarshalIndent(msg, "", "  ")
	
	// call mail service
	mailServiceUrl := "http://mail-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil{
		app.errorJson(w,  err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil{
		app.errorJson(w, err)
		return 
	}

	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted{
		app.errorJson(w, err)
		return 
	}

	// send back json
	payload := jsonResponse{
		Error: false,
		Message: "Message sent to " + msg.To,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}
