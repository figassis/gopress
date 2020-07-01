package query

import (
	"fmt"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/labstack/echo"
)

// List prepares data for list queries
func List(u *model.AuthUser, resource string) (*model.ListQuery, error) {
	//The smaller the role, the highter the level
	if u.Role <= model.SupportRole {
		return nil, nil
	}

	roleFilter := ""
	if u.Role > model.SuperAdminRole {
		roleFilter = fmt.Sprintf(" and role > %d", model.SuperAdminRole)
	}

	switch resource {
	case model.ResourceUser:
		if u.Role > model.CompanyAdminRole {
			return nil, echo.ErrForbidden
		}
		return &model.ListQuery{Query: "organization = ?" + roleFilter, ID: u.Organization}, nil
	case model.ResourcePost:
		return &model.ListQuery{Query: "status = ?", ID: model.StatusPublished}, nil
	case model.ResourceScholarship:
		return &model.ListQuery{Query: "status = ?", ID: model.StatusOpen}, nil
	case model.ResourceApplication:
		return &model.ListQuery{Query: "user_id = ?", ID: u.ID}, nil
	case model.ResourceFile:
		return &model.ListQuery{Query: "public = true OR user_id = ?", ID: u.ID}, nil
	case model.ResourceAppointment:
		return &model.ListQuery{Query: "admin = ?", ID: u.ID}, nil
	case model.ResourceCourse:
		return nil, nil
	default:
		return nil, echo.ErrForbidden
	}
}
