package routes

import (
	"encoding/json"
	"net/http"

	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
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
		if database.UserHasSession(userId) {
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
			utils.CreateHttpResponse(w, http.StatusOK, "Logging in "+loginData.Username)
		}

	} else {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Can't connect to database")
		utils.Logger.Error(err.Error())
		return
	}
	defer database.Close()

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Login as admin to create user")
		return
	}
	if !database.CookieAuthorized(cookie.Value){
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Only admins can create users")
		return
	}
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

func ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	data := &types.ChangePassword{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// comprare password and user password
	if database.CheckPassword(data.Username, data.Password) {
		if err := database.ChangePassword(data.Username, data.NewPassword); err != nil{
			utils.Logger.Error(err.Error())
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		}else{
			utils.CreateHttpResponse(w, http.StatusNoContent, "Password successfuly updated")
			return
		}

	} else {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}
