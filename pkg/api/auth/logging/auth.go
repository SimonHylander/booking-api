package auth

import (
	"github.com/labstack/echo"
	"github.com/simonhylander/booker"
	"github.com/simonhylander/booker/pkg/api/auth"
)

// New creates new auth logging service
func New(svc auth.Service, logger booker.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents auth logging service
type LogService struct {
	auth.Service
	logger booker.Logger
}

const name = "auth"

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, username, password string) (resp booker.AuthToken, err error) {
	/*defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Authenticate request", err,
			map[string]interface{}{
				"req":  user,
				"took": time.Since(begin),
			},
		)
	}(time.Now())*/

	return ls.Service.Authenticate(c, username, password)
}

// Refresh logging
/*func (ls *LogService) Refresh(c echo.Context, req string) (token string, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Refresh request", err,
			map[string]interface{}{
				"req":  req,
				"resp": token,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Refresh(c, req)
}

// Me logging
func (ls *LogService) Me(c echo.Context) (resp booker.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Me request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Me(c)
}*/