package schedule

import (
	"github.com/labstack/echo"
	"net/http"
)

type Schedule struct {
	Id string `json:"id"`
	NameOfDay string `json:"nameOfDay"`
	Day int `json:"day"`
	IsOpen bool `json:"isOpen"`
	Start string `json:"start"`
	End string `json:"end"`
}

type ScheduleHandler struct {
	ScheduleRepository ScheduleRepository
}

func (e *ScheduleHandler) Create(c echo.Context) error {
	schedule := new(Schedule)

	if err := c.Bind(schedule); err != nil {
		return nil
	}

	createdEmployee, err := e.ScheduleRepository.Create(schedule)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusCreated, createdEmployee)
}

func (s *ScheduleHandler) Find(c echo.Context) error {
	schedule, _ := s.ScheduleRepository.Find()
	return c.JSON(http.StatusOK, schedule)
}