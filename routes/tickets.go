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
	if !checkHTTPRequest(w, r, new(interface{})) {
		return
	}
	// get session cookie
	cookie, _ := r.Cookie("session")
	userId, _ := db.Redis.GetUserId(cookie.Value)

	ticketId := chi.URLParam(r, "id")

	if id, err := db.Mysql.GetTicketIssuer(ticketId); err != nil {
		createHttpResponse(w, http.StatusNotFound, "Ticket with id does not exit")
		return
	} else if id != userId {
		createHttpResponse(w, http.StatusUnauthorized, "You can't delete this ticket")
		return
	}
	// delete ticket in database
	if err := db.Mysql.DeleteTicket(ticketId); err != nil {
		createHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else {
		createHttpResponse(w, http.StatusBadRequest, "Successfuly deleted ticket '"+ticketId+"'")
		utils.Logger.Info("Delete ticket: " + ticketId)

		return
	}
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	body := types.HTTPTicket{}
	if !checkHTTPRequest(w, r, &body) {
		return
	}
	// get session cookie
	cookie, _ := r.Cookie("session")
	userId, _ := db.Redis.GetUserId(cookie.Value)

	// create ticket
	body.Issuer = userId
	if added, err := db.Mysql.AddTicket(&body); err != nil {
		createHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	} else if added {
		utils.Logger.Info(body.Issuer + " successfully created a ticket")
		createHttpResponse(w, http.StatusCreated, "Ticket successfully created")
		return
	} else {
		createHttpResponse(w, http.StatusInternalServerError, "Can't create ticket")
		return
	}
}

func AllTickets(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPRequest(w, r, new(interface{})) {
		return
	}

	tickets, err := db.Mysql.GetAllTickets(r.URL.Query()["user"])
	if err != nil {
		createHttpResponse(w, http.StatusBadRequest, "Database error")
		utils.Logger.Error(err.Error())
		return
	}

	j, _ := json.Marshal(tickets)
	w.Write([]byte(j))
}

func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	body := types.HTTPTicket{}
	if !checkHTTPRequest(w, r, &body) {
		return
	}

	err := db.Mysql.UpdateTicket(chi.URLParam(r, "id"), &body)
	if err != nil {
		utils.Logger.Error(err.Error())
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	createHttpResponse(w, http.StatusNoContent, "")
}

func UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPRequest(w, r, new(interface{})) {
		return
	}
	status := chi.URLParam(r, "status")

	switch status {
	case "inprog":
		InProgTicket(w, r)
	case "closed":
		CloseTicket(w, r)
	}

}

func CloseTicket(w http.ResponseWriter, r *http.Request) {
	if err := db.Mysql.UpdateTicketStatus(2, chi.URLParam(r, "id")); err != nil {
		utils.Logger.Error(err.Error())
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	} else {
		createHttpResponse(w, http.StatusNoContent, "")
		return
	}

}

func InProgTicket(w http.ResponseWriter, r *http.Request) {
	if err := db.Mysql.UpdateTicketStatus(1, chi.URLParam(r, "id")); err != nil {
		utils.Logger.Error(err.Error())
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		return
	} else {
		createHttpResponse(w, http.StatusNoContent, "")
		return
	}

}
