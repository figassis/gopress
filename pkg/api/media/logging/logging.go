package media

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/media"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new media logging service
func New(svc media.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents media logging service
type LogService struct {
	media.Service
	logger model.Logger
}

const name = "media"

// Create logging
func (ls *LogService) Create(c echo.Context, req *media.Create) (resp *model.File, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create media request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.File, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List media request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.File, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View media request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete media request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *media.Update) (resp *model.File, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update media request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
