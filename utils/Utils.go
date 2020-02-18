package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//MySQLClient for easier DB-Access without a special database
func MySQLClient(
	nInetAddress string, nInetPort int,
	nUsername string, nPassword string,
) (*sql.DB, error) {
	return sql.Open("mysql",
		nUsername + ":" + nPassword +
		"@tcp(" + nInetAddress + ":" +
		strconv.Itoa(nInetPort) + ")/",
	)
}

//MySQLClient for easier DB-Access with a special database
func MySQLClientDB(
	nInetAddress string, nInetPort int, nUsername string,
	nPassword string, nDatabase string,
) (*sql.DB, error) {
	return sql.Open("mysql",
		nUsername + ":" + nPassword +
		"@tcp(" + nInetAddress + ":" +
		strconv.Itoa(nInetPort) + ")/" +
		nDatabase,
	)
}