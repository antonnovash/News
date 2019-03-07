package main

import (
	"database/sql"
	"fmt"
	"log"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS library DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE library;`,
	`CREATE TABLE IF NOT EXISTS news (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		title TEXT NULL,
		author VARCHAR(255) NULL,
		publishedDate VARCHAR(255) NULL,
		description TEXT NULL,
		PRIMARY KEY (id)
	)`,
}

// mysqlDB persists books to a MySQL instance.
type mysqlDB struct {
	conn *sql.DB

	list   *sql.Stmt
	listBy *sql.Stmt
	insert *sql.Stmt
	get    *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
}
type rowScanner interface {
	Scan(dest ...interface{}) error
}

var NewsDatabase = &mysqlDB{}

func newMySQLDb() (*mysqlDB, error) {
	conn, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal("dbnews:news.go.Eror with open database : ", err)
	}
	db := &mysqlDB{
		conn: conn,
	}

	// Prepared statements. The actual SQL queries are in the code near the
	// relevant method (e.g. addBook).
	if db.list, err = conn.Prepare(listStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare list: %v", err)
	}
	/*
	if db.listBy, err = conn.Prepare(listByStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare listBy: %v", err)
	}
	if db.get, err = conn.Prepare(getStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare get: %v", err)
	}
	if db.insert, err = conn.Prepare(insertStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}
	if db.update, err = conn.Prepare(updateStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare update: %v", err)
	}
	if db.delete, err = conn.Prepare(deleteStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare delete: %v", err)
	}*/

	return db, nil
}

// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			log.Fatal("createTable")
		}
	}
	return nil
}
func scanNews(s rowScanner) (*Item, error) {
	var (
		id            int64
		title         sql.NullString
		author        sql.NullString
		publishedDate sql.NullString
		description   sql.NullString
	)
	if err := s.Scan(&id, &title, &author, &publishedDate, &description); err != nil {
		return nil, err
	}
	news := &Item{
		ID:            id,
		Title:         title.String,
		PublishedDate: publishedDate.String,
		Description:   description.String,
	}
	return news, nil
}

const listStatement = `SELECT * FROM books ORDER BY title`

func (db *mysqlDB) ListBooks() ([]*Item, error) {
	rows, err := db.list.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Item
	for rows.Next() {
		book, err := scanNews(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		books = append(books, book)
	}

	return books, nil
}
