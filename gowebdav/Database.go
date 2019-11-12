package gowebdav

import (
	"../utils"
	"database/sql"
)

//All Functions / Structs / Consts for the DB_Connection itself
const (
	crtUserTable = "CREATE TABLE IF NOT EXISTS `User` (" +
		"Userid int(11) NOT NULL AUTO_INCREMENT," +
		"Username varchar(50) NOT NULL," +
		"Password varchar(100) NOT NULL," +
		"PRIMARY KEY (Userid)," +
		"UNIQUE KEY Username (Username)" +
	");"

	crtGroupTable = "CREATE TABLE IF NOT EXISTS `Group` (" +
		"Groupid int(11) NOT NULL AUTO_INCREMENT," +
		"Groupname varchar(100) NOT NULL," +
		"PRIMARY KEY (Groupid)," +
		"UNIQUE KEY Groupname (Groupname)" +
	");"

	crtFKUserGroupTable = "CREATE TABLE IF NOT EXISTS `FK_UserGroup` (" +
		"Userid int(11) NOT NULL," +
		"Groupid int(11) NOT NULL," +
		"PRIMARY KEY (Userid,Groupid)" +
	");"

	crtFKUserFolderTable = "CREATE TABLE IF NOT EXISTS `FK_UserFolder` (" +
		"Userid int(11) NOT NULL," +
		"Folderpath varchar(255) NOT NULL," +
		"PRIMARY KEY (Userid,Folderpath)" +
	");"

	crtFKGroupFolderTable = "CREATE TABLE IF NOT EXISTS `FK_GroupFolder` (" +
		"Groupid int(11) NOT NULL," +
		"Folderpath varchar(255) NOT NULL," +
		"PRIMARY KEY (Groupid,Folderpath)" +
	");"

	getPasswordByUsername = "SELECT Password FROM `User` WHERE Username=?"
)

//Easy to use DB-Connection
type DatabaseConnection struct {
	dbClient *sql.DB
	stmtGetPwdByUser *sql.Stmt
}

func checkDatabase(dbClient *sql.DB) error {
	//Check and setup the Database and return occurred errors
	if _, err := dbClient.Exec("CREATE DATABASE IF NOT EXISTS `" + sqlDatabase + "`;"); err != nil {
		return err
	}

	//Call the MySQLClient to use the Database and return occurred errors
	if _, err := dbClient.Exec("USE `" + sqlDatabase + "`;"); err != nil {
		return err
	}

	//No error occurred
	return nil
}

func setupTables(dbClient *sql.DB) error {
	//Check and setup the UserDB
	if _, err := dbClient.Exec(crtUserTable); err != nil {
		return err
	}

	//Check and setup the GroupDB
	if _, err := dbClient.Exec(crtGroupTable); err != nil {
		return err
	}

	//Check and setup the UserGroupDB
	if _, err := dbClient.Exec(crtFKUserGroupTable); err != nil {
		return err
	}

	//Check and setup the UserFolderDB
	if _, err := dbClient.Exec(crtFKUserFolderTable); err != nil {
		return err
	}

	//Return the result of check and setup the GroupFolderDB
	_, err := dbClient.Exec(crtFKGroupFolderTable); return err
}

func createDatabaseConnection(dbClient *sql.DB) *DatabaseConnection {
	//Create the getPasswordByUser Statement
	stmt1, err1 := dbClient.Prepare(getPasswordByUsername)

	//Check if all prepared statements are ready
	if err1 != nil {
		return nil
	}

	//Return the prepared DatabaseConnection
	return &DatabaseConnection { dbClient,  stmt1 }
}

//All generic function for quick Setup and Check
func setupDatabase() (*DatabaseConnection, error) {
	//Create a new Connection the given DB and test it
	dbClient, err := utils.MySQLClient(sqlAddress, sqlPort, sqlUsername, sqlPassword)

	//Check if an error occurred and return
	if err != nil {
		return nil, err
	}

	//Check if the Database exist/create
	// and tell the client to use it
	if err := checkDatabase(dbClient); err != nil {
		return nil, err
	}

	//Check if the Tables exists/create
	if err := setupTables(dbClient); err != nil {
		return nil, err
	}

	//Return nil-Error
	return createDatabaseConnection(dbClient), nil
}

func (dbCon DatabaseConnection) getUserPassword(username string) (string, error) {
	//Create a local variable
	var tempPassword string

	//Get the row and then scan it for the given parameter and return the result
	return tempPassword, dbCon.stmtGetPwdByUser.QueryRow(username).Scan(&tempPassword)
}