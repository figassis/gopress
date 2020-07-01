package model

import (
	"fmt"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (a Course) Validate(db *gorm.DB) {
	if a.Name == "" || a.Domain == "" || a.Cluster == "" || a.DomainName == "" || a.ClusterName == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}

func (a CourseDomain) Validate(db *gorm.DB) {
	if a.Name == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}

func (a CourseCluster) Validate(db *gorm.DB) {
	if a.Name == "" || a.Domain == "" || a.DomainName == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}

func (a Organization) Validate(db *gorm.DB) {
	if a.Name == "" || a.Phone == "" || a.Email == "" || a.Logo == "" {
		db.AddError(fmt.Errorf("Campos em falta: %s", zaplog.JSON(a)))
	}
}

func (a Scholarship) Validate(db *gorm.DB) {
	now := time.Now()
	if a.Name == "" || a.Sponsor == "" || a.SponsorName == "" {
		db.AddError(fmt.Errorf("Campos em falta: %s", zaplog.JSON(a)))
	}

	if len(a.RequiredDocuments) == 0 {
		db.AddError(fmt.Errorf("A bolsa %s não tem documentos obrigatórios", zaplog.JSON(a)))
	}

	if a.Available < 0 {
		db.AddError(fmt.Errorf("A bolsa tem vagas negativas"))
	}

	if a.MaxAge < 18 {
		db.AddError(fmt.Errorf("Idade máxima inválida"))
	}

	if a.MinGrade >= 20 || a.MinGrade < 9 {
		db.AddError(fmt.Errorf("Nota mínima inválida"))
	}

	if a.Start.After(a.End) || a.End.Before(now) || a.End.Equal(now) {
		// db.AddError(fmt.Errorf("A bolsa tem datas inválidas"))
	}
}

func (a User) Validate(db *gorm.DB) {
	if a.Name == "" || a.Username == "" || a.Password == "" || a.Email == "" || a.Organization == "" || a.OrganizationName == "" || a.UnsubscribeID == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}

	if a.Role < 100 || a.Role > 130 {
		db.AddError(fmt.Errorf("Categoria inválida para utilizador"))
	}
}

func (a Application) Validate(db *gorm.DB) {
	if a.UserID == "" || a.IDNumber == "" || a.Scholarship == "" || a.ScholarshipName == "" || a.Name == "" || a.BirthProvince == "" ||
		a.BirthCity == "" || a.CurrentCountry == "" || a.CurrentProvince == "" || a.CurrentCity == "" || a.Phone == "" || a.Email == "" ||
		a.BankName == "" || a.BankAccountNumber == "" || a.BankAccountOwner == "" || a.IBAN == "" || a.EducationLevel == "" ||
		a.EducationSchool == "" || a.EducationCourse == "" || a.EducationCountry == "" || a.EducationProvince == "" ||
		a.ApplicationLevel == "" || a.CourseOption1 == "" {
		zaplog.ZLog(db.AddError(fmt.Errorf("A candidatura de %s tem campos em falta: %s", a.Email, zaplog.JSON(a))))
	}

	if a.EducationGrade.LessThan(decimal.NewFromFloat(0)) {
		zaplog.ZLog(db.AddError(fmt.Errorf("A candidatura de %s tem nota inválida: %v", a.Email, a.EducationGrade)))
	}

}

func (a ProvinceQuota) Validate(db *gorm.DB) {
	if a.Scholarship == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}

	if a.Bengo.LessThan(decimal.NewFromFloat(0)) || a.Bengo.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Benguela.LessThan(decimal.NewFromFloat(0)) || a.Benguela.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Bie.LessThan(decimal.NewFromFloat(0)) || a.Bie.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Cabinda.LessThan(decimal.NewFromFloat(0)) || a.Cabinda.GreaterThan(decimal.NewFromFloat(100)) ||
		a.CuandoCubango.LessThan(decimal.NewFromFloat(0)) || a.CuandoCubango.GreaterThan(decimal.NewFromFloat(100)) ||
		a.CuanzaNorte.LessThan(decimal.NewFromFloat(0)) || a.CuanzaNorte.GreaterThan(decimal.NewFromFloat(100)) ||
		a.CuanzaSul.LessThan(decimal.NewFromFloat(0)) || a.CuanzaSul.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Cunene.LessThan(decimal.NewFromFloat(0)) || a.Cunene.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Huambo.LessThan(decimal.NewFromFloat(0)) || a.Huambo.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Huila.LessThan(decimal.NewFromFloat(0)) || a.Huila.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Luanda.LessThan(decimal.NewFromFloat(0)) || a.Luanda.GreaterThan(decimal.NewFromFloat(100)) ||
		a.LundaNorte.LessThan(decimal.NewFromFloat(0)) || a.LundaNorte.GreaterThan(decimal.NewFromFloat(100)) ||
		a.LundaSul.LessThan(decimal.NewFromFloat(0)) || a.LundaSul.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Malanje.LessThan(decimal.NewFromFloat(0)) || a.Malanje.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Moxico.LessThan(decimal.NewFromFloat(0)) || a.Moxico.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Namibe.LessThan(decimal.NewFromFloat(0)) || a.Namibe.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Uige.LessThan(decimal.NewFromFloat(0)) || a.Uige.GreaterThan(decimal.NewFromFloat(100)) ||
		a.Zaire.LessThan(decimal.NewFromFloat(0)) || a.Zaire.GreaterThan(decimal.NewFromFloat(100)) {
		db.AddError(fmt.Errorf("Quotas inválidas"))
	}
}

func (a Statistic) Validate(db *gorm.DB) {
	if a.Name == "" || a.Resource == "" || a.ResourceID == "" || a.Type == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}

func (a Appointment) Validate(db *gorm.DB) {
	if a.User == "" || a.UserName == "" || a.Resource == "" || a.ContactEmail == "" || a.ContactName == "" || a.ContactNumber == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}

func (a Post) Validate(db *gorm.DB) {
	if a.Author == "" || a.AuthorName == "" || a.Title == "" || a.Slug == "" || a.Status == "" {
		db.AddError(fmt.Errorf("Campos em falta"))
	}
}
