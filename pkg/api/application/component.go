// Package user contains user application services
package application

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/query"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

// Create creates a new user account
func (u *App) Create(c echo.Context, req *Create) (*model.Application, error) {
	au := u.rbac.User(c)
	var user model.User
	var scholarship model.Scholarship

	if au.Role != model.CandidateRole {
		return nil, errors.New("Apenas candidatos podem submeter candidaturas")
	}

	if err := u.db.Model(&model.User{}).Where("uuid = ?", au.ID).First(&user).Error; err != nil {
		return nil, err
	}

	if user.Status != model.StatusActive {
		return nil, zaplog.ZLog(fmt.Errorf("User %s is %s, so cannot apply", user.Email, user.Status))
	}

	if err := u.db.Model(&model.Scholarship{}).Where("uuid = ?", req.Scholarship).First(&scholarship).Error; err != nil {
		return nil, err
	}

	if scholarship.Status != model.StatusOpen {
		err := zaplog.ZLog(errors.New("A bolsa não está aberta"))
		return nil, err
	}
	now := time.Now()
	if scholarship.Start.After(now) {
		err := zaplog.ZLog(errors.New("A bolsa ainda não está aberta"))
		return nil, err
	}
	if scholarship.End.Before(now) {
		err := zaplog.ZLog(errors.New("A bolsa já está encerrada"))
		return nil, err
	}

	id, err := util.GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return nil, err
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", req.Documents).Update(model.File{
		Resource:   model.ResourceApplication,
		ResourceID: id,
	}).Error)

	var docs []model.File

	if err := u.db.Model(&model.Scholarship{}).Where("uuid IN (?)", req.Documents).Find(&docs).Error; err != nil {
		return nil, err
	}

	if len(scholarship.RequiredDocuments) > len(docs) {
		return nil, zaplog.ZLog(fmt.Errorf("Documentos em falta: %s", strings.Join(scholarship.RequiredDocuments, ", ")))
	}

	for _, rdoc := range scholarship.RequiredDocuments {
		found := false
		for _, doc := range docs {
			if doc.Type == rdoc {
				found = true
				break
			}
		}
		if !found {
			return nil, zaplog.ZLog(fmt.Errorf("Documento em falta: %s", rdoc))
		}
	}

	var educationLevel = model.HighSchool
	switch scholarship.Level {
	case model.Masters, model.Specialization:
		educationLevel = model.GradOther
	case model.Doctorate:
		educationLevel = model.PostGrad
	}

	application := model.Application{
		Base:               model.Base{ID: id},
		Scholarship:        scholarship.ID,
		ScholarshipName:    fmt.Sprintf("Bolsa %s de %s", scholarship.Type, scholarship.Level),
		UserID:             u.rbac.User(c).ID,
		IDNumber:           req.IDNumber,
		PassportNumber:     req.PassportNumber,
		Name:               user.Name,
		BirthDate:          req.BirthDate,
		Gender:             req.Gender,
		BirthProvince:      req.BirthProvince,
		BirthCity:          req.BirthCity,
		CurrentCountry:     req.CurrentCountry,
		CurrentProvince:    req.CurrentProvince,
		CurrentCity:        req.CurrentCity,
		Phone:              req.Phone,
		Email:              req.Email,
		Skype:              req.Skype,
		BankName:           req.BankName,
		BankAccountNumber:  req.BankAccountNumber,
		BankAccountOwner:   req.BankAccountOwner,
		IBAN:               req.IBAN,
		RPE:                req.RPE,
		ScoreJustification: req.ScoreJustification,
		FamilyMembers:      req.FamilyMembers,
		Children:           req.Children,
		WorkingMembers:     req.WorkingMembers,

		EducationLevel:            educationLevel,
		EducationSchool:           req.EducationSchool,
		EducationSchoolDepartment: req.EducationSchoolDepartment,
		EducationSchoolContact:    req.EducationSchoolContact,
		EducationCourse:           req.EducationCourse,
		EducationCountry:          req.EducationCountry,
		EducationProvince:         req.EducationProvince,
		GraduationDate:            req.GraduationDate,
		EducationGrade:            req.EducationGrade,
		EducationEvaluation:       req.EducationEvaluation,

		ApplicationLevel: req.ApplicationLevel,
		CourseSchool:     req.CourseSchool,
		CourseDepartment: req.CourseDepartment,
		CourseOption1:    req.CourseOption1,
		CourseOption2:    req.CourseOption2,
		CourseOption3:    req.CourseOption3,
		CourseCountry:    req.CourseCountry,
		CourseProvince:   req.CourseProvince,
		CourseContact:    req.CourseContact,
		CourseYear:       req.CourseYear,
		CourseAverage:    req.CourseAverage,

		EducationCourseID:   req.EducationCourseID,
		ApplicationCourseID: req.ApplicationCourseID,

		Employer:                 req.Employer,
		EmployerAddress:          req.EmployerAddress,
		EmployerContact:          req.EmployerContact,
		JobTitle:                 req.JobTitle,
		Salary:                   req.Salary,
		HasEmployerAuthorization: req.HasEmployerAuthorization,
		Status:                   req.Status,
		Languages:                req.Languages,
		Family:                   req.Family,
	}
	return u.udb.Create(u.db, application)
}

// List returns list of users
func (u *App) List(c echo.Context, p *model.Pagination) ([]model.Application, string, string, int64, int64, error) {
	au := u.rbac.User(c)

	// if err := u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
	// 	return []model.Application{}, "", "", 0, 0, err
	// }

	q, err := query.List(au, model.ResourceApplication)
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	if c.QueryString() != "" {
		p.ApplicationQuery = &model.Application{}
		params := c.QueryParams()

		p.ApplicationQuery.Scholarship = params.Get("scholarship")
		p.ApplicationQuery.Operator = params.Get("operator")
		p.ApplicationQuery.Status = params.Get("status")
		p.ApplicationQuery.Email = params.Get("email")
		p.ApplicationQuery.IDNumber = params.Get("id_number")
		p.ApplicationQuery.UserID = params.Get("user")
		p.ApplicationQuery.RPE = params.Get("rpe") == "true"
		p.ApplicationQuery.CurrentProvince = params.Get("current_province")
		p.SearchQuery = params.Get("s")
	}

	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *App) View(c echo.Context, id string) (*model.Application, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u *App) Delete(c echo.Context, id string) error {
	if err := u.rbac.EnforceRole(c, model.AdminRole); err != nil {
		return err
	}

	application, err := u.udb.View(u.db, id)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	if application.Status != model.StatusDuplicate {
		return zaplog.ZLog(fmt.Errorf("Apenas se pode eliminar candidaturas duplicadas."))
	}

	return u.udb.Delete(u.db, id)
}

// Update updates user's contact information
func (u *App) Update(c echo.Context, r *Update) (result *model.Application, err error) {
	if err = u.rbac.EnforceRole(c, model.OperatorRole); err != nil {
		return
	}

	update, err := u.udb.View(u.db, r.ID)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if r.Status != "" && update.Status != r.Status && !update.AllowedStatuses(r.Status) {
		err = zaplog.ZLog(fmt.Errorf("Não é possível passar de %s para %s", update.Status, r.Status))
		return
	}

	var scholarship model.Scholarship
	var group []model.Scholarship
	if err = u.db.Model(&model.Scholarship{}).Where("uuid = ?", update.Scholarship).First(&scholarship).Error; err != nil {
		zaplog.ZLog(err)
		return
	}

	query := model.Scholarship{
		Sponsor:       scholarship.Sponsor,
		Name:          scholarship.Name,
		Start:         scholarship.Start,
		End:           scholarship.End,
		Available:     scholarship.Available,
		EnableQuotas:  scholarship.EnableQuotas,
		RPEQuota:      scholarship.RPEQuota,
		PriorityQuota: scholarship.PriorityQuota,
		Status:        scholarship.Status,
	}
	var ids = []string{scholarship.ID}
	if err = u.db.Model(&model.Scholarship{}).Where(query).Where("uuid != ?", scholarship.ID).Find(&group).Error; err == nil {
		for _, s := range group {
			ids = append(ids, s.ID)
		}
	}

	//If approving or awarding, check quotas
	if scholarship.Available > 0 && scholarship.Type == model.InternalScholarship && r.Status != "" && update.Status != r.Status && r.Status == model.StatusAwarded {
		province := update.CurrentProvince
		var provincePercentQuota decimal.Decimal
		if provincePercentQuota, err = scholarship.GetProvinceQuota(province); err != nil {
			zaplog.ZLog(err)
			return
		}

		var rpeCount, regularCount, provinceCount, rpeQuota, regularQuota, provinceQuota int64

		//When province quota and rpe quota are enabled, rpe quota is calculated within each province
		if provincePercentQuota.GreaterThan(decimal.NewFromFloat(0)) {
			provinceQuota = provincePercentQuota.Div(decimal.NewFromFloat(100)).Mul(decimal.NewFromFloat(float64(scholarship.Available))).IntPart()
			if err = u.db.Model(&model.Application{}).Where("scholarship IN ? AND status = ? AND current_province = ?", ids, r.Status, update.CurrentProvince).Count(&provinceCount).Error; err != nil {
				return
			}
			if provinceCount >= provinceQuota {
				err = zaplog.ZLog(fmt.Errorf("Já atingiu a quota de %s%% (%d/%d) para %s", provincePercentQuota.StringFixed(0), provinceCount, scholarship.Available, province))
				return
			}

			if scholarship.RPEQuota.GreaterThan(decimal.NewFromFloat(0)) {
				rpeQuota = scholarship.RPEQuota.Div(decimal.NewFromFloat(100)).Mul(decimal.NewFromFloat(float64(provinceQuota))).IntPart()
				if err = u.db.Model(&model.Application{}).Where("scholarship IN ? AND status = ? AND current_province = ? AND rpe = ?", ids, r.Status, update.CurrentProvince, true).Count(&rpeCount).Error; err != nil {
					return
				}

				regularCount = provinceCount - rpeCount
				regularQuota = provinceQuota - rpeQuota

				if update.RPE {
					if rpeCount >= rpeQuota {
						err = zaplog.ZLog(fmt.Errorf("Já atingiu a quota RPE de %s%% (%d/%d) para %s", scholarship.RPEQuota.StringFixed(2), rpeCount, provinceQuota, province))
						return
					}
				} else {
					if regularCount >= regularQuota {
						err = zaplog.ZLog(fmt.Errorf("Já atingiu a quota do regime normal de %s%% (%d/%d) para %s", decimal.NewFromFloat(100).Sub(scholarship.RPEQuota).StringFixed(2), regularCount, provinceQuota, province))
						return
					}
				}
			}

		} else {
			//When only rpe quota is enabled, it is calculated from global stats
			if scholarship.RPEQuota.GreaterThan(decimal.NewFromFloat(0)) {
				rpeQuota = scholarship.RPEQuota.Div(decimal.NewFromFloat(100)).Mul(decimal.NewFromFloat(float64(scholarship.Available))).IntPart()
				if err = u.db.Model(&model.Application{}).Where("scholarship IN ? AND status = ? AND rpe = ?", ids, r.Status, true).Count(&rpeCount).Error; err != nil {
					return
				}
				if err = u.db.Model(&model.Application{}).Where("scholarship IN ? AND status = ? AND rpe = ?", ids, r.Status, false).Count(&regularCount).Error; err != nil {
					return
				}

				regularQuota = scholarship.Available - rpeQuota

				if update.RPE {
					if rpeCount >= rpeQuota {
						err = zaplog.ZLog(fmt.Errorf("Já atingiu a quota RPE global de %s%% (%d/%d)", scholarship.RPEQuota.StringFixed(2), rpeCount, scholarship.Available))
						return
					}
				} else {
					if regularCount >= regularQuota {
						err = zaplog.ZLog(fmt.Errorf("Já atingiu a quota global do regime normal de %s%% (%d/%d)", decimal.NewFromFloat(100).Sub(scholarship.RPEQuota).StringFixed(2), regularCount, scholarship.Available))
						return
					}
				}
			}
		}
	}

	zaplog.ZLog(u.db.Model(&model.File{}).Where("uuid IN (?)", r.Documents).Update(model.File{
		Resource:   model.ResourceApplication,
		ResourceID: update.ID,
	}).Error)

	// update.BirthDate = r.BirthDate
	// update.Gender = r.Gender
	// update.BirthProvince = r.BirthProvince
	// update.BirthCity = r.BirthCity
	// update.CurrentCountry = r.CurrentCountry
	// update.CurrentProvince = r.CurrentProvince
	// update.CurrentCity = r.CurrentCity
	// update.Phone = r.Phone
	// update.Email = r.Email
	// update.Skype = r.Skype
	// update.BankName = r.BankName
	// update.BankAccountNumber = r.BankAccountNumber
	// update.BankAccountOwner = r.BankAccountOwner
	// update.IBAN = r.IBAN
	update.RPE = r.RPE
	update.Operator = r.Operator
	// update.FamilyMembers = r.FamilyMembers
	// update.Children = r.Children
	// update.WorkingMembers = r.WorkingMembers
	// update.EducationLevel = r.EducationLevel
	// update.EducationSchool = r.EducationSchool
	// update.EducationSchoolDepartment = r.EducationSchoolDepartment
	// update.EducationSchoolContact = r.EducationSchoolContact
	// update.EducationCourse = r.EducationCourse
	// update.EducationCountry = r.EducationCountry
	// update.EducationProvince = r.EducationProvince
	// update.GraduationDate = r.GraduationDate
	// update.EducationGrade = r.EducationGrade
	// update.EducationEvaluation = r.EducationEvaluation
	// update.ApplicationLevel = r.ApplicationLevel
	// update.CourseSchool = r.CourseSchool
	// update.CourseDepartment = r.CourseDepartment
	// update.CourseOption1 = r.CourseOption1
	// update.CourseOption2 = r.CourseOption2
	// update.CourseOption3 = r.CourseOption3
	// update.CourseCountry = r.CourseCountry
	// update.CourseProvince = r.CourseProvince
	// update.CourseContact = r.CourseContact
	// update.CourseYear = r.CourseYear
	// update.CourseAverage = r.CourseAverage
	// update.Employer = r.Employer
	// update.EmployerAddress = r.EmployerAddress
	// update.EmployerContact = r.EmployerContact
	// update.JobTitle = r.JobTitle
	// update.Salary = r.Salary
	// update.HasEmployerAuthorization = r.HasEmployerAuthorization
	update.Status = r.Status
	update.Languages = r.Languages
	// update.Family = r.Family

	var operator model.User
	if err = u.db.Model(&model.User{}).Where("uuid = ?", r.Operator).First(&operator).Error; err == nil {
		update.OperatorName = operator.Name
	}

	if err = u.udb.Update(u.db, update); err != nil {
		zaplog.ZLog(err)
		return
	}
	return u.udb.View(u.db, r.ID)
}
