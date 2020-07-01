package transport

import (
	"net/http"

	"github.com/figassis/goinagbe/pkg/api/password"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

// HTTP represents password http transport service
type HTTP struct {
	svc password.Service
}

// NewHTTP creates new password http service
func NewHTTP(svc password.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/password")
	pr.PATCH("/:id", h.change)
}

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
)

func (h *HTTP) change(c echo.Context) error {
	p := password.Update{}
	if err := c.Bind(&p); err != nil {
		return zaplog.ZLog(err)
	}

	if p.NewPassword != p.NewPasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if err := h.svc.Update(c, c.Param("id"), p.OldPassword, p.NewPassword); err != nil {
		return zaplog.ZLog(err)
	}

	return c.NoContent(http.StatusOK)
}
