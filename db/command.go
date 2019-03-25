package db

var createTableStatements = `CREATE TABLE IF NOT EXISTS news (
		id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		title TEXT NULL,
		publishedDate VARCHAR(255) NULL
	)`

var InsertIntoDatabase = `INSERT INTO news (
title, publishedDate
) VALUES (?, ?)`

var DropTable = `DROP TABLE news;`

var SelectFromDatabase = `SELECT * FROM news`
