package utils

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/KameeKaze/Ticketing-system/types"
)

func init() {
	Logging()
}

var Logger *zap.Logger

//create logger
func Logging() {
	// the log file
	logFile, _ := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// create config
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // set iso time format

	// create encoder for both terminal and logfile
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)
	Logger = zap.New(core, zap.AddCaller())

}

//hash the user's password
func HashPassword(password string) string {
	// Create a byte slice
	var passwordBytes = []byte(password)
	// Hash password
	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidateJSON(data interface{}) bool {
	validate := validator.New()
	return validate.Struct(data) != nil
}

func CreateHttpResponse(w http.ResponseWriter, statusCode int, text string) {
	//set status code
	w.WriteHeader(statusCode)
	//create json
	r, _ := json.Marshal(types.ResponseBody{
		Msg: text,
	})
	//send data
	w.Write([]byte(r))
}

func GenerateSessionCookie() *http.Cookie {
	// generate http cookie
	cookie := &http.Cookie{
		Name:     "session",
		Value:    uuid.New().String(),
		HttpOnly: true,
		Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
		Path:     "/",
	}
	return cookie
}
