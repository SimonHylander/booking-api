package employee

import (
	"github.com/labstack/echo"
	"net/http"
)

type Employee struct {
	Id string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Fullname string `json:"fullname"`
	Titles []string `json:"titles"`
	CompanyId string `json:"companyId,omitempty"`
	Company string `json:"company,omitempty"`
	Image string `json:"img,omitempty"`
}

type Title struct {
	Name string `json:"name"`
}

type EmployeeHandler struct {
	EmployeeRepository EmployeeRepository
}

// TODO swagger

func (e *EmployeeHandler) CreateEmployee(c echo.Context) error {
	employee := new(Employee)

	if err := c.Bind(employee); err != nil {
		return nil
	}

	createdEmployee, err := e.EmployeeRepository.Create(employee)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusCreated, createdEmployee)
}

func (e *EmployeeHandler) Find(c echo.Context) error {
	employees, _ := e.EmployeeRepository.Find()
	return c.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) Get(c echo.Context) error {
	employee, _ := e.EmployeeRepository.Get(c.Param("id"))
	return c.JSON(http.StatusOK, employee)
}