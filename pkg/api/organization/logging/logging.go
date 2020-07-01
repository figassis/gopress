package organization

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/organization"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new organization logging service
func New(svc organization.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents organization logging service
type LogService struct {
	organization.Service
	logger model.Logger
}

const name = "organization"

// Create logging
func (ls *LogService) Create(c echo.Context, req *organization.Create) (resp *model.Organization, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create organization request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Organization, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List organization request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Organization, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View organization request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete organization request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *organization.Update) (resp *model.Organization, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update organization request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
