package main

import (
	"fmt"
	"github.com/simonhylander/booking/db"
)

type EmployeeRepository interface {
	FindByEmailAndPassword(email string, password string) (*User, error)
}

func Connect () {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "booking"
	str, err := db.HelloWorld(uri, username, password, false)
	fmt.Println(str)
	fmt.Println(err)

	session, err := db.Connect(uri, username, password, false)

	if err != nil {
		fmt.Println(err)
		return
	}
}