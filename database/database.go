package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	_"github.com/lib/pq"
	"github.com/joho/godotenv"
)

var DB *sql.DB
func ConnectionDb(){
	err := godotenv.Load()
	if(err!=nil){
		fmt.Println("Couldnt Load the env file")
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	connstr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, pass)
	db, dberror := sql.Open("postgres", connstr)
	if(dberror!=nil){
		fmt.Println("Error while opening a connection to the database", err)
	}else{
		DB=db
		fmt.Println("connected successfully")
	}


}