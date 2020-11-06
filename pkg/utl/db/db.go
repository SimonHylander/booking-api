package db

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func New() (neo4j.Session, error) {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "booking"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = true
	})

	if err != nil {
		return nil, err
	}

	defer driver.Close()

	session, err := driver.Session(neo4j.AccessModeWrite)

	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session, nil

	/*greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil*/
}

/*import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	// DB adapter
	_ "github.com/lib/pq"
)

type dbLogger struct{}

// BeforeQuery hooks before pg queries
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery hooks after pg queries
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	fmt.Println(query)
	return err
}

// New creates new database connection to a postgres database
func New(psn string, timeout int, enableLog bool) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db = db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}*/