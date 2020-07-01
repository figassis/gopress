package transport

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/figassis/goinagbe/pkg/api/media"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/gabriel-vasile/mimetype"

	echo "github.com/labstack/echo/v4"
)

// HTTP represents user http service
type HTTP struct {
	svc media.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc media.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/media")
	ur.POST("", h.create)
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.PATCH("/:id", h.update)
	ur.DELETE("/:id", h.delete)
}

// User create request
func (h *HTTP) create(c echo.Context) error {
	file, err := c.FormFile("File")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fileID, err := util.GenerateUUID()
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", os.Getenv("UPLOADS"), fileID)

	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	mime, err := mimetype.DetectFile(filePath)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	req := media.Create{
		Resource:   c.FormValue("Resource"),
		ResourceID: c.FormValue("ResourceID"),
		Type:       c.FormValue("Type"),
		Extension:  mime.Extension(),
		Name:       file.Filename,
		Path:       fmt.Sprintf("%s/%s", os.Getenv("UPLOADS"), fileID),
		Public:     strings.EqualFold(c.FormValue("Access"), model.AccessPublic),
	}

	if req.Type == "" {
		switch req.Extension {
		case "jpg", "png", "bmp", "tiff", "webp":
			req.Type = model.ImageFile
		default:
			req.Type = model.GeneralFile
		}
	}

	usr, err := h.svc.Create(c, &req)

	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) list(c echo.Context) error {
	// p := new(model.Pagination)
	// if err := c.Bind(p); err != nil {
	// 	return err
	// }

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
	req := media.Update{}
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
