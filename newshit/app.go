package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"github.com/simonhylander/booking-api/company"
	"github.com/simonhylander/booking-api/db"
	"github.com/simonhylander/booking-api/employee"
	"github.com/simonhylander/booking-api/schedule"
	"net/http"
)

func main() {
	neo4jUri := "bolt://localhost:7687" //os.LookupEnv("NEO4J_URI")
	neo4jUsername := "neo4j"            //os.LookupEnv("NEO4J_USERNAME")
	neo4jPassword := "booking"          //os.LookupEnv("NEO4J_PASSWORD")

	neo4jDriver, err := neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))

	if err != nil {
		panic(err)
	}

	defer neo4jDriver.Close()

	companyHandler := &company.CompanyHandler{
		CompanyRepository: &company.CompanyNeo4jRepository{
			Driver: neo4jDriver,
		},
	}

	employeeHandler := &employee.EmployeeHandler{
		EmployeeRepository: &employee.EmployeeNeo4jRepository{
			Driver: neo4jDriver,
		},
	}

	scheduleHandler := &schedule.ScheduleHandler{
		/*ScheduleRepository: &schedule.ScheduleNeo4jRepository{
			Driver: neo4jDriver,
		},*/
	}

	e := echo.New()

	signingKey := []byte("secret")

	config := middleware.JWTConfig{
		TokenLookup: "query:token",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingKey, nil
			}

			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
	}

	e.Use(middleware.JWTWithConfig(config))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*","http://localhost", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	err = db.AutoMigrate(neo4jDriver)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrated neo4j")

	v1 := e.Group("/api")
	v1.POST("/companies", companyHandler.CreateCompany)
	v1.GET("/companies", companyHandler.FindCompanies)
	v1.POST("/employees", employeeHandler.CreateEmployee)
	v1.GET("/employees", employeeHandler.Find)
	v1.GET("/employees/:id", employeeHandler.Get)

	
	/*v1.POST("/schedule", scheduleHandler.CreateEmployee)*/
	v1.GET("/schedule", scheduleHandler.Find)

	//e.GET("/book/available", GetAvailableBookings)
	e.Logger.Fatal(e.Start(":8080"))
}

type AvailableBooking struct {
	EmployeeId   int    `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	Start        string `json:"start"`
}

func GetAvailableBookings(c echo.Context) error {
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	fmt.Println("GetAvailableBookings")
	fmt.Println(start)
	fmt.Println(end)

	availableBookings := []AvailableBooking{
		{
			Start:        "2021-07-12 09:00:00",
			EmployeeId:   1,
			EmployeeName: "Sabina Hylander",
		},
		{
			Start:        "2021-07-12 10:00:00",
			EmployeeId:   1,
			EmployeeName: "Sabina Hylander",
		},
		{
			Start:        "2021-07-12 12:00:00",
			EmployeeId:   1,
			EmployeeName: "Sabina Hylander",
		},
		{
			Start:        "2021-07-12 14:00:00",
			EmployeeId:   1,
			EmployeeName: "Sabina Hylander",
		},
		{
			Start:        "2021-07-12 15:00:00",
			EmployeeId:   1,
			EmployeeName: "Sabina Hylander",
		},
	}

	return c.JSON(http.StatusOK, availableBookings)
}
