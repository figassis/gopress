package auth

import (
	"errors"
	"net/http"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Incorrect username or password")
	ErrInsecurePassword   = echo.NewHTTPError(http.StatusBadRequest, "Senha insegura")
)

// Authenticate tries to authenticate the user provided by username and password
func (a *Auth) Authenticate(c echo.Context, user, pass string) (*model.AuthToken, error) {
	if user == model.AdminEmail {
		//System user cannot login
		return nil, model.ErrUnauthorized
	}

	u, err := a.udb.FindByUsername(a.db, user)
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	if !a.sec.HashMatchesPassword(u.Password, pass) {
		return nil, ErrInvalidCredentials
	}

	if util.CheckWordpressPassword(pass, u.Password) {
		u.ChangePassword(a.sec.Hash(pass))
		zaplog.ZLog(a.udb.Update(a.db, &model.User{Base: model.Base{ID: u.ID}, Password: u.Password}))
	}

	if u.Status != model.StatusActive && u.Status != model.StatusPending {
		return nil, model.ErrUnauthorized
	}

	token, expire, refreshToken, refreshExpire, err := a.tg.GenerateToken(u)
	if err = zaplog.ZLog(err); err != nil {
		return nil, model.ErrUnauthorized
	}

	u.UpdateLastLogin(a.sec.Token(token))

	if err := a.udb.Update(a.db, u); err != nil {
		return nil, zaplog.ZLog(err)
	}

	return &model.AuthToken{Token: token, Expires: expire, RefreshToken: refreshToken, RefreshTokenExpires: refreshExpire}, nil
}

func (a *Auth) Contact(c echo.Context, req model.Contact) error {
	return zaplog.ZLog(model.SendContactEmail(req))
}

func (a *Auth) Signup(c echo.Context, req Signup) error {
	email, err := util.ValidateEmail(req.Email)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if _, err := a.udb.FindByUsername(a.db, email); err == nil {
		return errors.New("Este utilizador já existe")
	}

	ids, err := util.GenerateUUIDS(2)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	var org model.Organization
	if err = zaplog.ZLog(a.db.Model(&model.Organization{}).Where("type = ?", model.OrgMain).First(&org).Error); err != nil {
		return err
	}

	req.Password = a.sec.Hash(req.Password)
	user := model.User{
		Base:             model.Base{ID: ids[0]},
		Name:             req.Name,
		Username:         email,
		Password:         req.Password,
		Email:            email,
		Phone:            req.Phone,
		Status:           model.StatusPending,
		Role:             model.CandidateRole,
		Organization:     org.ID,
		OrganizationName: org.Name,
		UnsubscribeID:    ids[1],
	}

	return zaplog.ZLog(a.udb.Signup(a.db, &user))
}

// Refresh refreshes jwt token and puts new claims inside
func (a *Auth) Refresh(c echo.Context) (*model.AuthToken, error) {
	au := a.rbac.User(c)
	user, err := a.udb.View(a.db, au.ID)
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	if user.Status != model.StatusActive {
		return nil, model.ErrUnauthorized
	}

	token, expire, refreshToken, refreshExpire, err := a.tg.GenerateToken(user)
	if err = zaplog.ZLog(err); err != nil {
		return nil, model.ErrUnauthorized
	}

	user.UpdateLastLogin(a.sec.Token(token))
	if err := a.udb.Update(a.db, user); err != nil {
		return nil, zaplog.ZLog(err)
	}

	return &model.AuthToken{Token: token, Expires: expire, RefreshToken: refreshToken, RefreshTokenExpires: refreshExpire}, nil
}

// Me returns info about currently logged user
func (a *Auth) Me(c echo.Context) (*model.User, error) {
	au := a.rbac.User(c)
	return a.udb.View(a.db, au.ID)
}

func (a *Auth) GetPublicData(c echo.Context) (*model.Public, error) {
	return a.udb.GetPublicData(a.db)
}

func (a *Auth) ConfirmEmail(c echo.Context, token string) error {
	return a.udb.ConfirmEmail(a.db, token)
}

// Me returns info about currently logged user
func (a *Auth) Resend(c echo.Context) error {
	au := a.rbac.User(c)
	user, err := a.udb.View(a.db, au.ID)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return user.SendConfirmationEmail()
}

// Me returns info about currently logged user
func (a *Auth) Unsubscribe(c echo.Context, token string) error {
	return a.udb.Unsubscribe(a.db, token)
}

func (a *Auth) Bounce(c echo.Context, n model.BounceNotification) error {
	return a.udb.Bounce(a.db, n)
}

func (p *Auth) Reset(c echo.Context, email string) error {
	u, err := p.udb.FindByUsername(p.db, email)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	return zaplog.ZLog(u.SendResetEmail())
}

func (p *Auth) CheckResetToken(c echo.Context, token string) error {
	var userID string
	if err := util.GetCache(token, &userID); err != nil {
		return err
	}

	_, err := p.udb.View(p.db, userID)
	return err
}

func (p *Auth) CompleteReset(c echo.Context, token, newPass string) error {
	var userID string
	if err := util.GetCache(token, &userID); err != nil {
		return err
	}

	u, err := p.udb.View(p.db, userID)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if !p.sec.Password(newPass, u.Name, u.Username, u.Email) {
		return ErrInsecurePassword
	}

	u.ChangePassword(p.sec.Hash(newPass))

	return p.udb.Update(p.db, u)
}
