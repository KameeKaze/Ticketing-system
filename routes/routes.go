package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)

func RoutesHandler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//routes
	r.Get("/", Home)
	r.Post("/api/login", Login)
	r.Post("/api/register", SignUp)
	r.Post("/api/changepassword", ChangePassword)
	r.Post("/api/tickets", CreateTicket)
	r.Get("/api/tickets", AllTickets)
	r.Put("/api/tickets/{id}", UpdateTicket)
	r.Put("/api/tickets/{id}/{status}", UpdateTicketStatus)
	r.Delete("/api/tickets/{id}", DeleteTicket)
	r.Delete("/api/logout", Logout)

	//start
	fmt.Println("Running on http://127.0.0.1:" + os.Getenv("PORT"))
	http.ListenAndServe(":3000", r)
}

func Home(w http.ResponseWriter, r *http.Request) {
	createHttpResponse(w, http.StatusOK, "Ticketing system")
}

func checkHTTPRequest[T any](w http.ResponseWriter, r *http.Request, body *T) bool {
	//set headers
	w.Header().Set("Content-Type", "application/json")

	// get session cookie
	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		createHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return false
	}
	// get user from cookie
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return false
	}
	//user not logged in or not exist
	if userId == "" {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return false
	}

	json.NewDecoder(r.Body).Decode(&body)
	//check if body requered
	var ret T
	switch any(&ret).(type) {
	case any:
		return true
	}

	// check if request was valid
	if utils.ValidateJSON(body) {
		createHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return false
	}
	return true
}

func createHttpResponse(w http.ResponseWriter, statusCode int, text string) {
	//set status code
	w.WriteHeader(statusCode)
	//create json
	r, _ := json.Marshal(types.ResponseBody{
		Msg: text,
	})
	//send data
	w.Write([]byte(r))
}
