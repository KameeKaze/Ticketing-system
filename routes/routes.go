package routes

import (
	"net/http"
	"fmt"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/utils"

)

const(
	PORT = ":3000"
)

func RoutesHandler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//routes
	r.Get("/",home)
	r.Post("/login",login)

	//start
	fmt.Println("Running on http://127.0.0.1" + PORT)
	http.ListenAndServe(PORT, r)
}

func home(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Ticketing system\n"))
	database, err := db.ConnectDB()
	if err != nil{
		utils.Logger.Error(err.Error())
	}
	defer database.Close()
	w.Write([]byte("Connected to database"))
}

func login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil{
		utils.Logger.Error(err.Error())
	}
	defer database.Close()

	//decode body data 
	loginData := &types.Login{}
	json.NewDecoder(r.Body).Decode(&loginData)

	// check if request was valid
	if loginData.Username == "" || loginData.Password == "" {
		w.Write([]byte("Invalid request"))
		return
	}

	if database.CheckPassword(loginData.Username, loginData.Password){
		w.Write([]byte("Succesful login"))
	}else{
		w.Write([]byte("Invalid login"))
	}
	
}
