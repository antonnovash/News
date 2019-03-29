package handlers

import (
	"encoding/xml"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"news-reader/cmd"
	"news-reader/db"
	"news-reader/errors"
)

type ctrl interface {
	Result() ([]cmd.Item, error)
}
type Server struct {
	Controller ctrl
}

func (s Server) HandleResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Controller.Result()
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			log.Fatal(cmd.TemplateWritingError)
		}
		err = tmpl.Execute(w, items)
		if err != nil {
			log.Fatal(cmd.TemplateWritingError)
		}
	}
}
func GetHttpResponse() (*http.Response, error) {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		return nil, errors.WrapError("GetHttpResponse", cmd.HttpGetError, err)
	}
	return resp, nil
}
func IndexHandler(m *db.MySQL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := m.ScanDb()
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			log.Fatal(cmd.TemplateWritingError)
		}
		err = tmpl.Execute(w, items)
		if err != nil {
			log.Fatal(cmd.TemplateWritingError)
		}
	}
}

//DecodeRss decode rss date in necessary format
func DecodeRss() (*cmd.Rss, error) {
	resp, err := GetHttpResponse()
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(cmd.CloseError, err)
		}
	}()
	if err != nil {
		return nil, errors.WrapError("DecodeRss", "er", err) //TODO Message
	}
	rss := &cmd.Rss{}
	decoder := xml.NewDecoder(resp.Body);
	err = decoder.Decode(rss)
	if err != nil {
		return nil, errors.WrapError("DecodeRss", "er", err) //TODO Message
	}
	return rss, nil
}
func NewRouter(s Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/result", s.HandleResult())
	return r
}
