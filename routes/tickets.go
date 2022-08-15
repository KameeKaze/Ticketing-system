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

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}
	ticketId := chi.URLParam(r, "id")

	if id, err := db.Mysql.GetTicketIssuer(ticketId); err != nil {
		utils.CreateHttpResponse(w, http.StatusNotFound, "Ticket with id does not exit")
		return
	} else if id != userId {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "You can't delete this ticket")
		return
	}
	// delete ticket in database
	if err := db.Mysql.DeleteTicket(ticketId); err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Successfuly deleted ticket '"+ticketId+"'")
		utils.Logger.Info("Delete ticket: " + ticketId)

		return
	}
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	// create ticket
	data.Issuer = userId
	if added, err := db.Mysql.AddTicket(data); err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else if added {
		utils.Logger.Info(data.Issuer + " successfully created a ticket")
		utils.CreateHttpResponse(w, http.StatusCreated, "Ticket successfully created")
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't create ticket")
		return
	}
}

func AllTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	tickets, err := db.Mysql.GetAllTickets(r.URL.Query()["user"])
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	}

	j, _ := json.Marshal(tickets)
	w.Write([]byte(j))
}

func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//query user parameter

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	data := &types.CreateTicket{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	err = db.Mysql.UpdateTicket(chi.URLParam(r, "id"), data)
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.CreateHttpResponse(w, http.StatusNoContent, "")
}

func CloseTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//query user parameter

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}
	if db.Mysql.UpdateTicketStatus(2, chi.URLParam(r, "id")) != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusNoContent, "")
		return
	}

}

func InProgTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//query user parameter

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}
	if db.Mysql.UpdateTicketStatus(1, chi.URLParam(r, "id")) != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusNoContent, "")
		return
	}

}
