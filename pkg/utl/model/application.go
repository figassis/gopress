package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/shopspring/decimal"
)

var (
	applicationStatusFlows = map[string]List{
		StatusPending:             List{StatusValidated, StatusRejectedApplication, StatusNeedsReview, StatusDuplicate},
		StatusValidated:           List{StatusNeedsReview, StatusApprovedApplication, StatusRejectedApplication, StatusDuplicate},
		StatusApprovedApplication: List{StatusNeedsReview, StatusAwarded, StatusDuplicate},
		StatusRejectedApplication: List{StatusNeedsReview},
		StatusNeedsReview:         List{StatusValidated, StatusRejectedApplication, StatusDuplicate, StatusApprovedApplication},
		StatusDuplicate:           List{StatusRejectedApplication},
		StatusAwarded:             List{StatusCanceledApplication},
		StatusCanceledApplication: List{},
	}
)

type (
	//To allow a single user to apply multiple times to the same scholarshipe (eg. parent applyign for kids)
	//We need to remove the user, email and phone indexes
	Application struct {
		Base
		UserID             string    `gorm:"index:user;not null"`
		IDNumber           string    `gorm:"index:id_number;not null"`
		PassportNumber     string    `gorm:"index:passport"`
		Scholarship        string    `gorm:"index:user,id_number,passport,email,phone;not null"`
		ScholarshipName    string    `gorm:"not null"`
		Name               string    `gorm:"not null"`
		BirthDate          time.Time `gorm:"not null"`
		Gender             string    `gorm:"type:ENUM('Masculino','Feminino','Não-Binário');default:'Masculino';not null"`
		BirthProvince      string    `gorm:"not null"`
		BirthCity          string    `gorm:"not null"`
		CurrentCountry     string    `gorm:"not null;default:'Angola'"`
		CurrentProvince    string    `gorm:"not null"`
		CurrentCity        string    `gorm:"not null"`
		Phone              string    `gorm:"not null;index:phone"`
		Email              string    `gorm:"not null;index:email"`
		Skype              string
		BankName           string `gorm:"not null"`
		BankAccountNumber  string `gorm:"not null"`
		BankAccountOwner   string `gorm:"not null"`
		IBAN               string `gorm:"not null"`
		RPE                bool
		Operator           string
		OperatorName       string
		Score              decimal.Decimal `sql:"type:decimal(5,2)" gorm:"not null;default:0"`
		ScoreGrades        decimal.Decimal `sql:"type:decimal(5,2)" gorm:"not null;default:0"`
		ScoreAge           decimal.Decimal `sql:"type:decimal(5,2)" gorm:"not null;default:0"`
		ScoreJustification string
		FamilyMembers      int64
		Children           int64
		WorkingMembers     int64

		EducationLevel            string `gorm:"type:ENUM('Ensino Médio','Graduação','Pós-Graduação');default:'Ensino Médio';not null"`
		EducationSchool           string `gorm:"not null"`
		EducationSchoolDepartment string
		EducationSchoolContact    string
		EducationCourse           string `gorm:"not null"`
		EducationCourseID         string
		EducationCountry          string          `gorm:"not null;default:'Angola'"`
		EducationProvince         string          `gorm:"not null"`
		GraduationDate            time.Time       `gorm:"not null"`
		EducationGrade            decimal.Decimal `sql:"type:decimal(4,2)" gorm:"not null;default:0"`
		EducationEvaluation       string          `gorm:"type:ENUM('Não Aplicável','Bom com distinção','Muito Bom','Excelente');default:'Não Aplicável';not null"`

		ApplicationLevel    string `gorm:"type:ENUM('Licenciatura','Mestrado','Especialidade','Doutoramento');default:'Licenciatura';not null"`
		CourseSchool        string `gorm:"not null"`
		CourseDepartment    string
		ApplicationCourseID string
		CourseOption1       string `gorm:"not null"`
		CourseOption2       string
		CourseOption3       string
		CourseCountry       string
		CourseProvince      string
		CourseContact       string
		CourseYear          int64
		CourseAverage       decimal.Decimal `sql:"type:decimal(4,2)" gorm:"not null;default:0"`

		Employer                 string
		EmployerAddress          string
		EmployerContact          string
		JobTitle                 string
		Salary                   decimal.Decimal `sql:"type:decimal(12,2)" gorm:"not null;default:0"`
		HasEmployerAuthorization bool
		Status                   string       `gorm:"type:ENUM('Pendente','Validada','Aprovada','Atribuída','Duplicada','Precisa de Revisão','Rejeitada');default:'Pendente';not null"`
		Documents                []File       `gorm:"foreignkey:ResourceID;association_foreignkey:uuid;association_autoupdate:false"`
		Languages                Languages    `gorm:"foreignkey:Application;association_foreignkey:uuid;association_autoupdate:false"`
		Family                   []FamilyInfo `gorm:"foreignkey:Application;association_foreignkey:uuid;association_autoupdate:false"`
		Unlock                   bool
	}

	Languages struct {
		Base
		Application string
		French      string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		Spanish     string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		English     string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		Italian     string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		German      string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		Russian     string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		Mandarin    string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		Other       string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
		OtherName   string `gorm:"type:ENUM('Não Aplicável','Fraco','Suficiente','Bom','Excelente');default:'Não Aplicável';not null"`
	}

	FamilyInfo struct {
		Base
		Application  string
		Relation     string `gorm:"type:ENUM('Pai','Mãe','Cônjuge');default:'Pai';not null"`
		Name         string
		BirthCountry string
		Profession   string
		Employer     string
		Salary       decimal.Decimal `sql:"type:decimal(12,2)"`
	}
)

func (a Application) ValidateDocuments(required []string) (err error) {
	if len(required) > len(a.Documents) {
		zaplog.ZLog(fmt.Errorf("A candidatura tem documentos em falta"))
	}

	documentMap := make(map[string]string, len(a.Documents))
	for _, d := range a.Documents {
		documentMap[d.Type] = d.URL
	}

	for _, r := range required {
		if documentMap[r] == "" {
			return zaplog.ZLog(fmt.Errorf("Documento em falta: %s", r))
		}
	}

	if a.BirthDate.Before(time.Now().AddDate(-18, 0, 0)) && a.Gender == GenderMale && documentMap[MilitaryDocument] == "" {
		return zaplog.ZLog(errors.New("O documento militar é obrigatório para homens maiores de 18 anos"))
	}

	//Enrolled
	if a.CourseYear >= 1 && documentMap[EnrollmentProof] == "" {
		return zaplog.ZLog(errors.New("A declaração de frequência é obrigatória para quem já frequenta o curso"))
	}
	//Homologado - postgrad application
	if a.EducationLevel != HighSchool && documentMap[CertificateValidationScan] == "" {
		return zaplog.ZLog(errors.New("O certificado homologado é obrigatório"))
	}
	return
}

func (FamilyInfo) TableName() string {
	return "family_info"
}

func (a Application) AllowedStatuses(newStatus string) bool {
	allowed, ok := applicationStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}
