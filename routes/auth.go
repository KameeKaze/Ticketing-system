package routes

import (
	"encoding/json"
	"net/http"

	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)

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

	// comprare password and user password
	if database.CheckPassword(loginData.Username, loginData.Password) {
		//generate cookie
		cookie := utils.GenerateSessionCookie()
		//get userId
		userId := database.GetUserId(loginData.Username)

		// update or create session based on user already has a session
		if hasCookie, err := database.UserHasSession(userId); err != nil {
			//check database error
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			utils.Logger.Error(err.Error())
		} else if hasCookie {
			err = database.UpdateCookie(userId, cookie.Value, &cookie.Expires)
		} else {
			err = database.SaveCookie(userId, cookie.Value, &cookie.Expires)
		}
		// check error creating new session
		if err != nil {
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			utils.Logger.Error(err.Error())
			return
		} else {
			http.SetCookie(w, cookie)
			utils.CreateHttpResponse(w, http.StatusOK, "Logging in")
		}

	} else {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
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
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}

	//check if session exist
	if !(database.SessionExist(cookie.Value)) {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	} else {
		//delete cookie
		if err := database.DeleteCookie(cookie.Value); err != nil {
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			utils.Logger.Error(err.Error())
			return
		}
		utils.CreateHttpResponse(w, http.StatusResetContent, "Logging out")
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		utils.Logger.Error(err.Error())
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
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		} else {
			utils.CreateHttpResponse(w, http.StatusCreated, "Creating user "+registerData.Username)
			return
		}
	}

}
