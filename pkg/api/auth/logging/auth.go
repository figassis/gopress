package auth

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/auth"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new auth logging service
func New(svc auth.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents auth logging service
type LogService struct {
	auth.Service
	logger model.Logger
}

const name = "auth"

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, user, password string) (resp *model.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Authenticate request", err,
			map[string]interface{}{
				"req":  user,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Authenticate(c, user, password)
}

// Refresh logging
func (ls *LogService) Refresh(c echo.Context) (resp *model.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Refresh request", err,
			map[string]interface{}{
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Refresh(c)
}

func (ls *LogService) Resend(c echo.Context) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Resend confirm email", err,
			map[string]interface{}{
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Resend(c)
}

// Me logging
func (ls *LogService) Me(c echo.Context) (resp *model.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Me request", err,
			map[string]interface{}{
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Me(c)
}
