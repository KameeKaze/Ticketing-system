package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KameeKaze/Ticketing-system/utils"
)

const (
	PORT = ":3000"
)

func RoutesHandler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//routes
	r.Get("/", Home)
	r.Post("/login", Login)
	r.Post("/register", SignUp)
	r.Post("/changepassword", ChangePassword)
	r.Post("/tickets/create", CreateTicket)
	r.Delete("/tickets/{id}", DeleteTicket)
	r.Delete("/logout", Logout)
	r.Get("/tickets", AllTickets)

	//start
	fmt.Println("Running on http://127.0.0.1" + PORT)
	http.ListenAndServe(PORT, r)
}

func Home(w http.ResponseWriter, r *http.Request) {
	utils.CreateHttpResponse(w, http.StatusOK, "Ticketing system")
}
