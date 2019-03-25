package db

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"news-reader/cmd"
	"news-reader/errors"
)

var database *sql.DB

func DatabaseInsert() error {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		return errors.WrapError("DatabaseInsert", httpGetError, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(CloseError, err)
		}
	}()
	rss := cmd.Rss{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return errors.WrapError("DatabaseInsert", DecodeError, err)
	}
	database, err = OpenMySQLDb()
	if err != nil {
		return errors.WrapError("DatabaseInsert", OpenDatabaseError, err)
	}
	err = CreateDbTable(database)
	if err != nil {
		return errors.WrapError("DatabaseInsert", CreateTableError, err)
	}
	for _, item := range rss.Channel.Items {
		stmt, err := database.Prepare(InsertIntoDatabase)
		if err != nil {
			return errors.WrapError("DatabaseInsert", DatabasePrepareError, err)
		}
		_, err = stmt.Exec(item.Title, item.PublishedDate)
		if err != nil {
			return errors.WrapError("DatabaseInsert", ExecError, err)
		}
	}
	return nil
}

//Open MySql database
func OpenMySQLDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:1234567@/News")
	if err != nil {
		return nil, errors.WrapError("OpenMySQLDb", OpenDatabaseError, err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.WrapError("OpenMySQLDb", PingError, err)
	}
	return db, nil
}

// Create MySql table
func CreateDbTable(db *sql.DB) error {
	stmt, err := db.Prepare(DropTable)
	if err != nil {
		return errors.WrapError("CreateDbTable", DropDatabaseError, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Fatal(CloseError, err)
		}
	}()
	_, err = stmt.Exec()
	if err != nil {
		return errors.WrapError("CreateDbTable", ExecError, err)
	}
	_, err = db.Exec(createTableStatements)
	if err != nil {
		return errors.WrapError("CreateDbTable", CreateTableError, err)
	}
	return nil
}

//Write news to server and read from database
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query(SelectFromDatabase)
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
		log.Fatal(TemplateWritingError)
	}
}

/*func CheckError(msg string, err error) {
	if err != nil {
		log.Fatal(msg)
	}
}*/
