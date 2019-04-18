package db

import (
	"database/sql"
	"fmt"
	"log"
	"news-reader/entity"
	"news-reader/errors"
)

type MySQL struct {
	database *sql.DB
}

// DatabaseInsert insert items(sport news) in database
func (m *MySQL) DatabaseInsert(rss *entity.Rss) error {
	for _, item := range rss.Channel.Items {
		stmt, err := m.database.Prepare(InsertIntoDatabase)
		if err != nil {
			return errors.WrapError("DatabaseInsert", errors.DatabasePrepareError, err)
		}
		_, err = stmt.Exec(item.Title, item.PublishedDate)
		if err != nil {
			return errors.WrapError("DatabaseInsert", errors.ExecError, err)
		}
	}
	return nil
}

//OpenMySQLDb open MySql database
func OpenMySQLDb() (*MySQL, error) {
	db, err := sql.Open("mysql", "root:1234567@/News")
	if err != nil {
		return nil, errors.WrapError("OpenMySQLDb", errors.OpenDatabaseError, err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.WrapError("OpenMySQLDb", errors.PingError, err)
	}
	return &MySQL{database: db}, nil
}

//CleanDbTable create table using database
func (m *MySQL) CleanDbTable() error {
	stmt, err := m.database.Prepare(CleanTable)
	if err != nil {
		return errors.WrapError("CreateDbTable", errors.DropDatabaseError, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Fatal(errors.CloseError, err)
		}
	}()
	_, err = stmt.Exec()
	if err != nil {
		return errors.WrapError("CreateDbTable", errors.ExecError, err)
	}
	_, err = m.database.Exec(createTableStatements)
	if err != nil {
		return errors.WrapError("CreateDbTable", errors.CreateTableError, err)
	}
	return nil
}

func (m *MySQL) ScanDb() ([]entity.Item, error) {
	rows, err := m.database.Query(SelectFromDatabase)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	items := []entity.Item{}

	for rows.Next() {
		p := entity.Item{}
		err := rows.Scan(&p.ID, &p.Title, &p.PublishedDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items, nil
}

//Controller
func (m *MySQL) Result() ([]entity.Item, error) {
	rows, err := m.database.Query(SelectFromDatabase)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	items := []entity.Item{}

	for rows.Next() {
		p := entity.Item{}
		err := rows.Scan(&p.ID, &p.Title, &p.PublishedDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items, nil
}

func (m *MySQL) Take(id string) (string, error) {
	var newsById string
	err := m.database.QueryRow(GetNewsByID,id).Scan(&newsById)
	if err!= nil {
		return "", errors.WrapError("Take", errors.GetNewsByIdError, err)
	}
	return newsById, nil
}
