package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/l3njo/dropnote-api/models"
	u "github.com/l3njo/dropnote-api/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// CreateNote is the handler function for adding a note to the database
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var id string
	if r.Context().Value(UserKey) == nil {
		id = ""
	} else {
		id = (r.Context().Value(UserKey).(uuid.UUID)).String()
	}
	user := uuid.FromStringOrNil(id)
	note := &models.Note{}
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	note.Creator = user
	resp := note.Create(App.DB)
	u.Respond(w, resp)
}

// GetNote is the handler function for getting a note from the database
func GetNote(w http.ResponseWriter, r *http.Request) {
	var userID string
	if r.Context().Value(UserKey) == nil {
		userID = ""
	} else {
		userID = (r.Context().Value(UserKey).(uuid.UUID)).String()
	}
	params := mux.Vars(r)
	id, err := uuid.FromString(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	data := models.GetNote(App.DB, id)
	if *data != (models.Note{}) {
		user := uuid.FromStringOrNil(userID)
		if ok := uuid.Equal(user, data.Creator); !data.Visible && !ok {
			data = nil
		}
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

// GetNotes is the handler function for getting all notes from database
func GetNotes(w http.ResponseWriter, r *http.Request) {
	data := models.GetNotes(App.DB)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
