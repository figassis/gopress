package sql

import (
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
func (u *ORM) Create(db *gorm.DB, usr model.File) (user *model.File, err error) {

	if err = db.Create(&usr).Error; err != nil {
		err = zaplog.ZLog(err)
		return nil, err
	}

	switch usr.Resource {
	case model.ResourcePost, model.ResourceScholarship:
		//Upload immediately to make them available in the Editor
		url, err2 := model.SaveFile(usr.Path, model.GetS3KeyFromID(usr.ID, usr.Public))
		if err = zaplog.ZLog(err2); err != nil {
			return nil, err
		}

		zaplog.ZLog(db.Model(&model.File{}).Where("uuid = ?", usr.ID).Update("url", url).Error)
	default:
		if err = db.Create(&model.Upload{Base: model.Base{ID: usr.ID}, Path: usr.Path, Public: usr.Public, Status: model.StatusPending}).Error; err != nil {
			err = zaplog.ZLog(err)
			return nil, err
		}
	}

	return &usr, nil
}

// View returns single user by ID
func (u *ORM) View(db *gorm.DB, id string) (user *model.File, err error) {
	user = new(model.File)
	if err = zaplog.ZLog(db.Where("uuid = ?", id).First(user).Error); err != nil {
		return nil, err
	}

	user.URL, err = user.GetURL()

	return
}

// Update updates user's contact info
func (u *ORM) Update(db *gorm.DB, user *model.File) (err error) {
	if err = zaplog.ZLog(db.Model(user).Where("uuid = ?", user.ID).Updates(*user).Error); err != nil {
		return
	}

	return
}

// List returns list of all users retrievable for the current user, depending on role
func (u *ORM) List(db *gorm.DB, qp *model.ListQuery, p *model.Pagination) (users []model.File, next, prev string, total, pages int64, err error) {
	q := db.Model(&model.File{})
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
	return zaplog.ZLog(db.Where("uuid = ?", id).Delete(&model.File{}).Error)
}
