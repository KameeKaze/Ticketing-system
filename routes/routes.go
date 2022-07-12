package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)

const (
	PORT = ":3000"
)

func RoutesHandler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//routes
	r.Get("/", home)
	r.Post("/login", login)
	r.Post("/register", register)
	r.Post("/tickets/create", createTicket)
	r.Get("/tickets", allTickets)

	//start
	fmt.Println("Running on http://127.0.0.1" + PORT)
	http.ListenAndServe(PORT, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	utils.CreateHttpResponse(w, 200, "Ticketing system")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, 500, "Can't connect to database")
		return
	}
	defer database.Close()

	//decode body data
	loginData := &types.Login{}
	json.NewDecoder(r.Body).Decode(&loginData)

	// check if request was valid
	if utils.ValidateJSON(loginData) {
		utils.CreateHttpResponse(w, 400, "Invalid request")
		return
	}

	if database.CheckPassword(loginData.Username, loginData.Password) {
		utils.CreateHttpResponse(w, 200, "Succesful login")
	} else {
		utils.CreateHttpResponse(w, 401, "Invalid credentials")
	}

}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, 500, "Can't connect to database")
		return
	}
	defer database.Close()

	//decode body data
	registerData := &types.Register{}
	json.NewDecoder(r.Body).Decode(&registerData)

	if utils.ValidateJSON(registerData) {
		utils.CreateHttpResponse(w, 400, "Invalid request")
		return
	}

	//check if user exoist
	if database.CheckUserExist(registerData.Username) {
		utils.CreateHttpResponse(w, 409, "Username already taken")
		return
	} else {
		err = database.AddUser(registerData)
		if err != nil {
			utils.Logger.Error(err.Error())
			utils.CreateHttpResponse(w, 400, "Invalid request")
			return
		} else {
			utils.CreateHttpResponse(w, 201, "Creating user "+registerData.Username)
		}
	}

}

func createTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, 500, "Can't connect to database")
	}
	defer database.Close()

	//decode body data
	data := &types.CreateTicket{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		utils.CreateHttpResponse(w, 400, "Invalid request")
		return
	}

	if database.AddTicket(data) != nil {
		utils.CreateHttpResponse(w, 500, "Can't create ticket")
		return
	} else {
		utils.CreateHttpResponse(w, 201, "Ticket successfully created")
	}
}

func allTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, 500, "Can't connect to database")
		return
	}
	defer database.Close()

	//decode body data
	getTickets := &types.GetTickets{}
	json.NewDecoder(r.Body).Decode(&getTickets)

	// check if request was valid
	if utils.ValidateJSON(getTickets) {
		utils.CreateHttpResponse(w, 400, "Invalid request")
		return
	}

}
