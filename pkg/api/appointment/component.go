// Package user contains user application services
package appointment

import (
	"fmt"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Appointment, error) {
	au := u.rbac.User(c)

	var operator model.User
	if err := u.db.Model(&model.User{}).Where("uuid = ?", req.Admin).First(&operator).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := u.db.Model(&model.User{}).Where("uuid = ?", au.ID).First(&user).Error; err != nil {
		return nil, err
	}

	appointment := model.Appointment{
		User:          au.ID,
		UserName:      user.Name,
		Resource:      req.Resource,
		ResourceID:    req.ResourceID,
		Date:          req.Date,
		Admin:         req.Admin,
		AdminName:     operator.Name,
		Status:        model.StatusPending,
		ContactName:   req.ContactName,
		ContactNumber: req.ContactNumber,
		ContactEmail:  req.ContactEmail,
		Message:       req.Message,
	}
	return u.udb.Create(u.db, appointment)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Appointment, string, string, int64, int64, error) {
	au := u.rbac.User(c)

	q, err := query.List(au, model.ResourceAppointment)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.AppointmentQuery = &model.Appointment{}
		params := c.QueryParams()

		p.AppointmentQuery.User = params.Get("user")
		p.AppointmentQuery.Admin = params.Get("admin")
		p.AppointmentQuery.Status = params.Get("status")
		p.SearchQuery = params.Get("s")
	}
	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Appointment, error) {
	appointment, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	if err := u.rbac.EnforceUser(c, appointment.Admin); err != nil {
		return nil, err
	}
	return appointment, nil
}

// Delete deletes a user
func (u *App) Delete(c echo.Context, id string) error {
	if err := u.rbac.EnforceRole(c, model.SupportRole); err != nil {
		return err
	}

	old, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if err = u.rbac.EnforceUser(c, old.Admin); err != nil {
		return err
	}

	if old.Status != model.StatusConcludedAppointment && old.Status != model.StatusCanceled {
		return zaplog.ZLog(fmt.Errorf("Não é possível eliminar audiências activas. Cancele ou conclua primeiro."))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Appointment, err error) {
	if err = u.rbac.EnforceRole(c, model.SupportRole); err != nil {
		return
	}

	old, err := u.udb.View(u.db, r.ID)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if err := u.rbac.EnforceUser(c, old.Admin); err != nil {
		return nil, err
	}

	if r.Status != "" && old.Status != r.Status && !old.AllowedStatuses(r.Status) {
		err = fmt.Errorf("Não é possível passar de %s para %s", old.Status, r.Status)
		return
	}

	update := model.Appointment{
		Base:     model.Base{ID: r.ID},
		Date:     r.Date,
		Admin:    r.Admin,
		Comments: r.Comments,
		Status:   r.Status,
	}

	var operator model.User
	if err = u.db.Model(&model.User{}).Where("uuid = ?", r.Admin).First(&operator).Error; err == nil {
		update.AdminName = operator.Name
	}
	if err = u.udb.Update(u.db, &update); err != nil {
		return
	}
	return u.udb.View(u.db, r.ID)
}
