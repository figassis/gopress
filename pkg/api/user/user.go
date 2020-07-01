// Package user contains user application services
package user

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *User) Create(c echo.Context, req Create) (*model.User, error) {

	if !req.Role.Exists() {
		return nil, errors.New("Invalid role")
	}

	if err := u.rbac.AccountCreate(c, req.Role, req.Organization); err != nil {
		return nil, err
	}

	email, err := util.ValidateEmail(req.Email)
	if err := zaplog.ZLog(err); err != nil {
		return nil, err
	}

	if _, err := u.udb.FindByUsername(u.db, email); err == nil {
		return nil, errors.New("Este utilizador já existe")
	}

	req.Password = u.sec.Hash(req.Password)

	uuids, err := util.GenerateUUIDS(2)
	if err != nil {
		return nil, err
	}

	var org model.Organization
	if err := u.db.Model(&model.Organization{}).Where("uuid = ?", req.Organization).First(&org).Error; err != nil {
		req.Organization = ""
	}

	user := model.User{
		Base:             model.Base{ID: uuids[0]},
		Name:             req.Name,
		Username:         email,
		Password:         req.Password,
		Email:            email,
		Phone:            req.Phone,
		Status:           model.StatusActive,
		Role:             req.Role,
		Organization:     req.Organization,
		OrganizationName: org.Name,
		UnsubscribeID:    uuids[1],
	}
	return u.udb.Create(u.db, user)
}

// List returns list of users
func (u *User) List(c echo.Context, p *model.Pagination) ([]model.User, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourceUser)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.UserQuery = &model.User{}
		params := c.QueryParams()
		if role, err := strconv.Atoi(params.Get("role")); err == nil {
			p.UserQuery.Role = model.AccessRole(role)
		}
		p.UserQuery.Organization = params.Get("organization")
		p.UserQuery.Status = params.Get("status")
		p.SearchQuery = params.Get("s")
	}

	// var users []model.User
	// if err = util.GetCache(fmt.Sprintf("/users_%s_%s", p.CacheKey,au.ID), &users); err == nil {

	// }

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *User) View(c echo.Context, id string) (*model.User, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}

	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u *User) Delete(c echo.Context, id string) error {
	user, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}

	if err := u.rbac.EnforceRole(c, model.CompanyAdminRole); err != nil {
		return err
	}

	if err := u.rbac.EnforceCompany(c, user.Organization); err != nil {
		return err
	}

	if err := u.rbac.IsLowerRole(c, user.Role); err != nil {
		return err
	}

	var count int
	if err = zaplog.ZLog(u.db.Model(&model.Application{}).Where("user_id = ?", id).Count(&count).Error); err == nil && count > 0 {
		return zaplog.ZLog(fmt.Errorf("Não é possível eliminar utilizadores com candidaturas associadas."))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *User) Update(c echo.Context, r *Update) (user *model.User, err error) {
	if err := u.rbac.EnforceUser(c, r.ID); err != nil {
		return nil, err
	}

	user, err = u.udb.View(u.db, r.ID)
	if err != nil {
		return
	}

	if r.Status != "" && user.Status != r.Status && !user.AllowedStatuses(r.Status) {
		err = fmt.Errorf("Não é possível passar de %s para %s", user.Status, r.Status)
		return
	}

	if err = u.rbac.IsLowerRole(c, user.Role); err != nil {
		return
	}

	//Never elevate a user to a higher role than current user
	if err = u.rbac.IsLowerRole(c, r.Role); err != nil {
		return
	}

	if _, err = u.udb.Update(u.db, &model.User{
		Base:         model.Base{ID: r.ID},
		Name:         r.Name,
		Phone:        r.Phone,
		Status:       r.Status,
		Role:         r.Role,
		Unsubscribed: r.Unsubscribed,
	}); err != nil {
		return nil, err
	}

	return u.udb.View(u.db, r.ID)
}
