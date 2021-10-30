package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type User struct {
	name     string
	password string
	address  string
}

func queryUserByName(name string, db *sql.DB) (*User, error) {
	var user User
	err := db.QueryRow("select name, password, address from userInfo where name=%s", name).Scan(&user.name, &user.password, &user.address)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "query failed!")
	} else {
		return &user, nil
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println("open database failed!")
		return
	}

	defer db.Close()

	userInfo, err := queryUserByName("lily", db)
	if err != nil {
		fmt.Printf("err code : %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace : \n%+v\n", err)
		return
	}

	if userInfo != nil {
		fmt.Printf("get user info, name:%s, address:%s\n", userInfo.name, userInfo.address)
	} else {
		fmt.Printf("user(%s) not exist, please register first!\n", userInfo.name)
	}

	return
}
