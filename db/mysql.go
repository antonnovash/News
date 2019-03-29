package db

import (
	"database/sql"
	"fmt"
	"log"
	"news-reader/cmd"
	"news-reader/errors"
)

type MySQL struct {
	database *sql.DB
}

// DatabaseInsert insert items(sport news) in database
func (m *MySQL) DatabaseInsert(rss *cmd.Rss) error {
	for _, item := range rss.Channel.Items {
		stmt, err := m.database.Prepare(InsertIntoDatabase)
		if err != nil {
			return errors.WrapError("DatabaseInsert", cmd.DatabasePrepareError, err)
		}
		_, err = stmt.Exec(item.Title, item.PublishedDate)
		if err != nil {
			return errors.WrapError("DatabaseInsert", cmd.ExecError, err)
		}
	}
	return nil
}

//OpenMySQLDb open MySql database
func OpenMySQLDb() (*MySQL, error) {
	db, err := sql.Open("mysql", "root:1234567@/News")
	if err != nil {
		return nil, errors.WrapError("OpenMySQLDb", cmd.OpenDatabaseError, err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.WrapError("OpenMySQLDb", cmd.PingError, err)
	}
	return &MySQL{database: db}, nil
}

//CreateDbTable create table using database
func (m *MySQL) CreateDbTable() error {
	_, err := m.database.Exec(createTableStatements)
	if err != nil {
		return errors.WrapError("CreateDbTable", cmd.CreateTableError, err)
	}
	stmt, err := m.database.Prepare(DropTable)
	if err != nil {
		return errors.WrapError("CreateDbTable", cmd.DropDatabaseError, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Fatal(cmd.CloseError, err)
		}
	}()
	_, err = stmt.Exec()
	if err != nil {
		return errors.WrapError("CreateDbTable", cmd.ExecError, err)
	}
	_, err = m.database.Exec(createTableStatements)
	if err != nil {
		return errors.WrapError("CreateDbTable", cmd.CreateTableError, err)
	}
	return nil
}

func (m *MySQL) ScanDb() ([]cmd.Item, error) {
	rows, err := m.database.Query(SelectFromDatabase)
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
	return items, nil
}
//Controller
func (m *MySQL) Result() ([]cmd.Item, error) {
	rows, err := m.database.Query(SelectFromDatabase)
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
	return items, nil
}
