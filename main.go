package main

import (
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"news-reader/cmd"
	"news-reader/db"
)
func main() {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		log.Fatal("main.Error with http.Get: ", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal("main.Error with resp.Body.Close: ", err)
		}
	}()
	rss := cmd.Rss{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)

	cmd.Database, err = db.NewMySQLDb() //TODO
	_ = db.CreateDbTable(cmd.Database)    //TODO
	for _, item := range rss.Channel.Items {
		stmt, err := cmd.Database.Prepare(db.InsertIntoDatabase)
		if err != nil {
			log.Fatal("main.go:Error with database.Prepare: ", err)
		}
		_, err = stmt.Exec(item.Title, item.PublishedDate)
		if err != nil {
			log.Fatal("main.go:Error with Exec data: ", err)
		}
	}
	http.HandleFunc("/", db.IndexHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":3006", http.DefaultServeMux)

}

