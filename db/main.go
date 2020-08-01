package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// useMysqlDriver()
}

func useMysqlDriver() {
	// user@unix(/path/to/socket)/dbname?charset=utf8
	// user:password@tcp(localhost:5555)/dbname?charset=utf8
	// user:password@/dbname
	// user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
	//
	// 如 user:password@tcp(127.0.0.1:3306)/test
	db, err := sql.Open("mysql", "root:root@tcp(192.168.1.104:3307)/micro?charset=utf8")
	checkErr(err)
	defer db.Close()

	insert(db)
}

func insert(db *sql.DB) int64 {
	// insert
	preparedStatement, err := db.Prepare("insert into userinfo set username=?,department=?,created=?")
	checkErr(err)
	res, err := preparedStatement.Exec("xiaoyu", "研发部门", "2012-12-09")
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("insert id = ", id)
	return id
}

func update(db *sql.DB, id int64) {
	// update
	preparedStatement, err := db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res, err := preparedStatement.Exec("xiaoyu_update", id)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("update affected rows = ", affect)
}

func queryAll(db *sql.DB) {
	// query
	rows, err := db.Query("select * from userinfo")
	checkErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Print(uid, " ")
		fmt.Print(username, " ")
		fmt.Print(department, "")
		fmt.Print(created, "\n")
	}
}

func delete(db *sql.DB, id int64) {
	//删除数据
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	checkErr(err)
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("delete affected rows = ", affect)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
