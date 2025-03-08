package main

import (
	"banco/db"
	"banco/routes"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	r := routes.SetupRouter()
	r.Run(":8080") // Inicia el servidor en el puerto 8080
}
