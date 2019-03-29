package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"news-reader/cmd"
	"testing"
)

type Item struct {
	ID          int
	Title       string
	PublishDate string
}
type getItems struct {
	item           Item
	expectedResult string
}

func checkItem(items getItems) string {
	if items.item.ID < 0 {
		return "Invalid user id"
	}
	if items.item.Title == "" {
		return "Field 'Title' is empty"
	}
	return "OK"

}

func TestItem(t *testing.T) {
	p := []getItems{
		{Item{ID: -1, Title: "first", PublishDate: "27.02.2019"}, "Invalid user id"},
		{Item{ID: 1, Title: "", PublishDate: "27.02.2019"}, "Field 'Title' is empty"},
	}

	for _, val := range p {
		if (checkItem(val) != val.expectedResult) {
			t.Error(val.expectedResult)
		}
	}
}

func TestHttpConnection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Server is start")
		if err != nil {
			t.Error("Error with write to NewServer")
		}
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error("Expected success http request ")
	}
	greeting, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			t.Error(cmd.CloseError, err)
		}
	}()
	fmt.Printf("%s", greeting)
}
func TestUserHandler(t *testing.T) {
	item := []struct {
		ID          int
		Title       string
		PublishDate string
	}{
		{
			ID:          1,
			Title:       "First News",
			PublishDate: "Wed, 27 Mar 2019 13:16:00 +0300",
		},
		{
			ID:          2,
			Title:       "Second News",
			PublishDate: "Wed, 27 Mar 2019 12:56:00 +0300",
		},
	}
	w := httptest.NewRecorder()
	for _, tc := range item {
		jsonData, _ := json.Marshal(tc)
		w.Write(jsonData)
	}
	for _, tc := range item {
		json.Unmarshal(w.Body.Bytes(), &tc) //TODO
		if tc.ID < 0 {
			t.Errorf("Invalid user id %d ", tc.ID)
		}
	}
}
