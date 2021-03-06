package user

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/user"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new user logging service
func New(svc user.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents user logging service
type LogService struct {
	user.Service
	logger model.Logger
}

const name = "user"

// Create logging
func (ls *LogService) Create(c echo.Context, req user.Create) (resp *model.User, err error) {
	defer func(begin time.Time) {
		req.Password = "xxx-redacted-xxx"
		ls.logger.Log(
			c,
			name, "Create user request", err,
			map[string]interface{}{
				"req": req,
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.User, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List user request", err,
			map[string]interface{}{
				"req": req,
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "View user request", err,
			map[string]interface{}{
				"req": req,
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Delete user request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *user.Update) (resp *model.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Update user request", err,
			map[string]interface{}{
				"req": req,
				// "resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
