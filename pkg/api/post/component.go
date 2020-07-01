// Package user contains user application services
package post

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Post, error) {
	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	var operator model.User
	if err = u.db.Model(&model.User{}).Where("uuid = ?", req.Author).First(&operator).Error; err != nil {
		return nil, err
	}

	var dupe model.Post
	if err = u.db.Model(&model.Post{}).Where("slug = ?", req.Slug).Order("id DESC").First(&dupe).Error; err == nil {
		fragment := strings.TrimPrefix(dupe.Slug, req.Slug)
		if fragment == "" {
			req.Slug += "-2"
		}

		counter, err := strconv.Atoi(strings.TrimPrefix(fragment, "-"))
		if err != nil {
			req.Slug += "-2"
		} else {
			req.Slug += fmt.Sprintf("-%d", counter+1)
		}
	}

	if len(req.Excerpt) > 255 {
		req.Excerpt = req.Excerpt[:250] + "..."
	}

	post := model.Post{
		Base:       model.Base{ID: id},
		Author:     req.Author,
		AuthorName: operator.Name,
		Category:   req.Category,
		Tags:       req.Tags,
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		Status:     req.Status,
	}
	return u.udb.Create(u.db, post)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Post, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourcePost)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.PostQuery = &model.Post{}
		params := c.QueryParams()

		p.PostQuery.Status = params.Get("status")
		p.SearchQuery = params.Get("s")
	}

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Post, error) {
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

	post, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if post.Status != model.StatusDraft {
		return zaplog.ZLog(errors.New("Apenas é possível eliminar artigos em rascunho"))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Post, err error) {
	if err = u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		zaplog.ZLog(err)
		return
	}

	if len(r.Excerpt) > 255 {
		r.Excerpt = r.Excerpt[:250] + "..."
	}

	old, err := u.udb.View(u.db, r.ID)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if r.Status != "" && old.Status != r.Status && !old.AllowedStatuses(r.Status) {
		err = zaplog.ZLog(fmt.Errorf("Não é possível passar de %s para %s", old.Status, r.Status))
		return
	}

	update := model.Post{
		Base:     model.Base{ID: r.ID},
		Author:   r.Author,
		Category: r.Category,
		Tags:     r.Tags,
		Title:    r.Title,
		Slug:     r.Slug,
		Content:  r.Content,
		Excerpt:  r.Excerpt,
		Status:   r.Status,
	}

	var operator model.User
	if err = u.db.Model(&model.User{}).Where("uuid = ?", r.Author).First(&operator).Error; err == nil {
		update.AuthorName = operator.Name
	}

	if err = zaplog.ZLog(u.udb.Update(u.db, &update)); err != nil {
		return
	}
	return u.udb.View(u.db, r.ID)
}
