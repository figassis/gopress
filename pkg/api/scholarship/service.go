package scholarship

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/scholarship/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type (
	//Update represents the update structure
	Update struct {
		ID                string
		Name              string
		Sponsor           string
		Start             time.Time
		End               time.Time
		Available         int64
		MaxAge            int64
		MinGrade          int64
		RPEQuota          decimal.Decimal
		PriorityQuota     decimal.Decimal
		ProvinceQuota     model.ProvinceQuota
		Status            string
		Content           string
		Type              string
		Level             string
		EnableQuotas      bool
		Documents         []string
		RequiredDocuments []string
	}

	//Create represents the update structure
	Create struct {
		Name              string
		Sponsor           string
		Start             time.Time
		End               time.Time
		Available         int64
		MaxAge            int64
		MinGrade          int64
		RPEQuota          decimal.Decimal
		PriorityQuota     decimal.Decimal
		ProvinceQuota     model.ProvinceQuota
		EnableQuotas      bool
		Type              string
		Level             string
		Content           string
		Status            string
		Documents         []string
		RequiredDocuments []string
	}

	//Service represents the HTTP service interface
	Service interface {
		Create(echo.Context, *Create) (*model.Scholarship, error)
		List(echo.Context, *model.Pagination) ([]model.Scholarship, string, string, int64, int64, error)
		View(echo.Context, string) (*model.Scholarship, error)
		Export(echo.Context, string) (string, error)
		Delete(echo.Context, string) error
		Update(echo.Context, *Update) (*model.Scholarship, error)
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
		Create(*gorm.DB, model.Scholarship) (*model.Scholarship, error)
		View(*gorm.DB, string) (*model.Scholarship, error)
		List(*gorm.DB, *model.ListQuery, *model.Pagination) ([]model.Scholarship, string, string, int64, int64, error)
		Update(*gorm.DB, *model.Scholarship) error
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

// Initialize initalizes Scholarship application service with defaults
func Initialize(db *gorm.DB, app *config.Application, rbac RBAC, sec Securer) (u *App) {
	u = New(db, sql.New(), rbac, sec)
	return
}
