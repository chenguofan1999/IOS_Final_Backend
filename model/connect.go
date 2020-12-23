package model

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//数据库连接信息
const (
	USERNAME = "root"
	PASSWORD = "Tbd_2333"
	DATABASE = "ios"
)

// DB : 全局的数据库对象
var DB *sql.DB

// Connect : 连接到数据库
func Connect() {
	dataSourceName := fmt.Sprintf("%s:%s@/%s?charset=utf8", USERNAME, PASSWORD, DATABASE)

	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
