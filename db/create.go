package db

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"news-reader/cmd"
)

//Open Sql db
func NewMySQLDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:1234567@/News")
	if err != nil {
		log.Fatal("news.go:Error with open database : ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db, nil
}

// createTable creates the table
func CreateDbTable(db *sql.DB) error {
	stmt, err := db.Prepare(DropTable)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createTableStatements)
	if err != nil {
		log.Fatal("Error with create table: %v", err)
	}
	return nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := cmd.Database.Query(SelectFromDatabase)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	items := []cmd.Item{}

	for rows.Next() {
		p := cmd.Item{}
		err := rows.Scan(&p.ID, &p.Title, &p.PublishedDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}

	tmpl, err := template.ParseFiles("template/index.html")
	err = tmpl.Execute(w, items)
	if err != nil {
		log.Fatal("create.go:Error with writing template,execution stops.")
	}
	/*if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var b bytes.Buffer
	tmpl.Execute(&b, items)
	log.Println(b.String())*/
}
func drop(w http.ResponseWriter, r *http.Request) {
	log.Println("drop collection")

	stmt, err := cmd.Database.Prepare(DropTable)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	//CreateTable()
}