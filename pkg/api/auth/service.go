package auth

import (
	"github.com/labstack/echo"
	"github.com/simonhylander/booker"
	"github.com/simonhylander/booker/pkg/api/auth/platform/neo4j"
)

// New creates new iam service
func New(udb UserDB, tokenGenerator TokenGenerator, securer Securer, rbac RBAC) Auth {
	return Auth{
		udb:            udb,
		tokenGenerator: tokenGenerator,
		securer:        securer,
		rbac:           rbac,
	}
}

// Initialize initializes auth application service
func Initialize(tokenGenerator TokenGenerator, sec Securer, rbac RBAC) Auth {
	return New(neo4j.User{}, tokenGenerator, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (booker.AuthToken, error)
	//Refresh(echo.Context, string) (string, error)
	//Me(echo.Context) (booker.User, error)
}

// Auth represents auth application service
type Auth struct {
	//db   *pg.DB
	udb            UserDB
	tokenGenerator TokenGenerator
	securer        Securer
	rbac           RBAC
}

// UserDB represents user repository interface
type UserDB interface {
	/*View(orm.DB, int) (booker.User, error)
	FindByToken(orm.DB, string) (booker.User, error)*/
	FindByUsername(string) (booker.User, error)
	Update(booker.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(booker.User) (string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) booker.AuthUser
}
