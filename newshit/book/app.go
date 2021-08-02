package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	//lambda.Start(HandleRequest)
	e := echo.New()
	e.GET("/book/available", GetAvailableBookings)
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
	//return c.JSON(http.StatusBadRequest, map[string]string{"error": "Please specify the data type as Sting or JSON",})
}
