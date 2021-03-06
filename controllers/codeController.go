package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/l3njo/dropnote-api/models"
	u "github.com/l3njo/dropnote-api/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GenerateCode is the handler function for creating a new reset token
func GenerateCode(w http.ResponseWriter, r *http.Request) {
	body := &struct {
		Mail string `json:"mail"`
	}{}
	err := json.NewDecoder(r.Body).Decode(body)
	uData := models.GetUserByMail(App.DB, body.Mail)
	if uData == nil {
		u.Respond(w, u.Message(true, "success"))
		return
	}
	user := uData.ID
	code, err := models.NewCode(App.DB, models.Actions["reset"], user)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}
	if err := emailTokenToUser(uData, code); err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Unable to send email"))
		return
	}
	resp := u.Message(true, "success")
	u.Respond(w, resp)
}

// ExecuteCode is the handler function for carrying out a token action
func ExecuteCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := uuid.FromString(params["user"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	code, err := uuid.FromString(params["code"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	if err = models.Execute(App.DB, r.Body, code, user); err != nil {
		log.Println(err.Error())
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	resp := u.Message(true, "success")
	u.Respond(w, resp)
}
