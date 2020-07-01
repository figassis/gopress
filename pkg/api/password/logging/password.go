package password

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/password"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new password logging service
func New(svc password.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	password.Service
	logger model.Logger
}

const name = "password"

// Change logging
func (ls *LogService) Change(c echo.Context, id string, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, id, oldPass, newPass)
}
