package transport

import (
	"net/http"
	"strconv"

	"github.com/figassis/goinagbe/pkg/api/course"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

// HTTP represents user http service
type HTTP struct {
	svc course.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc course.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/courses")
	ur.POST("", h.create)
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.PATCH("/:id", h.update)
	ur.DELETE("/:id", h.delete)
}

// User create request
func (h *HTTP) create(c echo.Context) error {
	req := course.Create{}

	if err := c.Bind(&req); err != nil {

		return err
	}

	usr, err := h.svc.Create(c, &req)

	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) list(c echo.Context) error {
	page, _ := strconv.Atoi(c.Request().Header.Get("Page"))
	limit, _ := strconv.Atoi(c.Request().Header.Get("Limit"))
	if limit < 0 {
		limit = 20
	}
	p := model.Pagination{Page: page, Limit: limit, Cursor: c.Request().Header.Get("Cursor"), CacheKey: c.QueryString()}
	result, next, prev, total, pages, err := h.svc.List(c, &p)

	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	c.Response().Header().Set("Page", strconv.Itoa(p.Page))
	c.Response().Header().Set("NextCursor", next)
	c.Response().Header().Set("PreviousCursor", prev)
	c.Response().Header().Set("TotalResults", strconv.Itoa(int(total)))
	c.Response().Header().Set("TotalPages", strconv.Itoa(int(pages)))
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) view(c echo.Context) error {
	result, err := h.svc.View(c, c.Param("id"))
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) update(c echo.Context) error {
	req := course.Update{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	req.ID = c.Param("id")

	usr, err := h.svc.Update(c, &req)

	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	if err := h.svc.Delete(c, c.Param("id")); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
