package appointment

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/appointment/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

type (
	//Update represents the update structure
	Update struct {
		ID       string
		Date     time.Time
		Admin    string
		Comments string
		Status   string
	}

	//Create represents the update structure
	Create struct {
		Resource      string
		ResourceID    string
		ContactName   string
		ContactNumber string
		ContactEmail  string
		Message       string
		Date          time.Time
		Admin         string
	}

	//Service represents the HTTP service interface
	Service interface {
		Create(echo.Context, *Create) (*model.Appointment, error)
		List(echo.Context, *model.Pagination) ([]model.Appointment, string, string, int64, int64, error)
		View(echo.Context, string) (*model.Appointment, error)
		Delete(echo.Context, string) error
		Update(echo.Context, *Update) (*model.Appointment, error)
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
		Create(*gorm.DB, model.Appointment) (*model.Appointment, error)
		View(*gorm.DB, string) (*model.Appointment, error)
		List(*gorm.DB, *model.ListQuery, *model.Pagination) ([]model.Appointment, string, string, int64, int64, error)
		Update(*gorm.DB, *model.Appointment) error
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

// Initialize initalizes Appointment application service with defaults
func Initialize(db *gorm.DB, app *config.Application, rbac RBAC, sec Securer) (u *App) {
	u = New(db, sql.New(), rbac, sec)
	return
}
