package sql

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/labstack/echo"
)

// NewUser returns a new user database instance
func New() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user on database
func (u *User) Create(db *gorm.DB, usr model.User) (user *model.User, err error) {

	if err = db.Create(&usr).Error; err != nil {
		err = zaplog.ZLog(err)
		return nil, err
	}

	return &usr, nil
}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id string) (user *model.User, err error) {
	user = new(model.User)
	if err = zaplog.ZLog(db.Where("uuid = ?", id).First(user).Error); err != nil {
		return
	}

	return
}

// Update updates user's contact info
func (u *User) Update(db *gorm.DB, user *model.User) (*model.User, error) {
	if err := zaplog.ZLog(db.Model(user).Where("uuid = ?", user.ID).Updates(*user).Error); err != nil {
		return nil, err
	}

	return user, nil
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db *gorm.DB, qp *model.ListQuery, p *model.Pagination) (users []model.User, next, prev string, total, pages int64, err error) {
	q := db.Model(&model.User{})
	if p.UserQuery != nil {
		q = q.Where(*p.UserQuery)
	}
	if p.SearchQuery != "" {
		str := fmt.Sprintf("%%%s%%", p.SearchQuery)
		q = q.Where("name like ? OR email like ? OR phone like ?", str, str, str)
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
	err = zaplog.ZLog(q.Order("id DESC").Order("name DESC").Limit(limit).Find(&users).Error)
	return
}

// Delete sets deleted_at for a user
func (u *User) Delete(db *gorm.DB, id string) (err error) {

	if err = db.Where("uuid = ?", id).Delete(&model.User{}).Error; err != nil {
		return
	}

	return
}

func (u *User) FindByUsername(db *gorm.DB, uname string) (user *model.User, err error) {
	user = new(model.User)
	err = db.Where("username = ?", uname).First(user).Error
	err = zaplog.ZLog(err)
	return
}
