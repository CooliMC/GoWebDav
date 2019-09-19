package gowebdav

import (
	"../utils"
	"database/sql"
)

func setupDatabase() error {
	//Create a new Connection the given DB and test it
	dbClient, err := utils.MySQLClient(sqlAddress, sqlPort, sqlUsername, sqlPassword)

	//Check if an error occurred and return
	if err != nil {
		return err
	}

	//Check if the Database exist/create
	// and tell the client to use it
	if err := checkDatabase(dbClient); err != nil {
		return err
	}

	//Check if the Tables exists/create
	if err := setupTables(dbClient); err != nil {
		return err
	}

	//Return nil-Error
	return nil
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
	_, err1 := dbClient.Exec(
		"CREATE TABLE IF NOT EXISTS `User` (" +
			"Userid int(11) NOT NULL AUTO_INCREMENT," +
			"Username varchar(50) NOT NULL," +
			"Password varchar(100) NOT NULL," +
			"PRIMARY KEY (Userid)," +
			"UNIQUE KEY Username (Username)" +
		");",
	)

	if err1 != nil {
		return err1
	}

	//Check and setup the GroupDB
	_, err2 := dbClient.Exec(
		"CREATE TABLE IF NOT EXISTS `Group` (" +
			"Groupid int(11) NOT NULL AUTO_INCREMENT," +
			"Groupname varchar(100) NOT NULL," +
			"PRIMARY KEY (Groupid)," +
			"UNIQUE KEY Groupname (Groupname)" +
		");",
	)

	if err2 != nil {
		return err2
	}

	//Check and setup the UserGroupDB
	_, err3 := dbClient.Exec(
		"CREATE TABLE IF NOT EXISTS `FK_UserGroup` (" +
			"Userid int(11) NOT NULL," +
			"Groupid int(11) NOT NULL," +
			"PRIMARY KEY (Userid,Groupid)" +
		");",
	)

	if err3 != nil {
		return err3
	}

	//Check and setup the UserFolderDB
	_, err4 := dbClient.Exec(
		"CREATE TABLE IF NOT EXISTS `FK_UserFolder` (" +
			"Userid int(11) NOT NULL," +
			"Folderpath varchar(255) NOT NULL," +
			"PRIMARY KEY (Userid,Folderpath)" +
		");",
	)

	if err4 != nil {
		return err4
	}

	//Check and setup the GroupFolderDB
	_, err5 := dbClient.Exec(
		"CREATE TABLE IF NOT EXISTS `FK_GroupFolder` (" +
			"Groupid int(11) NOT NULL," +
			"Folderpath varchar(255) NOT NULL," +
			"PRIMARY KEY (Groupid,Folderpath)" +
		");",
	)

	//Return the result of the operation
	return err5
}