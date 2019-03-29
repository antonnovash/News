package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"news-reader/cmd"
	"news-reader/db"
	"news-reader/errors"
	"news-reader/handlers"
	"time"
)

func main() {
	t, err := GetMySQL()
	if err != nil {
		log.Fatal("Error with Get MySql: ", err)
	}
	server := handlers.Server{Controller: t}
	r := handlers.NewRouter(server)
	s := http.Server{
		Addr:         ":3006",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      r,
	}
	s.ListenAndServe()
	/*h := http.NewServeMux()
	h.HandleFunc("/", handlers.IndexHandler(t))
	http.ListenAndServe(":3006", h)*/
}

func GetMySQL() (*db.MySQL, error) { //TODO
	rss, err := handlers.DecodeRss()
	m, err := db.OpenMySQLDb()
	err = m.CreateDbTable()
	err = m.DatabaseInsert(rss)
	if err != nil {
		return nil, errors.WrapError("getMySQL", cmd.OpenDatabaseError, err)
	}
	return m, nil
}
