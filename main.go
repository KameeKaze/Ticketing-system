package main

import (
	"github.com/KameeKaze/Ticketing-system/db"
	"github.com/KameeKaze/Ticketing-system/routes"
)

func main() {
	defer db.Mysql.Close()
	//start
	routes.RoutesHandler()
}
