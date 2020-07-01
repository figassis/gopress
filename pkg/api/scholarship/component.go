// Package user contains user application services
package scholarship

import (
	"errors"
	"fmt"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Scholarship, error) {
	var org, mainOrg model.Organization

	if err := u.db.Model(&model.Organization{}).Where("type = ?", model.OrgMain).First(&mainOrg).Error; err != nil {
		return nil, err
	}

	if err := u.db.Model(&model.Organization{}).Where("uuid = ?", req.Sponsor).First(&org).Error; err == nil {
		req.Sponsor = mainOrg.ID
	}

	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", req.Documents).Update(model.File{
		Resource:   model.ResourceScholarship,
		ResourceID: id,
	}).Error)

	if !req.ProvinceQuota.Valid() {
		req.ProvinceQuota = model.ProvinceQuota{}
	}

	req.ProvinceQuota.Scholarship = id

	now := time.Now()
	if req.Status == model.StatusScholarshipPublished && (req.Start.Before(now) || req.Start.Equal(now)) && req.End.After(now) {
		req.Status = model.StatusOpen
	}

	scholarship := model.Scholarship{
		Base:              model.Base{ID: id},
		Name:              req.Name,
		Sponsor:           req.Sponsor,
		SponsorName:       org.Name,
		Start:             req.Start,
		End:               req.End,
		Available:         req.Available,
		MaxAge:            req.MaxAge,
		MinGrade:          req.MinGrade,
		RPEQuota:          req.RPEQuota,
		PriorityQuota:     req.PriorityQuota,
		Content:           req.Content,
		ProvinceQuota:     req.ProvinceQuota,
		Type:              req.Type,
		Level:             req.Level,
		Status:            req.Status,
		EnableQuotas:      req.EnableQuotas,
		RequiredDocuments: req.RequiredDocuments,
	}

	return u.udb.Create(u.db, scholarship)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Scholarship, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourceScholarship)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.ScholarshipQuery = &model.Scholarship{}
		params := c.QueryParams()

		p.ScholarshipQuery.Level = params.Get("level")
		p.ScholarshipQuery.Type = params.Get("type")
		p.ScholarshipQuery.Status = params.Get("status")
		p.ScholarshipQuery.Sponsor = params.Get("sponsor")
		p.ScholarshipQuery.SponsorName = params.Get("sponsor_name")
		p.SearchQuery = params.Get("s")
	}
	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Scholarship, error) {
	if err := u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
		return nil, err
	}
	return u.udb.View(u.db, id)
}

func (u *App) Export(c echo.Context, id string) (string, error) {
	if err := u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
		return "", err
	}
	scholarship, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return "", err
	}

	return scholarship.Export(u.db, false)
}

// Delete deletes a user
func (u *App) Delete(c echo.Context, id string) error {
	if err := u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return err
	}

	item, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if item.Status != model.StatusDraft {
		return errors.New("Apenas é possível eliminar bolsas em rascunho")
	}

	var count int
	if err := u.db.Model(&model.Application{}).Where("scholarship = ?", id).Count(&count).Error; err == nil && count > 0 {
		return errors.New("Elimine primeiro as candidaturas associadas a esta bolsa")
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Scholarship, err error) {
	if err = u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return
	}

	update, err := u.udb.View(u.db, c.Param("id"))
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if r.Status != "" && update.Status != r.Status && !update.AllowedStatuses(r.Status) {
		err = fmt.Errorf("Não é possível passar de %s para %s", update.Status, r.Status)
		return
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", r.Documents).Update(model.File{
		Resource:   model.ResourceScholarship,
		ResourceID: update.ID,
	}).Error)

	if !r.ProvinceQuota.Valid() && r.EnableQuotas {
		return nil, errors.New("As quotas são inválidas")
	}

	var quota model.ProvinceQuota
	if err = zaplog.ZLog(u.db.Model(&model.ProvinceQuota{}).Where("scholarship = ?", r.ID).First(&quota).Error); err != nil {
		quota.ID, err = util.GenerateUUID()
		if err = zaplog.ZLog(err); err != nil {
			return
		}

		quota.Scholarship = update.ID

		if err = zaplog.ZLog(u.db.Debug().Create(&quota).Error); err != nil {
			return
		}

	}

	if r.EnableQuotas {
		r.ProvinceQuota.IntID, r.ProvinceQuota.ID, r.ProvinceQuota.Scholarship = quota.IntID, quota.ID, update.ID
	}

	update.Name = r.Name
	update.Sponsor = r.Sponsor
	update.Available = r.Available
	update.MaxAge = r.MaxAge
	update.MinGrade = r.MinGrade
	update.Content = r.Content
	update.Status = r.Status
	update.RequiredDocuments = r.RequiredDocuments
	update.RPEQuota = r.RPEQuota
	update.PriorityQuota = r.PriorityQuota
	update.EnableQuotas = r.EnableQuotas
	update.ProvinceQuota = r.ProvinceQuota

	if r.Start.After(time.Now()) {
		update.Start = r.Start

	}
	if r.End.After(update.Start) {
		update.End = r.End
	}

	if update.Status == model.StatusDraft {
		var count int
		if err = u.db.Model(&model.Application{}).Where("scholarship = ?", update.ID).Count(&count).Error; err != nil {
			return
		}

		if (r.Level != "" && update.Level != r.Level) || (r.Type != "" && update.Type != r.Type) {
			if count > 0 {
				err = zaplog.ZLog(errors.New("Não pode alterar o tipo ou nível de uma bolsa com candidaturas"))
				return
			}
			update.Level = r.Level
			update.Type = r.Type
		}
	}

	now := time.Now()
	if update.Status == model.StatusScholarshipPublished && (update.Start.Before(now) || update.Start.Equal(now)) && update.End.After(now) {
		update.Status = model.StatusOpen
	}

	var org model.Organization
	if err = u.db.Model(&model.Organization{}).Where("uuid = ?", r.Sponsor).First(&org).Error; err == nil {
		update.SponsorName = org.Name
	}

	if err = zaplog.ZLog(u.udb.Update(u.db, update)); err != nil {
		return
	}

	result, err = u.udb.View(u.db, update.ID)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	return u.udb.View(u.db, update.ID)
}
