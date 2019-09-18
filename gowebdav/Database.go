package gowebdav

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type mysqlClient struct {
	Connection sql.DB
}

func MySQLClient(nInetAddress string, nInetPort int, nUsername string, nPassword string, nDatabase string) (*mysqlClient, error) {

	db, err := sql.Open("mysql",
		nUsername + ":" +
		nPassword + "@tcp(" +
		nInetAddress + ":" +
		strconv.Itoa(nInetPort) + ")/" +
		nDatabase,
	)

	if err != nil {
		return nil, err
	} else {
		return &mysqlClient{ Connection: *db }, nil
	}
}

func (e mysqlClient) Close() error {
	return e.Connection.Close()
}
