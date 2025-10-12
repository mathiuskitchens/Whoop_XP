package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Conect opens a connection to the MySQL database
func Connect() *sql.DB {
	dsn := "root:secret@tcp(127.0.0.1:3306)/whoop_xp?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Connected to MySQL!")
	return db
}
