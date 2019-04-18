package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
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
	_ = s.ListenAndServe()//TODO
}

func GetMySQL() (*db.MySQL, error) { //TODO
	rss, err := handlers.DecodeRss()
	m, err := db.OpenMySQLDb()
	if err != nil {
		log.Fatal(err)
	}
	err = m.CreateDbTable()
	err = m.DatabaseInsert(rss)
	if err != nil {
		return nil, errors.WrapError("GetMySQL", errors.OpenDatabaseError, err)
	}
	return m, nil
}
