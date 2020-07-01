package organization

import (
	"github.com/figassis/goinagbe/pkg/api/organization/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

type (
	//Update represents the update structure
	Update struct {
		ID        string
		Name      string
		Status    string
		Type      string
		Country   string
		Province  string
		City      string
		Logo      string
		Phone     string
		Email     string
		Website   string
		Facebook  string
		Twitter   string
		Linkedin  string
		Address   string
		Documents []string
	}

	//Create represents the update structure
	Create struct {
		Name      string
		Status    string
		Type      string
		Country   string
		Province  string
		City      string
		Phone     string
		Email     string
		Website   string
		Facebook  string
		Twitter   string
		Linkedin  string
		Address   string
		Logo      string
		Documents []string
	}

	//Service represents the HTTP service interface
	Service interface {
		Create(echo.Context, *Create) (*model.Organization, error)
		List(echo.Context, *model.Pagination) ([]model.Organization, string, string, int64, int64, error)
		View(echo.Context, string) (*model.Organization, error)
		Delete(echo.Context, string) error
		Update(echo.Context, *Update) (*model.Organization, error)
	}

	//App represents the application service
	App struct {
		db   *gorm.DB
		udb  UDB
		rbac RBAC
		sec  Securer
	}

	// Securer represents security interface
	Securer interface {
		Hash(string) string
	}

	// UDB represents the repository interface
	UDB interface {
		Create(*gorm.DB, model.Organization) (*model.Organization, error)
		View(*gorm.DB, string) (*model.Organization, error)
		List(*gorm.DB, *model.ListQuery, *model.Pagination) ([]model.Organization, string, string, int64, int64, error)
		Update(*gorm.DB, *model.Organization) error
		Delete(*gorm.DB, string) error
	}

	// RBAC represents role-based-access-control interface
	RBAC interface {
		User(echo.Context) *model.AuthUser
		EnforceUser(echo.Context, string) error
		EnforceRole(echo.Context, model.AccessRole) error
		EnforceCompany(echo.Context, string) error
		IsLowerRole(echo.Context, model.AccessRole) error
	}
)

// New creates new user application service
func New(db *gorm.DB, udb UDB, rbac RBAC, sec Securer) *App {
	return &App{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes Organization application service with defaults
func Initialize(db *gorm.DB, app *config.Application, rbac RBAC, sec Securer) (u *App) {
	u = New(db, sql.New(), rbac, sec)
	return
}
