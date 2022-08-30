package types

import (
	"os"

	"github.com/KameeKaze/Ticketing-system/utils"
)

var (
	ROLES          = []string{"programmer", "admin"}
	PORT           = envVar("PORT")
	MYSQL_PASSWORD = envVar("MYSQL_PASSWORD")
	REDIS_PASSWORD = envVar("REDIS_PASSWORD")
)

func envVar(name string) (value string) {
	value, err := os.LookupEnv(name)
	if !err {
		utils.Logger.Error("No env var set " + name)
	}
	return
}
