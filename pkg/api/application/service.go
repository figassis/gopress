package application

import (
	"time"

	"github.com/figassis/goinagbe/pkg/api/application/platform/sql"
	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type (
	//Update represents the update structure
	Update struct {
		ID   string
		Name string
		// BirthDate         time.Time
		Gender            string
		BirthProvince     string
		BirthCity         string
		CurrentCountry    string
		CurrentProvince   string
		CurrentCity       string
		Phone             string
		Email             string
		Skype             string
		BankName          string
		BankAccountNumber string
		BankAccountOwner  string
		IBAN              string
		RPE               bool
		Operator          string
		Score             decimal.Decimal
		FamilyMembers     int64
		Children          int64
		WorkingMembers    int64

		EducationLevel            string
		EducationSchool           string
		EducationSchoolDepartment string
		EducationSchoolContact    string
		EducationCourse           string
		EducationCountry          string
		EducationProvince         string
		// GraduationDate            time.Time
		EducationGrade      decimal.Decimal
		EducationEvaluation string

		ApplicationLevel string
		CourseSchool     string
		CourseDepartment string
		CourseOption1    string
		CourseOption2    string
		CourseOption3    string
		CourseCountry    string
		CourseProvince   string
		CourseContact    string
		CourseYear       int64
		CourseAverage    decimal.Decimal

		Employer                 string
		EmployerAddress          string
		EmployerContact          string
		JobTitle                 string
		Salary                   decimal.Decimal
		HasEmployerAuthorization bool
		Status                   string
		Languages                model.Languages
		Family                   []model.FamilyInfo
		Documents                model.List
	}

	//Create represents the update structure
	Create struct {
		Scholarship        string
		UserID             string
		IDNumber           string
		PassportNumber     string
		Name               string
		BirthDate          time.Time
		Gender             string
		BirthProvince      string
		BirthCity          string
		CurrentCountry     string
		CurrentProvince    string
		CurrentCity        string
		Phone              string
		Email              string
		Skype              string
		BankName           string
		BankAccountNumber  string
		BankAccountOwner   string
		IBAN               string
		RPE                bool
		Operator           string
		Score              decimal.Decimal
		ScoreGrades        decimal.Decimal
		ScoreAge           decimal.Decimal
		ScoreJustification string
		FamilyMembers      int64
		Children           int64
		WorkingMembers     int64

		EducationLevel            string
		EducationSchool           string
		EducationSchoolDepartment string
		EducationSchoolContact    string
		EducationCourse           string
		EducationCourseID         string
		EducationCountry          string
		EducationProvince         string
		GraduationDate            time.Time
		EducationGrade            decimal.Decimal
		EducationEvaluation       string

		ApplicationLevel    string
		CourseSchool        string
		CourseDepartment    string
		ApplicationCourseID string
		CourseOption1       string
		CourseOption2       string
		CourseOption3       string
		CourseCountry       string
		CourseProvince      string
		CourseContact       string
		CourseYear          int64
		CourseAverage       decimal.Decimal

		Employer                 string
		EmployerAddress          string
		EmployerContact          string
		JobTitle                 string
		Salary                   decimal.Decimal
		HasEmployerAuthorization bool
		Status                   string
		Languages                model.Languages
		Family                   []model.FamilyInfo
		Documents                []string
	}

	//Service represents the HTTP service interface
	Service interface {
		Create(echo.Context, *Create) (*model.Application, error)
		List(echo.Context, *model.Pagination) ([]model.Application, string, string, int64, int64, error)
		View(echo.Context, string) (*model.Application, error)
		Delete(echo.Context, string) error
		Update(echo.Context, *Update) (*model.Application, error)
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
		Create(*gorm.DB, model.Application) (*model.Application, error)
		View(*gorm.DB, string) (*model.Application, error)
		List(*gorm.DB, *model.ListQuery, *model.Pagination) ([]model.Application, string, string, int64, int64, error)
		Update(*gorm.DB, *model.Application) error
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

// Initialize initalizes Application application service with defaults
func Initialize(db *gorm.DB, app *config.Application, rbac RBAC, sec Securer) (u *App) {
	u = New(db, sql.New(), rbac, sec)
	return
}
