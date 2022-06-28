package routes

import (
	"net/http"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/utils"

)

const(
	PORT = ":3000"
)

func RoutesHandler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/",home)
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