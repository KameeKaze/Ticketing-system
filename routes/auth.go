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
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// comprare password and user password
	user, err := db.Mysql.GetUser(db.Mysql.GetUserId(loginData.Username))
	if user.Name == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
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

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "No sesion cookie specified")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	//delete cookie
	if err := db.Redis.DeleteCookie(cookie.Value); err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	utils.CreateHttpResponse(w, http.StatusResetContent, "Logging out")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session")
	// no sesion cookie set
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Login as admin to create user")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	user, err := db.Mysql.GetUser(userId)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if user.Role != "admin" {
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
	if db.Mysql.GetUserId(registerData.Username) != "" {
		utils.CreateHttpResponse(w, http.StatusConflict, "Username already taken")
		return
	} else {
		err = db.Mysql.AddUser(registerData)
		if err != nil {
			utils.Logger.Error(err.Error())
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		} else {
			utils.CreateHttpResponse(w, http.StatusCreated, "Creating user "+registerData.Username)
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
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Login as admin to create user")
		return
	}

	//check if session exist
	userId, err := db.Redis.GetUserId(cookie.Value)
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if userId == "" {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	//decode body data
	data := &types.ChangePassword{}
	json.NewDecoder(r.Body).Decode(&data)

	// check if request was valid
	if utils.ValidateJSON(data) {
		utils.CreateHttpResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// comprare password and user password
	user, err := db.Mysql.GetUser(db.Mysql.GetUserId(data.Username))
	if err != nil {
		utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
		utils.Logger.Error(err.Error())
		return
	}
	if utils.ComparePassword(user.Password, data.Password) {
		if err := db.Mysql.ChangePassword(data.Username, data.NewPassword); err != nil {
			utils.Logger.Error(err.Error())
			utils.CreateHttpResponse(w, http.StatusInternalServerError, "Database error")
			return
		} else {
			utils.CreateHttpResponse(w, http.StatusNoContent, "Password successfuly updated")
			utils.Logger.Info(user.Name + " password update")
			return
		}

	} else {
		utils.CreateHttpResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
}
