package employee

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strconv"
)

type EmployeeRepository interface {
	Create(e *Employee) (*Employee, error)
	Find() ([]*Employee, error)
	Get(id string) (*Employee, error)
}

type EmployeeNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *EmployeeNeo4jRepository) Create(e *Employee) (employee *Employee, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		err = session.Close()
	}()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "CREATE (n:Employee {id: apoc.create.uuid(), firstname: $firstname, lastname: $lastname}) RETURN n.id as id, n.firstname as firstname, n.lastname as lastname"
		parameters := map[string]interface{}{
			"firstname": e.Firstname,
			"lastname":  e.Lastname,
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
		firstname, _ := record.Get("firstname")
		lastname, _ := record.Get("lastname")
		fullname := firstname.(string) + " " + lastname.(string)

		return &Employee{
			Id:        id.(string),
			Firstname: firstname.(string),
			Lastname:  lastname.(string),
			Fullname:  fullname,
		}, nil
	})

	if result == nil {
		return nil, err
	}

	fmt.Println(result)

	employee = result.(*Employee)

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (a:Company), (b:Employee) " +
			"WHERE a.id = $companyId AND b.id = $employeeId " +
			"CREATE (a)-[r:EMPLOYED_BY]->(b) " +
			"RETURN type(r)"

		_, err := tx.Run(query, map[string]interface{}{
			"employeeId": employee.Id,
			"companyId":  e.CompanyId,
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return employee, nil
}

func (r *EmployeeNeo4jRepository) Find() (employees []*Employee, err error) {
	session := r.Driver.NewSession(neo4j.SessionConfig{})

	defer func() {
		err = session.Close()
	}()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (e:Employee)--(c:Company) RETURN e.id as id, e.firstname AS firstname, e.lastname AS lastname, c.name as company"
		result, err := tx.Run(query, map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect()

		if err != nil {
			return nil, err
		}

		var employees []*Employee

		for _, record := range records {
			id, _ := record.Get("id")
			firstname, _ := record.Get("firstname")
			lastname, _ := record.Get("lastname")
			fullname := firstname.(string) + " " + lastname.(string)
			company, _ := record.Get("company")

			titleQuery := "MATCH (t:Title)-[r:BELONGS_TO]->(e:Employee {id: $employeeId}) RETURN t.name as name"
			titleResult, err := tx.Run(titleQuery, map[string]interface{}{
				"employeeId": id.(string),
			})

			if err != nil {
				return nil, err
			}

			titleRecords, err := titleResult.Collect()

			if err != nil {
				return nil, err
			}

			var titles []string

			for _, titleRecord := range titleRecords {
				title, _ := titleRecord.Get("name")
				titleStr := title.(string)
				titles = append(titles, titleStr)
			}

			/*matchBuilder := match()
			matchBuilder.node("Title", "t")
			rel := matchBuilder.relationship("BELONGS_TO", "r")
			rel.direction("->")
			rel.node = node{
				name:     "Employee",
				variable: "e",
			}
			matchBuilder.buildQuery()
			fmt.Println(matchBuilder.getQuery().query)*/

			employees = append(employees, &Employee{
				Id:        id.(string),
				Firstname: firstname.(string),
				Lastname:  lastname.(string),
				Fullname:  fullname,
				Company:   company.(string),
				Titles: titles,
			})
		}

		return employees, nil
	})

	employees = result.([]*Employee)

	return employees, err
}

func (r *EmployeeNeo4jRepository) Get(id string) (employee *Employee, err error) {
	session := r.Driver.NewSession(neo4j.SessionConfig{})

	defer func() {
		err = session.Close()
	}()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (e:Employee {id: $employeeId})--(c:Company) RETURN e.id as id, e.firstname AS firstname, e.lastname AS lastname, c.name as company"
		result, err := tx.Run(query, map[string]interface{}{
			"employeeId": id,
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
		fullname := firstname.(string) + " " + lastname.(string)
		company, _ := record.Get("company")

		titleQuery := "MATCH (t:Title)-[r:BELONGS_TO]->(e:Employee {id: $employeeId}) RETURN t.name as name"
		titleResult, err := tx.Run(titleQuery, map[string]interface{}{
			"employeeId": id.(string),
		})

		if err != nil {
			return nil, err
		}

		titleRecords, err := titleResult.Collect()

		if err != nil {
			return nil, err
		}

		var titles []string
		for _, titleRecord := range titleRecords {
			title, _ := titleRecord.Get("name")
			titleStr := title.(string)
			titles = append(titles, titleStr)
		}

		var employee = &Employee{
			Id:        id.(string),
			Firstname: firstname.(string),
			Lastname:  lastname.(string),
			Fullname:  fullname,
			Company:   company.(string),
			Titles: titles,
			Image: "avatar1",
		}

		return employee, nil
	})

	employee = result.(*Employee)

	return employee, err
}

type iBuilder interface {
	//match()
	node(string, string)
	withProperties(map[string]interface{})
	relationship(string, string) relationship
	buildQuery()
	getQuery() CypherQuery
}

func match() iBuilder {
	return &matchBuilder{}
}

type matchBuilder struct {
	query        string
	nodeName     string
	nodeVariable string
	//nodeProperties interface{}
	nodeProperties map[string]interface{}
	rel            relationship
}

type relationship struct {
	node      node
	reltype   string
	variable  string
	relDirection string
}

type node struct {
	name     string
	variable string
}

func (b *matchBuilder) node(node string, variable string) {
	b.nodeName = node
	b.nodeVariable = variable
}

func (b *matchBuilder) withProperties(properties map[string]interface{}) {
	b.nodeProperties = properties
}

func (b *matchBuilder) relationship(reltype string, variable string) relationship {
	b.rel = relationship{
		node:      node{
			name: "",
			variable: "",
		},
		reltype:   reltype,
		variable:  variable,
		relDirection: "",
	}

	return b.rel
}

func (r *relationship) direction(relDirection string) {
	r.relDirection = relDirection
}

func (b *matchBuilder) buildQuery() {
	node := b.nodeName
	variable := b.nodeVariable
	b.query = fmt.Sprintf("MATCH (%s:%s", variable, node)

	if len(b.nodeProperties) > 0 {
		b.query += " {"
		for key, value := range b.nodeProperties {
			str := fmt.Sprintf("%v", value)
			if _, err := strconv.Atoi(str); err == nil {
				b.query += fmt.Sprintf("%s: %s", key, value)
			} else {
				b.query += fmt.Sprintf("%s: '%s'", key, value)
			}
		}
		b.query += "}"
	}

	b.query += ")"

	/*if len(b.reltype) > 0 {
		b.query += fmt.Sprintf("-[%s:%s]-", b.relvar, b.reltype)
		if len(b.relvar) > 0 {

		} else {
			b.query += fmt.Sprintf("-[:%s]-", b.reltype)
		}
	}*/

	fmt.Println(b.rel.relDirection)

	rel := b.rel
	b.query += fmt.Sprintf("-[%s:%s]%s", rel.variable, rel.reltype, rel.direction)
	//b.query += fmt.Sprintf("(%s:%s)", rel.node.variable, rel.node.name)
}

func (b *matchBuilder) getQuery() CypherQuery {
	return CypherQuery{
		query: b.query,
	}
}

type House struct {
	windowType string
	doorType   string
	floor      int
}

type CypherQuery struct {
	query string
}

/*type director struct {
	builder iBuilder
}

func newDirector(b iBuilder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b iBuilder) {
	d.builder = b
}

func (d *director) buildHouse() House {
	d.builder.setDoorType()
	d.builder.setWindowType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}*/
