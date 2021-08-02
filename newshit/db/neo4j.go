package db

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Category struct {
	Name     string
	Price    int
	Category []Category
}

/*type Type struct {
	Name string `json:"name"`
}*/

func HelloWorld(uri, username, password string, encrypted bool) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		return "", err
	}

	defer driver.Close()

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var defaultCategories = []Category{
		{
			Name: "Klippning",
			Category: []Category{
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

	for _, category := range defaultCategories {
		dbCategory, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run("CREATE (c:Category {Name: $name})", map[string]interface{}{"name": category.Name})

			if err != nil {
				return nil, err
			}

			if result.Next() {
				return result.Record().GetByIndex(0), nil
			}

			return nil, result.Err()
		})

		fmt.Println(dbCategory)

		if err != nil {
			fmt.Println(err)
		}
	}

	// CREATE (c:Category) SET c.name
	// CREATE (m:Menu) SET m.category

	/*if err != nil {
		return "", err
	}

	return greeting.(string), nil*/

	return "", nil
}

/*func Connect(uri, username, password string, encrypted bool) (neo4j.Session, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = encrypted
	})

	if err != nil {
		return nil, err
	}

	defer driver.Close()

	driver.NewSession(neo4j.SessionConfig{})

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session, nil
}*/