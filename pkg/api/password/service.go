package password

import (
	"github.com/figassis/goinagbe/pkg/api/password/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

type (
	Reset struct {
		Email string `json:"email" validate:"required"`
	}

	Update struct {
		OldPassword        string `validate:"required"`
		NewPassword        string `validate:"required,min=8"`
		NewPasswordConfirm string `validate:"required"`
	}

	// Service represents password application interface
	Service interface {
		Update(echo.Context, string, string, string) error
	}

	// Password represents password application service
	Password struct {
		db   *gorm.DB
		udb  UserDB
		rbac RBAC
		sec  Securer
	}

	// UserDB represents user repository interface
	UserDB interface {
		View(*gorm.DB, string) (*model.User, error)
		Update(*gorm.DB, *model.User) error
	}

	// Securer represents security interface
	Securer interface {
		Hash(string) string
		HashMatchesPassword(string, string) bool
		Password(string, ...string) bool
	}

	// RBAC represents role-based-access-control interface
	RBAC interface {
		EnforceUser(echo.Context, string) error
		User(echo.Context) *model.AuthUser
	}
)

// New creates new password application service
func New(db *gorm.DB, udb UserDB, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		udb:  udb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initalizes password application service with defaults
func Initialize(db *gorm.DB, rbac RBAC, sec Securer) *Password {
	return New(db, sql.NewUser(), rbac, sec)
}
