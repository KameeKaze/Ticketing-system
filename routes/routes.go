package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

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
	r.Delete("/logout", logout)
	r.Get("/tickets", allTickets)

	//start
	fmt.Println("Running on http://127.0.0.1" + PORT)
	http.ListenAndServe(PORT, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	utils.CreateHttpResponse(w, http.StatusOK, "Ticketing system")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		return
	}
	defer database.Close()

	//decode body data
	loginData := &types.Login{}
	json.NewDecoder(r.Body).Decode(&loginData)

	// check if request was valid
	if utils.ValidateJSON(loginData) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if database.CheckPassword(loginData.Username, loginData.Password) {
		//check if cookie already set
		cookie, err := r.Cookie("session")
		if err != nil || !(database.ValidateUserSession(cookie.Value, database.GetUserId(loginData.Username))){
			setSessionCookie(w, database.GetUserId(loginData.Username))
			utils.CreateHttpResponse(w, http.StatusOK, "Succesful login")
			return
		}else{ // validate cookie
				utils.CreateHttpResponse(w, http.StatusNotAcceptable, "Already logged in")
				return
		}
	} else {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}


func setSessionCookie(w http.ResponseWriter, userId string){
	database, err := db.ConnectDB()
	if err != nil {
		return
	}
	defer database.Close()

	// generate http cookie
	cookie := &http.Cookie{
		Name:     "session",
		Value:    uuid.New().String(),
		HttpOnly: true,
		Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
		Path:     "/",
	}
	// save or update cookie in database
	if !(database.SessionExist(cookie.Value)){
		err = database.SaveCookie(userId, cookie.Value, &cookie.Expires)
	}else{
		err = database.UpdateCookie(userId, cookie.Value, &cookie.Expires)
	}
	if err != nil{
		utils.Logger.Error(err.Error())
	}

	// set cookie
	http.SetCookie(w, cookie)
}

func logout(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		return
	}
	defer database.Close()

	cookie, err := r.Cookie("session")
	if err != nil || !(database.ValidateSession(cookie.Value)){
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Not logged in")
		return
	}else{
		//delete cookie
		if err := database.DeleteCookie(cookie.Value); err != nil{
			utils.Logger.Error(err.Error())
		}

		utils.CreateHttpResponse(w, http.StatusResetContent, "Logging out")

	}
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		return
	}
	defer database.Close()

	//decode body data
	registerData := &types.Register{}
	json.NewDecoder(r.Body).Decode(&registerData)

	if utils.ValidateJSON(registerData) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	//check if user exist
	if database.CheckUserExist(registerData.Username) {
		utils.CreateHttpResponse(w, http.StatusConflict, "Username already taken")
		return
	} else {
		err = database.AddUser(registerData)
		if err != nil {
			utils.Logger.Error(err.Error())
			utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
			return
		} else {
			utils.CreateHttpResponse(w, http.StatusCreated, "Creating user "+registerData.Username)
		}
	}

}

func createTicket(w http.ResponseWriter, r *http.Request) {
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

	if database.AddTicket(data) != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't create ticket")
		return
	} else {
		utils.CreateHttpResponse(w, http.StatusCreated, "Ticket successfully created")
	}
}

func allTickets(w http.ResponseWriter, r *http.Request) {
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

	tickets := database.GetAllTickets(r.URL.Query()["user"])

	j, err := json.Marshal(tickets)
	w.Write([]byte(j))
}