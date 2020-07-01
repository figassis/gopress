package sql

import (
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id string) (user *model.User, err error) {
	user = new(model.User)
	err = db.Model(&model.User{}).Where("uuid = ?", id).First(user).Error
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates user's info
func (u *User) Update(db *gorm.DB, user *model.User) (err error) {
	if err = zaplog.ZLog(db.Model(user).Updates(*user).Error); err != nil {
		return
	}

	return
}
