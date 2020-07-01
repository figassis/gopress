// Package user contains user application services
package media

import (
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.File, error) {
	au := u.rbac.User(c)
	id, err := util.GenerateUUID()
	if err != nil {
		return nil, err
	}

	status := model.StatusPending
	if au.Role <= model.SupportRole {
		status = model.StatusApproved
	}

	if req.Resource == "" {
		req.Resource = model.ResourceUser
		req.ResourceID = au.ID
	}

	if req.Type == "" {
		req.Type = model.GeneralFile
	}

	var operator model.User
	if err = u.db.Model(&model.User{}).Where("uuid = ?", au.ID).First(&operator).Error; err != nil {
		return nil, err
	}

	media := model.File{
		Base:       model.Base{ID: id},
		UserID:     au.ID,
		UserName:   operator.Name,
		Name:       req.Name,
		Resource:   req.Resource,
		ResourceID: req.ResourceID,
		Type:       req.Type,
		URL:        model.GetS3URLFromID(id, req.Public),
		Status:     status,
		Public:     req.Public,
		Location:   "s3",
		Path:       req.Path,
	}
	return u.udb.Create(u.db, media)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.File, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourceFile)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.File, error) {
	file, err := u.udb.View(u.db, id)
	if err != nil {
		return nil, err
	}

	if file.Type == model.ResourceApplication {
		if err = u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
			return nil, err
		}
	}

	if err := u.rbac.EnforceUser(c, file.UserID); err != nil {
		return nil, err
	}
	return file, nil
}

// Delete deletes a user
func (u *App) Delete(c echo.Context, id string) error {
	file, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}

	if err := u.rbac.EnforceUser(c, file.UserID); err != nil {
		return err
	}

	if file.Type == model.ResourceApplication {
		if err = u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
			return err
		}
	}

	go model.DeleteFiles(&[]model.File{*file})
	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.File, err error) {
	file, err := u.udb.View(u.db, r.ID)
	if err != nil {
		return nil, err
	}

	if err := u.rbac.EnforceUser(c, file.UserID); err != nil {
		return nil, err
	}

	if file.Type == model.ResourceApplication {
		if err = u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
			return
		}
	}

	file.Comment = r.Comment
	file.Status = r.Status
	// file.Public = r.Public

	if err = u.udb.Update(u.db, file); err != nil {
		return
	}
	return u.udb.View(u.db, r.ID)
}
