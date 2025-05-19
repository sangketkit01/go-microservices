package main

import (
	"net/http"

	"github.com/sangketkit01/logger-service/data"
)

type JSONPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request){
	/// read json into var
	var req JSONPayload
	_ = app.readJson(w, r, &req)

	// insert data
	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	response := jsonReponse{
		Error: false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, response)
}