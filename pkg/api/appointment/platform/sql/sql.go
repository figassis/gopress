package sql

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
)

// New returns a new user database instance
func New() *ORM {
	return &ORM{}
}

// DB represents the client for user table
type ORM struct{}

// Create creates a new user on database
func (u *ORM) Create(db *gorm.DB, usr model.Appointment) (user *model.Appointment, err error) {

	if err = db.Create(&usr).Error; err != nil {
		err = zaplog.ZLog(err)
		return nil, err
	}

	return &usr, nil
}

// View returns single user by ID
func (u *ORM) View(db *gorm.DB, id string) (user *model.Appointment, err error) {
	user = new(model.Appointment)
	err = db.Where("uuid = ?", id).First(user).Error
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	return
}

// Update updates user's contact info
func (u *ORM) Update(db *gorm.DB, user *model.Appointment) (err error) {
	if err = zaplog.ZLog(db.Model(user).Where("uuid = ?", user.ID).Updates(*user).Error); err != nil {
		return
	}

	return
}

// List returns list of all users retrievable for the current user, depending on role
func (u *ORM) List(db *gorm.DB, qp *model.ListQuery, p *model.Pagination) (users []model.Appointment, next, prev string, total, pages int64, err error) {
	q := db.Model(&model.Appointment{})
	if p.AppointmentQuery != nil {
		q = q.Where(*p.AppointmentQuery)
	}
	if p.SearchQuery != "" {
		str := fmt.Sprintf("%%%s%%", p.SearchQuery)
		q = q.Where("user_name like ? OR admin_name like ? OR contact_name like ? OR contact_number like ? OR contact_email like ?", str, str, str, str, str)
	}
	if qp != nil {
		q = q.Where(qp.Query, qp.ID)
	}

	zaplog.ZLog(q.Count(&total).Error)
	limit, cursor, prev, next := p.DbPagination(q)

	if cursor.ID != "" {
		q = q.Where("(id,uuid) <= (?,?)", cursor.IntID, cursor.ID)
	}

	pages = total / int64(limit)
	err = zaplog.ZLog(q.Order("id DESC").Limit(limit).Find(&users).Error)
	return
}

// Delete sets deleted_at for a user
func (u *ORM) Delete(db *gorm.DB, id string) (err error) {
	return zaplog.ZLog(db.Where("uuid = ?", id).Delete(&model.Appointment{}).Error)
}
