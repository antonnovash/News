package db

var (
	/*createTableStatements = `CREATE TABLE IF NOT EXISTS sportNews (
		id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		title TEXT NULL,
		publishedDate VARCHAR(255) NULL
	)`*/
	InsertIntoDatabase = `INSERT INTO sportNews (
		title, publishedDate
		) VALUES (?, ?)`
	CleanTable = `TRUNCATE TABLE sportNews;`

	SelectFromDatabase = `SELECT * FROM sportNews`

	GetNewsByID = `SELECT title FROM sportNews WHERE id=?`
)
