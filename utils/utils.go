package utils

import(
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
)

func init(){
	Logging()
}


var Logger *zap.Logger

//create logger
func Logging(){
	// the log file
	logFile, _ := os.OpenFile("ticketing.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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
	Logger = zap.New(core,zap.AddCaller())
	
}

func Comparepassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password))
	return err == nil
}