package post

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/post"
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new post logging service
func New(svc post.Service, logger model.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents post logging service
type LogService struct {
	post.Service
	logger model.Logger
}

const name = "post"

// Create logging
func (ls *LogService) Create(c echo.Context, req *post.Create) (resp *model.Post, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Create post request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *model.Pagination) (resp []model.Post, next, prev string, total, pages int64, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "List post request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req string) (resp *model.Post, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "View post request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Delete post request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *post.Update) (resp *model.Post, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(c, name, "Update post request", err, map[string]interface{}{"req": req, "took": time.Since(begin)})
	}(time.Now())
	return ls.Service.Update(c, req)
}
