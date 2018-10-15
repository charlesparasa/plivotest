package mysqlManager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/charlesparasa/plivotest/service/model"
)

const plivotable  =  "plivo_contact"
var localdb  *sql.DB

func GetDB() *sql.DB {
	return localdb
}
func init()  {
	ldb, err := GetConnection()
	if err != nil {
		fmt.Println("no connection available")
		return
	}
	localdb = ldb
}
func InsertOne(data interface{})error  {
	query := fmt.Sprintf("insert into %s (id,name,email,phone) values(?,?,?,?)",plivotable)
	stmt , prepErr := localdb.Prepare(query)
	if prepErr != nil {
		err := fmt.Errorf("Invalid mysql statments", prepErr)
		return  err
	}
	var c model.Contact
	bytes, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("Unable to marshall the data +%v", err)
		return err
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		err = fmt.Errorf("unbale to unmarshall data +%v", err)
		return err
	}

	result, execErr := stmt.Exec(c.Id, c.Name,c.Email,c.Phone)
	if execErr != nil {
		err := fmt.Errorf("InsertContact: error inserting %+v; err: %v", execErr)
		return err
	}

	id, lastErr := result.LastInsertId()
	rows, rowErr := result.RowsAffected()
	fmt.Printf("result: %v, %v, \n %v, %v\n", id, lastErr, rows, rowErr)
	defer localdb.Close()
	return  nil
}

func GetData(from string , to string )(data interface{}, err error)  {
	fmt.Println(from ,"from", to)
	if to == "0" {
		to = "10"
	}
	str := "select * from " + plivotable + " LIMIT " + from + " , " + to
	query := fmt.Sprintf(str)
	// Execute the query
	results, err := localdb.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	 idPtr := new(int)
	 namePtr := new(string)
	 emailPtr := new(string)
	 phonePtr := new(string)
	 var c model.Contact
	 var ca []model.Contact
	for results.Next() {
		scanErr := results.Scan(&idPtr, &namePtr, &emailPtr, &phonePtr)
		if scanErr != nil {
			fmt.Println("scanErr",scanErr)
			return nil, err
		}
		c.Id = *idPtr
		c.Name = *namePtr
		c.Email = *emailPtr
		c.Phone = *phonePtr
		ca = append(ca , c)
	}

	 return ca , nil
}

func DeleteContact(email string) error {
	query := fmt.Sprintf("DELETE FROM "+plivotable+ " WHERE email=?")
	stmt , prepErr := localdb.Prepare(query)
	if prepErr != nil {
		err := fmt.Errorf("Invalid mysql statments", prepErr)
		return  err
	}
	result , err := stmt.Exec(email)
	if err != nil {
		fmt.Println("error" , err)
		return err
	}
	rowsEffected, err  := result.RowsAffected()
	if err != nil {
		fmt.Println("error", err)
		return  err
	}
	if rowsEffected == 0 {
		fmt.Println("no rows are effected")
	}
	return nil
}
func UpdateContact(contact model.Contact )error  {
	query := fmt.Sprintf("Update "+plivotable+ " Set  name=?, email=?, phone=? where id=?")
	stmt , prepErr := localdb.Prepare(query)
	if prepErr != nil {
		err := fmt.Errorf("Invalid mysql statments", prepErr)
		return  err
	}
	result, err := stmt.Exec(contact.Name, contact.Email, contact.Phone, contact.Id)
	if err != nil {
		fmt.Println("Contact not updated " , err)
		return err
	}
	rowsEffected, err  := result.RowsAffected()
	if err != nil {
		fmt.Println("error", err)
		return  err
	}
	if rowsEffected == 0 {
		fmt.Println("no rows are effected")
	}

	return nil
}
