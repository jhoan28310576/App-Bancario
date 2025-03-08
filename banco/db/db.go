package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	// Cambia "root:@tcp(127.0.0.1:3306)/Banco" por tus credenciales
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/Banco")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error conectando a la DB:", err)
	}
	fmt.Println("✅ Conexión a DB exitosa!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
