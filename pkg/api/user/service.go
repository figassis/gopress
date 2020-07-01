package user

import (
	"github.com/brianvoe/gofakeit"
	"github.com/figassis/goinagbe/pkg/api/user/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

type (
	Update struct {
		ID           string
		Name         string
		Phone        string
		Status       string
		Role         model.AccessRole
		Unsubscribed bool
	}

	Create struct {
		Name            string
		Password        string
		PasswordConfirm string
		Email           string
		Phone           string
		Role            model.AccessRole
		Organization    string
	}

	// Service represents user application interface
	Service interface {
		Create(echo.Context, Create) (*model.User, error)
		List(echo.Context, *model.Pagination) ([]model.User, string, string, int64, int64, error)
		View(echo.Context, string) (*model.User, error)
		Delete(echo.Context, string) error
		Update(echo.Context, *Update) (*model.User, error)
	}

	// User represents user application service
	User struct {
		db   *gorm.DB
		udb  UDB
		rbac RBAC
		sec  Securer
	}

	// Securer represents security interface
	Securer interface {
		Hash(string) string
	}

	// UDB represents user repository interface
	UDB interface {
		Create(*gorm.DB, model.User) (*model.User, error)
		FindByUsername(*gorm.DB, string) (*model.User, error)
		View(*gorm.DB, string) (*model.User, error)
		List(*gorm.DB, *model.ListQuery, *model.Pagination) ([]model.User, string, string, int64, int64, error)
		Update(*gorm.DB, *model.User) (*model.User, error)
		Delete(*gorm.DB, string) error
	}

	// RBAC represents role-based-access-control interface
	RBAC interface {
		User(echo.Context) *model.AuthUser
		EnforceUser(echo.Context, string) error
		AccountCreate(echo.Context, model.AccessRole, string) error
		EnforceRole(echo.Context, model.AccessRole) error
		EnforceCompany(echo.Context, string) error
		IsLowerRole(echo.Context, model.AccessRole) error
	}
)

// New creates new user application service
func New(db *gorm.DB, udb UDB, rbac RBAC, sec Securer) *User {

	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *gorm.DB, app *config.Application, rbac RBAC, sec Securer) (u *User) {
	u = New(db, sql.New(), rbac, sec)
	initializeAdminUser(db, app, u)
	initializeSystemUser(db, app, u)
	return
}

func initializeAdminUser(db *gorm.DB, app *config.Application, u *User) (err error) {
	var user model.User
	if err = db.Model(&model.User{}).Where("email = ?", app.AdminEmail).First(&user).Error; user.ID != "" {
		return nil
	}

	ids, err := util.GenerateUUIDS(2)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var org model.Organization
	if err = db.Model(&model.Organization{}).Where("type = ?", model.OrgMain).First(&org).Error; err != nil {
		org = model.Organization{
			Base:     model.Base{ID: ids[0]},
			Name:     app.Company,
			Status:   model.StatusActive,
			Type:     model.OrgMain,
			Phone:    gofakeit.Phone(),
			Email:    "inagbe@inagbeonline.com",
			Country:  "Angola",
			Province: "Luanda",
			City:     "Luanda",
			Logo:     "assets/images/main.jpg",
		}

		if err = zaplog.ZLog(u.db.Create(&org).Error); err != nil {
			return
		}
	}

	user = model.User{
		Username:         app.AdminEmail,
		Password:         u.sec.Hash(app.AdminPassword),
		Email:            app.AdminEmail,
		Name:             "Admin",
		Phone:            gofakeit.Phone(),
		Role:             model.SuperAdminRole,
		Status:           model.StatusActive,
		Organization:     ids[0],
		OrganizationName: "INAGBE",
		UnsubscribeID:    ids[1],
	}

	_, err = u.udb.Create(db, user)

	return
}

func initializeSystemUser(db *gorm.DB, app *config.Application, u *User) (err error) {
	var user model.User
	if err = db.Model(&model.User{}).Where("uuid = ?", model.AdminUser).First(&user).Error; err == nil && user.ID != "" {
		return nil
	}

	var org model.Organization
	if err = zaplog.ZLog(db.Model(&model.Organization{}).Where("type = ?", model.OrgMain).First(&org).Error); err != nil {
		return
	}

	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	user = model.User{
		Base:             model.Base{ID: model.AdminUser},
		Username:         model.AdminEmail,
		Password:         u.sec.Hash(util.ShortID()),
		Email:            model.AdminEmail,
		Phone:            gofakeit.Phone(),
		Name:             "System",
		Role:             model.SuperAdminRole,
		OrganizationName: "Nellcorp",
		Status:           model.StatusActive,
		Organization:     org.ID,
		UnsubscribeID:    id,
	}

	_, err = u.udb.Create(db, user)

	return
}
