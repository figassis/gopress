package transport

import (
	"errors"
	"net/http"

	"github.com/figassis/goinagbe/pkg/api/auth"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

// NewHTTP creates new auth http service
func NewHTTP(svc auth.Service, e *echo.Echo, mw echo.MiddlewareFunc) {
	h := HTTP{svc}

	e.POST("/login", h.login)
	e.POST("/signup", h.signup)
	e.POST("/refresh", h.refresh, mw)
	e.GET("/me", h.me, mw)
	e.GET("/public", h.public)

	e.POST("/reset", h.reset)                  //frontend will submit the first reset request to this endpoint, providing the email address
	e.GET("/reset/:token", h.checkResetToken)  //after the user clicks on the email link, the frontend will check the token at this endpoint
	e.POST("/reset/complete", h.completeReset) //after the token is validated, the frontend will submit the token and new password

	e.GET("/confirm/:token", h.confirm)
	e.POST("/resend", h.resend, mw)
	e.GET("/unsubscribe/:token", h.unsubscribe)
	e.POST("/bounce", h.bounce)
	e.POST("/contact", h.contact)
}

func (h *HTTP) login(c echo.Context) error {
	cred := auth.Credentials{}
	if err := c.Bind(&cred); err != nil {
		return err
	}
	r, err := h.svc.Authenticate(c, cred.Username, cred.Password)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (h *HTTP) signup(c echo.Context) error {
	cred := auth.Signup{}
	if err := c.Bind(&cred); err != nil {
		return err
	}

	//Honeypot must be md5 hash of email
	if util.MD5(cred.Email) != cred.Honeypot {
		return errors.New("Caught a bee!")
	}

	if err := zaplog.ZLog(h.svc.Signup(c, cred)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) contact(c echo.Context) error {
	cred := model.Contact{}
	if err := c.Bind(&cred); err != nil {
		return err
	}

	//Honeypot must be md5 hash of email
	if util.MD5(cred.Email) != cred.Honeypot {
		return errors.New("Caught a bee!")
	}

	if err := zaplog.ZLog(h.svc.Contact(c, cred)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) refresh(c echo.Context) error {
	r, err := h.svc.Refresh(c)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (h *HTTP) resend(c echo.Context) error {
	if err := zaplog.ZLog(h.svc.Resend(c)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) me(c echo.Context) error {
	user, err := h.svc.Me(c)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HTTP) public(c echo.Context) error {
	data, err := h.svc.GetPublicData(c)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (h *HTTP) confirm(c echo.Context) error {
	if err := zaplog.ZLog(h.svc.ConfirmEmail(c, c.Param("token"))); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) unsubscribe(c echo.Context) error {
	if err := zaplog.ZLog(h.svc.Unsubscribe(c, c.Param("token"))); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) bounce(c echo.Context) error {
	n := model.BounceNotification{}
	if err := c.Bind(&n); err != nil {
		return err
	}
	if err := zaplog.ZLog(h.svc.Bounce(c, n)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) reset(c echo.Context) error {
	p := auth.Reset{}
	if err := c.Bind(&p); err != nil {
		return zaplog.ZLog(err)
	}

	if err := h.svc.Reset(c, p.Email); err != nil {
		return zaplog.ZLog(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) checkResetToken(c echo.Context) error {
	if err := h.svc.CheckResetToken(c, c.Param("token")); err != nil {
		return zaplog.ZLog(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) completeReset(c echo.Context) error {
	p := auth.CompleteReset{}
	if err := c.Bind(&p); err != nil {
		return zaplog.ZLog(err)
	}

	if p.NewPassword != p.NewPasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if err := h.svc.CompleteReset(c, p.Token, p.NewPassword); err != nil {
		return zaplog.ZLog(err)
	}

	return c.NoContent(http.StatusOK)
}
