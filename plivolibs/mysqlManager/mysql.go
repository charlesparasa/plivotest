package mysqlManager

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var _db *sql.DB

func initiateSQL() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	var openErr error
	_db, openErr = sql.Open("postgres", dsn)
	if openErr != nil {
		err := fmt.Errorf("error initializing mysql DB connection %s; err: %v", openErr)
		panic(err.Error())
	}

	pingErr := _db.Ping()
	if pingErr != nil {
		err := fmt.Errorf("ping test failed on mysql DB connection %s:%s@%s/%s; err:%v", dsn, pingErr)
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
