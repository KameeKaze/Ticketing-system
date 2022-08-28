package types

import (
	"os"

	"github.com/KameeKaze/Ticketing-system/utils"
)

var (
	ROLES             = []string{"programmer", "admin"}
	PORT              = envVar("PORT")
	DATABASE_PASSWORD = envVar("DATABASE_PASSWORD")
)

func envVar(name string) (value string) {
	value, err := os.LookupEnv(name)
	if !err {
		utils.Logger.Error("No env var set " + name)
	}
	return
}
