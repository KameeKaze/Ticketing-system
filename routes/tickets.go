package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)


func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
	}
	defer database.Close()

	ticketId := chi.URLParam(r, "id")
	// delete ticket in database
	if err := database.DeleteTicket(ticketId); err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Successfuly deleted ticket '"+ticketId+"'")
		return
	}
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
	}
	defer database.Close()



	//decode body data
	data := &types.CreateTicket{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	data.Issuer = database.CookieUserId(cookie.Value)

	// create ticket
	if added, err := database.AddTicket(data); err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else if added {
		utils.CreateHttpResponse(w, http.StatusCreated, "Ticket successfully created")
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't create ticket")
		return
	}
}

func AllTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//query user parameter

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		return
	}
	defer database.Close()

	tickets, err := database.GetAllTickets(r.URL.Query()["user"])
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	}

	j, err := json.Marshal(tickets)
	w.Write([]byte(j))
}
