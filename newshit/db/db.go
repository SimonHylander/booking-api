package db

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/simonhylander/booking-api/company"
	"github.com/simonhylander/booking-api/employee"
)

func AutoMigrate(driver neo4j.Driver) error {
	fmt.Println("AutoMigrate")

	err := clearNodesWithRelationships(driver)
	if err != nil {
		panic(err)
	}

	err = clearNodesWithoutRelationships(driver)
	if err != nil {
		panic(err)
	}

	c, err := addCompanies(driver)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)

	err = addEmployees(driver, c)
	if err != nil {
		panic(err)
	}

	err = addSchedule(driver, c)
	if err != nil {
		panic(err)
	}

	return nil
}

func clearNodesWithRelationships(driver neo4j.Driver) error {
	session := driver.NewSession(neo4j.SessionConfig{})

	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a) -[r] -> () DELETE a, r", map[string]interface{}{})

		if err != nil {
			return result, err
		}

		return result, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func clearNodesWithoutRelationships(driver neo4j.Driver) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a) DELETE a", map[string]interface{}{})

		if err != nil {
			return result, err
		}

		return result, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func addCompanies(driver neo4j.Driver) (*company.Company, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return saveCompany(tx)
	})

	if result == nil {
		return nil, err
	}

	c := result.(*company.Company)

	return c, nil
}

func saveCompany(tx neo4j.Transaction) (*company.Company, error) {
	//p := "CREATE p = (e:Employee {id: apoc.create.uuid(), name:'Andy'})-[:EMPLOYED_BY]->(c:Company { id: apoc.create.uuid(), name: 'Highlights2'})<-[:EMPLOYED_BY]-(b:Employee {id: apoc.create.uuid(), name: 'Michael'}) return p"

	result, err := tx.Run("CREATE (n:Company {id: apoc.create.uuid(), name: $name}) RETURN n.id as id, n.name as name", map[string]interface{}{
		"name": "Highlights",
	})

	if result == nil {
		return nil, err
	}

	record, err := result.Single()

	if record == nil {
		return nil, err
	}

	id, _ := record.Get("id")
	name, _ := record.Get("name")

	return &company.Company{
		Id:   id.(string),
		Name: name.(string),
	}, nil
}


func addEmployees(driver neo4j.Driver, company *company.Company) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "CREATE (n:Employee {id: apoc.create.uuid(), firstname: $firstname, lastname: $lastname}) RETURN n.id as id, n.firstname as firstname, n.lastname as lastname"
		result, err := tx.Run(query, map[string]interface{}{
			"firstname": "Sabina",
			"lastname":  "Hylander",
		})

		if err != nil {
			return nil, err
		}

		record, err := result.Single()

		if err != nil {
			return nil, err
		}

		id, _ := record.Get("id")
		firstname, _ := record.Get("firstname")
		lastname, _ := record.Get("lastname")

		return &employee.Employee{
			Id:        id.(string),
			Firstname: firstname.(string),
			Lastname:  lastname.(string),
		}, nil
	})

	if result == nil {
		return err
	}

	employee := result.(*employee.Employee)

	// Company Relationship
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (c:Company), (e:Employee) " +
				 "WHERE c.id = $companyId AND e.id = $employeeId " +
				 "CREATE (e)-[r:EMPLOYED_BY]->(c) " +
				 "RETURN type(r)"

		_, err := tx.Run(query, map[string]interface{}{
			"employeeId": employee.Id,
			"companyId": company.Id,
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	// Title
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (e:Employee {id: $employeeId}) CREATE (t:Title {name: 'VD'})-[r:BELONGS_TO]->(e) return e "

		_, err := tx.Run(query, map[string]interface{}{
			"employeeId": employee.Id,
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return nil
}


//MATCH (c:Company) WHERE c.name = 'Highlights' CREATE (s:Schedule {name: 'Monday', day: 1, isOpen: true, start: '09:00:00', end: '18:00:00'})-[r:BELONGS_TO]->(c) return type(r)
func addSchedule(driver neo4j.Driver, company *company.Company) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (c:Company) " +
				 "WHERE c.id = $companyId " +
				 //"CREATE (s:Schedule {name: 'Monday', day: 1, isOpen: true, start: '09:00:00', end: '18:00:00'})-[r:BELONGS_TO]->(c) " +
				 "CREATE (s:Schedule {day: 1, name: 'Monday', isOpen: true, start: '09:00:00', end: '18:00:00'})-[r:BELONGS_TO]->(c), " +
				 "(s2:Schedule {day: 2, name: 'Tueday', isOpen: true, start: '09:00:00', end: '18:00:00'})-[r2:BELONGS_TO]->(c)," +
				 "(s3:Schedule {day: 3, name: 'Wednesday', isOpen: true, start: '09:00:00', end: '18:00:00'})-[r3:BELONGS_TO]->(c)," +
				 "(s4:Schedule {day: 4, name: 'Thursday', isOpen: true, start: '09:00:00', end: '18:00:00'})-[r4:BELONGS_TO]->(c)," +
				 "(s5:Schedule {day: 5, name: 'Thursday', isOpen: true, start: '09:00:00', end: '18:00:00'})-[r5:BELONGS_TO]->(c)," +
				 "(s6:Schedule {day: 6, name: 'Saturday', isOpen: false})-[r6:BELONGS_TO]->(c)," +
				 "(s7:Schedule {day: 7, name: 'Sunday', isOpen: false})-[r7:BELONGS_TO]->(c)" +
				 "RETURN c"

		_, err := tx.Run(query, map[string]interface{}{
			"companyId": company.Id,
			"day1Name": "",
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

/*func saveSchedule(tx neo4j.Transaction) (*company.Company, error) {
	//p := "CREATE p = (e:Employee {id: apoc.create.uuid(), name:'Andy'})-[:EMPLOYED_BY]->(c:Company { id: apoc.create.uuid(), name: 'Highlights2'})<-[:EMPLOYED_BY]-(b:Employee {id: apoc.create.uuid(), name: 'Michael'}) return p"

	result, err := tx.Run("CREATE (n:Company {id: apoc.create.uuid(), name: $name}) RETURN n.id as id, n.name as name", map[string]interface{}{
		"name": "Highlights",
	})

	if result == nil {
		return nil, err
	}

	record, err := result.Single()

	if record == nil {
		return nil, err
	}

	id, _ := record.Get("id")
	name, _ := record.Get("name")

	return &company.Company{
		Id:   id.(string),
		Name: name.(string),
	}, nil
}*/