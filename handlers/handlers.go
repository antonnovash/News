package handlers

import (
	"encoding/xml"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"news-reader/entity"
	"news-reader/errors"
)

type ctrl interface {
	Result() ([]entity.Item, error)
	Take(id string) (string, error)
}
type Server struct {
	Controller ctrl
}

func (s Server) HandleResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Controller.Result()
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			log.Fatal(errors.TemplateWritingError)
		}
		err = tmpl.Execute(w, items)
		if err != nil {
			log.Fatal(errors.TemplateWritingError)
		}
	}
}
func (s Server) HandleTake() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		news, err := s.Controller.Take("1")
		if err != nil {
			log.Fatal(errors.TemplateWritingError)
		}
		_, _ = w.Write([]byte(news)) //TODO
	}
}
func GetHttpResponse() (*http.Response, error) {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		return nil, errors.WrapError("GetHttpResponse", errors.HttpGetError, err)
	}
	return resp, nil
}

//DecodeRss decode rss date in necessary format
func DecodeRss() (*entity.Rss, error) {
	resp, err := GetHttpResponse()
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(errors.CloseError, err)
		}
	}()
	if err != nil {
		return nil, errors.WrapError("DecodeRss", "er", err) //TODO Message
	}
	rss := &entity.Rss{}
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
	r.HandleFunc("/take", s.HandleTake())
	return r
}
