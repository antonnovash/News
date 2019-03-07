package main

import (
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	ID            int64
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Guid          string    `xml:"guid"`
	Enclosure     Enclosure `xml:"enclosure"`
	PublishedDate string    `xml:"pubDate"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

func main() {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		log.Fatal("main.Eror with http.Get: ", err)
	}
	if err != nil {
		log.Fatal("main.Eror with database connection : ", err)
	}
	if err != nil {
		log.Fatal("main.Eror with Prepare database : ", err)
	}
	if err != nil {
		log.Fatal("main.Eror with create database : ", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal("main.Eror with resp.Body.Close: ", err)
		}
	}()
	rss := Rss{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatalf("Error Decode: %v\n", err)
	}
	_, _ = newMySQLDb()
	fmt.Printf("Title: %v\n", rss.Channel.Title)
	fmt.Printf("Description: %v\n", rss.Channel.Desc)
	for i, item := range rss.Channel.Items {
		fmt.Printf("%d.\t Title: %v\n\t Description: %v\n\t PubDate: %v\n", i+1, item.Title, item.Description, item.PublishedDate)
	}
}
