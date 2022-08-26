package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	utils.CreateHttpResponse(w, http.StatusOK, "Ticketing system")
}
