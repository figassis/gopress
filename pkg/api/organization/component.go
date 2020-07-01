// Package user contains user application services
package organization

import (
	"errors"
	"fmt"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Organization, error) {
	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", req.Documents).Update(model.File{
		Resource:   model.ResourceOrganization,
		ResourceID: id,
	}).Error)

	if req.Logo == "" {
		req.Logo = "assets/images/thumb-1.jpg"
	}

	var logo model.File
	if err = zaplog.ZLog(u.db.Model(&model.File{}).Where("resource = ? AND type = ?", model.ResourceOrganization, model.ProfileImage).Last(&logo).Error); err == nil {
		url, _ := logo.GetURL()
		if url != "" {
			req.Logo = url
		}
	}

	organization := model.Organization{
		Base:     model.Base{ID: id},
		Name:     req.Name,
		Status:   req.Status,
		Type:     req.Type,
		Country:  req.Country,
		Province: req.Province,
		City:     req.City,
		Logo:     req.Logo,
		Phone:    req.Phone,
		Email:    req.Email,
		Website:  req.Website,
		Facebook: req.Facebook,
		Twitter:  req.Twitter,
		Linkedin: req.Linkedin,
		Address:  req.Address,
	}
	return u.udb.Create(u.db, organization)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Organization, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourceOrganization)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.OrganizationQuery = &model.Organization{}
		params := c.QueryParams()

		p.OrganizationQuery.Type = params.Get("type")
		p.OrganizationQuery.Status = params.Get("status")
		p.OrganizationQuery.Email = params.Get("email")
		p.SearchQuery = params.Get("s")
	}

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Organization, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u *App) Delete(c echo.Context, id string) error {
	if err := u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return err
	}

	org, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if org.Type == model.OrgMain {
		return errors.New("Esta organização não pode ser eliminada")
	}

	var user model.User
	if err := u.db.Model(&model.User{}).Where("organization = ?", id).First(&user).Error; err == nil && user.ID != "" {
		return errors.New("Elimine primeiro os utilizadores associados a esta organização")
	}

	var scholarship model.Scholarship
	if err := u.db.Model(&model.Scholarship{}).Where("sponsor = ?", id).First(&scholarship).Error; err == nil && scholarship.ID != "" {
		return errors.New("Elimine primeiro as bolsas associadas a esta organização")
	}

	var count int
	if err = zaplog.ZLog(u.db.Model(&model.Course{}).Where("school = ?", id).Count(&count).Error); err == nil && count > 0 {
		return zaplog.ZLog(fmt.Errorf("Não é possível eliminar organizações com cursos associados."))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Organization, err error) {
	if err = u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", r.Documents).Update(model.File{
		Resource:   model.ResourceOrganization,
		ResourceID: r.ID,
	}).Error)

	var logo model.File
	if err = zaplog.ZLog(u.db.Model(&model.File{}).Where("resource = ? AND resource_id = ? AND type = ?", model.ResourceOrganization, r.ID, model.ProfileImage).Last(&logo).Error); err == nil {
		url, _ := logo.GetURL()
		if url != "" {
			r.Logo = url
		}
	}

	old, err := u.udb.View(u.db, r.ID)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if r.Status != "" && old.Status != r.Status && !old.AllowedStatuses(r.Status) {
		err = fmt.Errorf("Não é possível passar de %s para %s", old.Status, r.Status)
		return
	}

	update := model.Organization{
		Base:     model.Base{ID: r.ID},
		Name:     r.Name,
		Status:   r.Status,
		Type:     r.Type,
		Country:  r.Country,
		Province: r.Province,
		City:     r.City,
		Logo:     r.Logo,
		Phone:    r.Phone,
		Email:    r.Email,
		Website:  r.Website,
		Facebook: r.Facebook,
		Twitter:  r.Twitter,
		Linkedin: r.Linkedin,
		Address:  r.Address,
	}
	if err = u.udb.Update(u.db, &update); err != nil {
		return
	}
	return u.udb.View(u.db, r.ID)
}
