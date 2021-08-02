package company

import (
	"github.com/labstack/echo"
	"net/http"
)

type Company struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type CompanyHandler struct {
	CompanyRepository CompanyRepository
}

func (c *CompanyHandler) CreateCompany(ctx echo.Context) error {
	company := new(Company)

	if err := ctx.Bind(company); err != nil {
		return nil
	}

	err := c.CompanyRepository.Create(company)

	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (e *CompanyHandler) FindCompanies(c echo.Context) error {
	companies, _ := e.CompanyRepository.FindAll()
	return c.JSON(http.StatusOK, companies)
}