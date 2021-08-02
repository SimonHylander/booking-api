package company

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type CompanyRepository interface {
	Create(company *Company) error
	FindAll() ([]*Company, error)
}

type CompanyNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *CompanyNeo4jRepository) Create(company *Company) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		err = session.Close()
	}()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "CREATE (n:Company {id: apoc.create.uuid(), name: $name}) RETURN n.id, n.name"
		parameters := map[string]interface{}{
			"name": company.Name,
		}

		_, err := tx.Run(query, parameters)
		return nil, err
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyNeo4jRepository) FindAll() (companies []*Company, err error) {
	session := r.Driver.NewSession(neo4j.SessionConfig{})

	defer func() {
		err = session.Close()
	}()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (e:Company) RETURN e.id as id, e.name AS name", map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect()

		if err != nil {
			return nil, err
		}

		var companies []*Company

		for _, record := range records {
			id, _ := record.Get("id")
			name, _ := record.Get("name")

			companies = append(companies, &Company{
				Id:   id.(string),
				Name: name.(string),
			})
		}

		return companies, nil
	})

	companies = result.([]*Company)

	return companies, err
}