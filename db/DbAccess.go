package db

import (
	"database/sql"
)

var db *sql.DB

func Open() error {
	var err error
	//参数根据自己的数据库进行修改
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=test123 dbname=test sslmode=disable")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
	//var rows *Rows
	//rows, err = db.Query("select * from t_user_info where f_account=$1", "admin")

	//defer rows.close()

}

func SelectById() error {
	return nil
}
