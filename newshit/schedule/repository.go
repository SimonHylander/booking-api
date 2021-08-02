package schedule

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ScheduleRepository interface {
	Create(e *Schedule) (*Schedule, error)
	Find() ([]*Schedule, error)
}

type ScheduleNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *ScheduleNeo4jRepository) Create(e *Schedule) (schedule *Schedule, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		err = session.Close()
	}()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "CREATE (s:Schedule {id: apoc.create.uuid(), nameOfDay: $nameOfDay, day: $day, isOpen: $isOpen}) RETURN s.id as id, s.nameOfDay as nameOfDay, s.day as day, s.isOpen as isOpen"
		parameters := map[string]interface{}{
			"nameOfDay": e.NameOfDay,
			"day":  e.Day,
			"isOpen": e.IsOpen,
		}

		result, err := tx.Run(query, parameters)

		if err != nil {
			return nil, err
		}

		record, err := result.Single()

		if err != nil {
			return nil, err
		}

		id, _ := record.Get("id")
		nameOfDay, _ := record.Get("nameOfDay")
		day, _ := record.Get("day")
		isOpen, _ := record.Get("isOpen")

		return &Schedule{
			Id:        id.(string),
			NameOfDay: nameOfDay.(string),
			Day: day.(int),
			IsOpen:  isOpen.(bool),
		}, nil
	})

	if result == nil {
		return nil, err
	}

	schedule = result.(*Schedule)

	// Relationship

	/*_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (a:Company), (b:Schedule) " +
			"WHERE a.id = $companyId AND b.id = $employeeId " +
			"CREATE (a)-[r:EMPLOYED_BY]->(b) " +
			"RETURN type(r)"

		_, err := tx.Run(query, map[string]interface{}{
			"employeeId": employee.Id,
			"companyId": e.CompanyId,
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})*/

	return schedule, nil
}

func (r *ScheduleNeo4jRepository) FindAll() (scheduleDays []*Schedule, err error) {
	session := r.Driver.NewSession(neo4j.SessionConfig{})

	defer func() {
		err = session.Close()
	}()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (e:Schedule)--(c:Company) RETURN e.id as id, e.firstname AS firstname, e.lastname AS lastname, c.name as company"
		result, err := tx.Run(query, map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect()

		if err != nil {
			return nil, err
		}

		var scheduleDays []*Schedule

		for _, record := range records {
			id, _ := record.Get("id")
			nameOfDay, _ := record.Get("nameOfDay")
			day, _ := record.Get("day")
			isOpen, _ := record.Get("isOpen")

			scheduleDays = append(scheduleDays, &Schedule{
				Id:   id.(string),
				NameOfDay: nameOfDay.(string),
				Day: day.(int),
				IsOpen: isOpen.(bool),
			})
		}

		return scheduleDays, nil
	})

	scheduleDays = result.([]*Schedule)

	return scheduleDays, err
}