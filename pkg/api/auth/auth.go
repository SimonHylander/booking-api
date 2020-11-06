package auth

import (
	"github.com/labstack/echo"
	"github.com/simonhylander/booker"
	"net/http"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
)

// Authenticate tries to authenticate the user provided by username and password
func (auth Auth) Authenticate(c echo.Context, username, pass string) (booker.AuthToken, error) {
	user, err := auth.udb.FindByUsername(username)

	if err != nil {
		return booker.AuthToken{}, err
	}

	if !auth.securer.HashMatchesPassword(user.Password, pass) {
		return booker.AuthToken{}, ErrInvalidCredentials
	}

	if !user.Active {
		return booker.AuthToken{}, booker.ErrUnauthorized
	}

	token, err := auth.tokenGenerator.GenerateToken(user)
	if err != nil {
		return booker.AuthToken{}, booker.ErrUnauthorized
	}

	user.UpdateLastLogin(auth.securer.Token(token))

	if err := auth.udb.Update(user); err != nil {
		return booker.AuthToken{}, err
	}

	return booker.AuthToken{Token: token, RefreshToken: user.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
/*func (a Auth) Refresh(c echo.Context, refreshToken string) (string, error) {
	user, err := a.udb.FindByToken(a.db, refreshToken)
	if err != nil {
		return "", err
	}
	return a.tg.GenerateToken(user)
}

// Me returns info about currently logged user
func (a Auth) Me(c echo.Context) (booker.User, error) {
	au := a.rbac.User(c)
	return a.udb.View(a.db, au.ID)
}*/