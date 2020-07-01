package application

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/application"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new application logging service
func New(svc application.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents application logging service
type LogService struct {
	application.Service
	logger model.Logger
}

const name = "application"

// Create logging
func (ls *LogService) Create(c echo.Context, req *application.Create) (resp *model.Application, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create application request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Application, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List application request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Application, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View application request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete application request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *application.Update) (resp *model.Application, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update application request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
