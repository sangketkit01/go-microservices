package main

import "net/http"

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request){
	type mailMessage struct{
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var req mailMessage
	err := app.readJson(w, r, &req)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	msg := Message{
		From: req.From,
		To: req.To,
		Subject: req.Subject,
		Data: req.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil{
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "sent to " + req.To,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}