package rbac

import (
	"github.com/figassis/goinagbe/pkg/utl/model"
	echo "github.com/labstack/echo/v4"
)

// New creates new RBAC service
func New() *Service {
	return &Service{}
}

// Service is RBAC application service
type Service struct{}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

// User returns user data stored in jwt token
func (s *Service) User(c echo.Context) *model.AuthUser {
	id := c.Get("id").(string)
	organization := c.Get("organization").(string)
	user := c.Get("username").(string)
	email := c.Get("email").(string)
	role := c.Get("role").(model.AccessRole)
	return &model.AuthUser{
		ID:           id,
		Username:     user,
		Organization: organization,
		Email:        email,
		Role:         role,
	}
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r model.AccessRole) error {
	return checkBool(!(c.Get("role").(model.AccessRole) > r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID string) error {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(string) == ID)
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *Service) EnforceCompany(c echo.Context, ID string) error {
	if s.isAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, model.CompanyAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("organization").(string) == ID)
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(model.AccessRole) > model.AdminRole)
}

func (s *Service) isCompanyAdmin(c echo.Context) bool {
	// Must query company ID in database for the given user
	return !(c.Get("role").(model.AccessRole) > model.CompanyAdminRole)
}

// AccountCreate performs auth check when creating a new account
// Location admin cannot create accounts, needs to be fixed on EnforceLocation function
func (s *Service) AccountCreate(c echo.Context, role model.AccessRole, organization string) error {
	if err := s.IsLowerRole(c, role); err != nil {
		return err
	}

	if err := s.EnforceCompany(c, organization); err != nil {
		return err
	}

	return nil
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r model.AccessRole) error {
	return checkBool(c.Get("role").(model.AccessRole) < r || c.Get("role").(model.AccessRole) == model.SuperAdminRole)
}
