package mysqlManager

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/charlesparasa/plivotest/plivolibs/mysqlManager/config"
	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

var (
	TableNameClientCredential = "plivocredentials"
)

func initiateSQL() (*sql.DB, error) {

	config.Init()
	username := config.AppConfig.Mysql.Username
	password := config.AppConfig.Mysql.Password
	host := config.AppConfig.Mysql.Host
	databaseName := config.AppConfig.Mysql.DatabaseName
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, databaseName)
	var openErr error
	_db, openErr = sql.Open("mysql", dsn)
	if openErr != nil {
		err := fmt.Errorf("error initializing mysql DB connection %s; err: %v", host, openErr)
		panic(err.Error())
	}

	pingErr := _db.Ping()
	if pingErr != nil {
		maskedPassword := "*********"
		if len(password) >= 5 {
			maskedPassword = fmt.Sprintf("%s***%s****", password[0:1], password[4:5])
		}
		err := fmt.Errorf("ping test failed on mysql DB connection %s:%s@%s/%s; err:%v", username,
			maskedPassword, host, databaseName, pingErr)
		panic(err.Error())
		return nil, err
	}
	return _db, nil
}

//GetConnection get the mysql DB connection
func GetConnection() (*sql.DB, error) {
	_db, err := initiateSQL()
	if _db == nil {
		err = errors.New("GetConnection: mysql connection is not initialised")
		return nil, err
	}
	return _db, nil
}

//Close close the mssql connection
func Close() {
	_db.Close()
}
