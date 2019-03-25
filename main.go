package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"news-reader/db"
)

func main() {
	err := db.DatabaseInsert()
	if err != nil {
		log.Fatal("Error with Insert into database: ", err)
	}
	http.HandleFunc("/", db.IndexHandler)
	http.ListenAndServe(":3006", http.DefaultServeMux)
}
