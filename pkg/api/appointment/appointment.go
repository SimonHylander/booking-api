package appointment

import (
	"github.com/labstack/echo"
	"github.com/simonhylander/booker"
)

// Service represents appointment application interface
type Service interface {
	Create(echo.Context, booker.User) (booker.User, error)
	List(echo.Context, booker.Pagination) ([]booker.User, error)
	View(echo.Context, int) (booker.User, error)
	Delete(echo.Context, int) error
	//Update(echo.Context, int) (booker.User, error)
}

// New creates new user application service
func New(rbac RBAC, sec Securer) *User {
	return &User{rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(rbac RBAC, sec Securer) *User {
	return New(rbac, sec)
}

// User represents user application service
type User struct {
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(booker.User) (booker.User, error)
	View(int) (booker.User, error)
	List(*booker.ListQuery, booker.Pagination) ([]booker.User, error)
	//Update(booker.User) error
	Delete(booker.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) booker.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, booker.AccessRole, int, int) error
	IsLowerRole(echo.Context, booker.AccessRole) error
}


func (u User) Create(echo.Context, booker.User) (booker.User, error) {
	panic("implement me")
}

func (u User) List(echo.Context, booker.Pagination) ([]booker.User, error) {
	panic("implement me")
}

func (u User) View(echo.Context, int) (booker.User, error) {
	panic("implement me")
}

func (u User) Delete(echo.Context, int) error {
	panic("implement me")
}
