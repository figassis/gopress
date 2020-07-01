package appointment

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/appointment"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new appointment logging service
func New(svc appointment.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents appointment logging service
type LogService struct {
	appointment.Service
	logger model.Logger
}

const name = "appointment"

// Create logging
func (ls *LogService) Create(c echo.Context, req *appointment.Create) (resp *model.Appointment, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create appointment request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Appointment, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List appointment request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Appointment, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View appointment request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete appointment request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *appointment.Update) (resp *model.Appointment, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update appointment request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
