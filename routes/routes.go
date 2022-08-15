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
	r.Post("/login", Login)
	r.Post("/register", SignUp)
	r.Post("/changepassword", ChangePassword)
	r.Put("/tickets/create", CreateTicket)
	r.Post("/tickets/{id}", UpdateTicket)
	r.Delete("/tickets/{id}", DeleteTicket)
	r.Post("/tickets/closed/{id}", CloseTicket)
	r.Post("/tickets/inprog/{id}", InProgTicket)
	r.Delete("/logout", Logout)
	r.Get("/tickets", AllTickets)

	//start
	fmt.Println("Running on http://127.0.0.1:" + os.Getenv("PORT"))
	http.ListenAndServe(":3000", r)
}

func Home(w http.ResponseWriter, r *http.Request) {
	utils.CreateHttpResponse(w, http.StatusOK, "Ticketing system")
}
