package course

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/course"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new course logging service
func New(svc course.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents course logging service
type LogService struct {
	course.Service
	logger model.Logger
}

const name = "course"

// Create logging
func (ls *LogService) Create(c echo.Context, req *course.Create) (resp *model.Course, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create course request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Course, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List course request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Course, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View course request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete course request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *course.Update) (resp *model.Course, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update course request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
