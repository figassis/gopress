// Package user contains user application services
package course

import (
	"fmt"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Course, error) {
	if err := u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return nil, err
	}

	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	schoolName := ""
	var school model.Organization
	if err := u.db.Model(&model.Organization{}).Where("uuid = ?", req.School).First(&school).Error; err == nil {
		schoolName = school.Name
	}

	course := model.Course{
		Base:       model.Base{ID: id},
		Name:       req.Name,
		School:     req.School,
		SchoolName: schoolName,
		Department: req.Department,
		Domain:     req.Domain,
		Cluster:    req.Cluster,
		Type:       req.Type,
		Level:      req.Level,
	}

	var domain model.CourseDomain
	if err := u.db.Model(&model.CourseDomain{}).Where("uuid = ?", req.Domain).First(&domain).Error; err == nil {
		course.DomainName = domain.Name
		for _, cluster := range domain.Clusters {
			if course.Cluster == cluster.ID {
				course.ClusterName = cluster.Name
			}
		}
	}

	return u.udb.Create(u.db, course)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Course, string, string, int64, int64, error) {
	au := u.rbac.User(c)
	q, err := query.List(au, model.ResourceCourse)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.CourseQuery = &model.Course{}
		params := c.QueryParams()

		p.CourseQuery.Level = params.Get("level")
		p.CourseQuery.Type = params.Get("type")
		p.CourseQuery.Domain = params.Get("domain")
		p.CourseQuery.Cluster = params.Get("cluster")
		p.CourseQuery.School = params.Get("school")
		p.CourseQuery.Department = params.Get("department")
		p.SearchQuery = params.Get("s")
	}

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Course, error) {
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

	var applications int
	if err := u.db.Model(&model.Application{}).Where("application_course_id = ?", id).Count(&applications).Error; err == nil && applications > 0 {
		return zaplog.ZLog(fmt.Errorf("Existem %d candidaturas associadas a este curso. Não é possível eliminar", applications))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Course, err error) {
	if err = u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return
	}

	schoolName := ""
	var school model.Organization
	if err := u.db.Model(&model.Organization{}).Where("uuid = ?", r.School).First(&school).Error; err == nil {
		schoolName = school.Name
	}

	update := model.Course{
		Base:       model.Base{ID: r.ID},
		Name:       r.Name,
		Domain:     r.Domain,
		Cluster:    r.Cluster,
		Type:       r.Type,
		Level:      r.Level,
		School:     r.School,
		SchoolName: schoolName,
	}

	if err = u.udb.Update(u.db, &update); err != nil {
		return
	}
	return u.udb.View(u.db, r.ID)
}
