package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func MySQLClient(
	nInetAddress string,nInetPort int,
	nUsername string, nPassword string,
) (*sql.DB, error) {
	return sql.Open("mysql",
		nUsername + ":" + nPassword +
		"@tcp(" + nInetAddress + ":" +
		strconv.Itoa(nInetPort) + ")/",
	)
}

func MySQLClientDB(
	nInetAddress string,nInetPort int, nUsername string,
	nPassword string, nDatabase string,
) (*sql.DB, error) {
	return sql.Open("mysql",
		nUsername + ":" + nPassword +
		"@tcp(" + nInetAddress + ":" +
		strconv.Itoa(nInetPort) + ")/" +
		nDatabase,
	)
}