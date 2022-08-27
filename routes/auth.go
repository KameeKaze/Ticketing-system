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

	//decode body data
	loginData := &types.Login{}
	json.NewDecoder(r.Body).Decode(&loginData)

	// check if request was valid
	if utils.ValidateJSON(loginData) {
		createHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// comprare password and user password
	user, err := db.Mysql.GetUser(db.Mysql.GetUserId(loginData.Username))
	if user.Name == "" {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if utils.ComparePassword(user.Password, loginData.Password) {
		//generate cookie
		cookie := utils.GenerateSessionCookie()
		//get userId
		userId := db.Mysql.GetUserId(loginData.Username)

		// update or create session based on user already has a session
		err = db.Redis.SetCookie(userId, cookie.Value, &cookie.Expires)
		// check error creating new session
		if err != nil {
			createHttpResponse(w, http.StatusInternalServerError, "Database error")
			utils.Logger.Error(err.Error())
			return
		} else {
			http.SetCookie(w, cookie)
			createHttpResponse(w, http.StatusOK, "Logging in "+loginData.Username)
		}

	} else {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		createHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	//delete cookie
	if err := db.Redis.DeleteCookie(cookie.Value); err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	createHttpResponse(w, http.StatusResetContent, "Logging out")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		createHttpResponse(w, http.StatusUnauthorized, "Login as admin to create user")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	user, err := db.Mysql.GetUser(userId)
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if user.Role != "admin" {
		createHttpResponse(w, http.StatusUnauthorized, "Only admins can create users")
		return
	}
	//decode body data
	registerData := &types.Register{}
	json.NewDecoder(r.Body).Decode(&registerData)

	if utils.ValidateJSON(registerData) {
		createHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	//check if user exist
	if db.Mysql.GetUserId(registerData.Username) != "" {
		createHttpResponse(w, http.StatusConflict, "Username already taken")
		return
	} else {
		err = db.Mysql.AddUser(registerData)
		if err != nil {
			utils.Logger.Error(err.Error())
			createHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		} else {
			createHttpResponse(w, http.StatusCreated, "Creating user "+registerData.Username)
			utils.Logger.Info("Creating user " + registerData.Username)
			return
		}
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		createHttpResponse(w, http.StatusUnauthorized, "Login as admin to create user")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	//decode body data
	data := &types.ChangePassword{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		createHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// comprare password and user password
	user, err := db.Mysql.GetUser(db.Mysql.GetUserId(data.Username))
	if err != nil {
		createHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if utils.ComparePassword(user.Password, data.Password) {
		if err := db.Mysql.ChangePassword(data.Username, data.NewPassword); err != nil {
			utils.Logger.Error(err.Error())
			createHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		} else {
			createHttpResponse(w, http.StatusNoContent, "Password successfuly updated")
			utils.Logger.Info(user.Name + " password update")
			return
		}

	} else {
		createHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}
