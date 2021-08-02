package main

import (
	"context"
	"fmt"
	"github.com/simonhylander/booking/db"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", name.Name ), nil
}

func main() {
	//lambda.Start(HandleRequest)

	var defaultCategories = []db.Category{
		{
			Name: "Klippning",
			Category: []db.Category{
				{
					Name:  "Barn 0 - 8 år",
					Price: 290,
				},
				{
					Name:  "Klippning - 30 min",
					Price: 300,
				},
				{
					Name:  "Klippning - 45 min",
					Price: 450,
				},
				{
					Name:  "Klippning - 60 min",
					Price: 490,
				},
			},
		},
	}

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

	for _, category := range defaultCategories {
		db.Create(session)
	}
}