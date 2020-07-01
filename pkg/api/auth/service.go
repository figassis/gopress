package auth

import (
	"github.com/figassis/goinagbe/pkg/api/auth/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

// New creates new iam service
func New(db *gorm.DB, udb UserDB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *gorm.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, sql.New(), j, sec, rbac)
}

// Service represents auth service interface
type (
	Credentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	Reset struct {
		Email string `json:"email" validate:"required"`
	}

	CompleteReset struct {
		Token              string `validate:"required"`
		NewPassword        string `validate:"required,min=8"`
		NewPasswordConfirm string `validate:"required"`
	}

	Signup struct {
		Name            string
		Password        string
		PasswordConfirm string
		Email           string
		Honeypot        string
		Phone           string
	}
	Service interface {
		Authenticate(echo.Context, string, string) (*model.AuthToken, error)
		Signup(echo.Context, Signup) error
		Contact(echo.Context, model.Contact) error
		Refresh(echo.Context) (*model.AuthToken, error)
		Resend(echo.Context) error
		Me(echo.Context) (*model.User, error)
		GetPublicData(echo.Context) (*model.Public, error)
		Unsubscribe(echo.Context, string) error
		ConfirmEmail(echo.Context, string) error
		Bounce(echo.Context, model.BounceNotification) error
		Reset(echo.Context, string) error
		CheckResetToken(echo.Context, string) error
		CompleteReset(echo.Context, string, string) error
	}

	// Auth represents auth application service
	Auth struct {
		db   *gorm.DB
		udb  UserDB
		tg   TokenGenerator
		sec  Securer
		rbac RBAC
	}

	// UserDB represents user repository interface
	UserDB interface {
		View(*gorm.DB, string) (*model.User, error)
		FindByUsername(*gorm.DB, string) (*model.User, error)
		FindByToken(*gorm.DB, string) (*model.User, error)
		GetPublicData(*gorm.DB) (*model.Public, error)
		Update(*gorm.DB, *model.User) error
		Signup(*gorm.DB, *model.User) error
		Unsubscribe(*gorm.DB, string) error
		ConfirmEmail(*gorm.DB, string) error
		Bounce(*gorm.DB, model.BounceNotification) error
	}

	// TokenGenerator represents token generator (jwt) interface
	TokenGenerator interface {
		GenerateToken(*model.User) (string, string, string, string, error)
		// ParseToken(echo.Context) (*jwt.Token, error)
	}

	// Securer represents security interface
	Securer interface {
		HashMatchesPassword(string, string) bool
		Token(string) string
		Hash(string) string
		Password(string, ...string) bool
	}

	// RBAC represents role-based-access-control interface
	RBAC interface {
		User(echo.Context) *model.AuthUser
	}
)
