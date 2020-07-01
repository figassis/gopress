package scholarship

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/scholarship"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new scholarship logging service
func New(svc scholarship.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents scholarship logging service
type LogService struct {
	scholarship.Service
	logger model.Logger
}

const name = "scholarship"

// Create logging
func (ls *LogService) Create(c echo.Context, req *scholarship.Create) (resp *model.Scholarship, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Scholarship, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Scholarship, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Export logging
func (ls *LogService) Export(c echo.Context, req string) (resp string, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Export scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Export(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *scholarship.Update) (resp *model.Scholarship, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update scholarship request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
