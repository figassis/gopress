package model

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/mholt/archiver/v3"
	"github.com/shopspring/decimal"
)

var (
	scholarshipStatusFlows = map[string]List{
		StatusDraft:                List{StatusTrash, StatusOpen, StatusScholarshipPublished},
		StatusScholarshipPublished: List{StatusDraft, StatusOpen},
		StatusOpen:                 List{StatusDraft, StatusScholarshipPublished, StatusCompleted, StatusClosed},
		StatusClosed:               List{StatusDraft, StatusScholarshipPublished, StatusCompleted, StatusOpen},
		StatusCompleted:            List{},
		StatusTrash:                List{StatusDraft, StatusScholarshipPublished},
	}

	statsCache = time.Minute * 10
)

type (
	Scholarship struct {
		Base
		Name              string
		Sponsor           string
		SponsorName       string
		Start             time.Time
		End               time.Time
		Available         int64
		MaxAge            int64
		MinGrade          int64
		TotalApplications int64
		Content           string `sql:"type:longtext"`
		html              string
		EnableQuotas      bool
		RPEQuota          decimal.Decimal             `sql:"type:decimal(5,2)"`
		PriorityQuota     decimal.Decimal             `sql:"type:decimal(5,2)"`
		ProvinceQuota     ProvinceQuota               `gorm:"foreignkey:Scholarship;association_foreignkey:uuid;association_autoupdate:false"`
		Type              string                      `gorm:"type:ENUM('Interna','Externa','Externa de Mérito');default:'Interna';not null"`
		Level             string                      `gorm:"type:ENUM('Licenciatura','Mestrado','Especialidade','Doutoramento');default:'Licenciatura';not null"`
		Status            string                      `gorm:"type:ENUM('Rascunho','Publicada','Aberta','Encerrada','Concluída','Lixeira');default:'Rascunho';not null"`
		Documents         []File                      `gorm:"foreignkey:ResourceID;association_foreignkey:uuid;association_autoupdate:false"`
		Applications      []Application               `gorm:"foreignkey:Scholarship;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
		RequiredDocuments List                        `gorm:"type:varchar(256)" sql:"type:varchar(256)"`
		Stats             ApplicationStats            `gorm:"-"`
		ProvinceStats     map[string]ApplicationStats `gorm:"-"`
	}

	ApplicationStats struct {
		Total        int64
		TotalRPE     int64
		Pending      int64
		Validated    int64
		Approved     int64
		Awarded      int64
		Rejected     int64
		Review       int64
		PendingRPE   int64
		ValidatedRPE int64
		ApprovedRPE  int64
		AwardedRPE   int64
		RejectedRPE  int64
		ReviewRPE    int64
	}
)

func (p Scholarship) GetHtml() (html string, err error) {
	data, err := ioutil.ReadFile(os.Getenv("ASSETS") + "/html/" + p.ID)
	if err != nil {
		if err = p.generateHtml(); err != nil {
			return
		}
		return p.html, nil
	}

	return string(data), nil
}

func (p Scholarship) generateHtml() (err error) {
	p.html, err = wysiwygToHTML(p.Content)
	if err != nil {
		return
	}

	return ioutil.WriteFile(os.Getenv("ASSETS")+"/html/"+p.ID, []byte(p.html), 0644)
}

func (a Scholarship) AllowedStatuses(newStatus string) bool {
	allowed, ok := scholarshipStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}

func (s *Scholarship) GetStats(db *gorm.DB, realtime bool) {
	var stats ApplicationStats

	if s.ProvinceStats == nil {
		s.ProvinceStats = make(map[string]ApplicationStats, len(Angola.Cities))
	}

	if !realtime {
		if err := zaplog.ZLog(util.GetCache(fmt.Sprintf("/scholarships/%s/stats/global", s.ID), &s.Stats)); err == nil {
			if err := zaplog.ZLog(util.GetCache(fmt.Sprintf("/scholarships/%s/stats/provinces", s.ID), &s.ProvinceStats)); err == nil {
				return
			}
		}
	}

	q := db.Model(&Application{}).Where("scholarship = ?", s.ID)

	zaplog.ZLog(q.Count(&stats.Total).Error)
	zaplog.ZLog(q.Where("rpe = ?", true).Count(&stats.TotalRPE).Error)
	zaplog.ZLog(q.Where("status = ?", StatusPending).Count(&stats.Pending).Error)
	zaplog.ZLog(q.Where("status = ?", StatusValidated).Count(&stats.Validated).Error)
	zaplog.ZLog(q.Where("status = ?", StatusApprovedApplication).Count(&stats.Approved).Error)
	zaplog.ZLog(q.Where("status = ?", StatusAwarded).Count(&stats.Awarded).Error)
	zaplog.ZLog(q.Where("status = ?", StatusRejected).Count(&stats.Rejected).Error)
	zaplog.ZLog(q.Where("status = ?", StatusNeedsReview).Count(&stats.Review).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusPending, true).Count(&stats.PendingRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusValidated, true).Count(&stats.ValidatedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusApprovedApplication, true).Count(&stats.ApprovedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusAwarded, true).Count(&stats.AwardedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusRejected, true).Count(&stats.RejectedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusNeedsReview, true).Count(&stats.ReviewRPE).Error)
	s.Stats = stats

	for province := range Angola.Cities {
		s.ProvinceStats[province] = getProvinceStats(province, q)
	}

	go util.CacheTTL(fmt.Sprintf("/scholarships/%s/stats/global", s.ID), s.Stats, statsCache)
	go util.CacheTTL(fmt.Sprintf("/scholarships/%s/stats/provinces", s.ID), s.ProvinceStats, statsCache)

}

func getProvinceStats(province string, db *gorm.DB) (stats ApplicationStats) {
	q := db.Where("current_province = ?", province)

	zaplog.ZLog(q.Count(&stats.Total).Error)
	zaplog.ZLog(q.Where("rpe = ?", true).Count(&stats.TotalRPE).Error)
	zaplog.ZLog(q.Where("status = ?", StatusPending).Count(&stats.Pending).Error)
	zaplog.ZLog(q.Where("status = ?", StatusValidated).Count(&stats.Validated).Error)
	zaplog.ZLog(q.Where("status = ?", StatusApprovedApplication).Count(&stats.Approved).Error)
	zaplog.ZLog(q.Where("status = ?", StatusAwarded).Count(&stats.Awarded).Error)
	zaplog.ZLog(q.Where("status = ?", StatusRejected).Count(&stats.Rejected).Error)
	zaplog.ZLog(q.Where("status = ?", StatusNeedsReview).Count(&stats.Review).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusPending, true).Count(&stats.PendingRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusValidated, true).Count(&stats.ValidatedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusApprovedApplication, true).Count(&stats.ApprovedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusAwarded, true).Count(&stats.AwardedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusRejected, true).Count(&stats.RejectedRPE).Error)
	zaplog.ZLog(q.Where("status = ? AND rpe = ?", StatusNeedsReview, true).Count(&stats.ReviewRPE).Error)
	return
}

func (s Scholarship) GetProvinceQuota(province string) (quota decimal.Decimal, err error) {
	switch province {
	case "Bengo":
		return s.ProvinceQuota.Bengo, nil
	case "Benguela":
		return s.ProvinceQuota.Benguela, nil
	case "Bié":
		return s.ProvinceQuota.Bie, nil
	case "Cabinda":
		return s.ProvinceQuota.Cabinda, nil
	case "Cuando Cubango":
		return s.ProvinceQuota.CuandoCubango, nil
	case "Cuanza Norte":
		return s.ProvinceQuota.CuanzaNorte, nil
	case "Cuanza Sul":
		return s.ProvinceQuota.CuanzaSul, nil
	case "Cunene":
		return s.ProvinceQuota.Cunene, nil
	case "Huambo":
		return s.ProvinceQuota.Huambo, nil
	case "Huíla":
		return s.ProvinceQuota.Huila, nil
	case "Luanda":
		return s.ProvinceQuota.Luanda, nil
	case "Lunda Norte":
		return s.ProvinceQuota.LundaNorte, nil
	case "Lunda Sul":
		return s.ProvinceQuota.LundaSul, nil
	case "Malanje":
		return s.ProvinceQuota.Malanje, nil
	case "Moxico":
		return s.ProvinceQuota.Moxico, nil
	case "Namibe":
		return s.ProvinceQuota.Namibe, nil
	case "Uíge":
		return s.ProvinceQuota.Uige, nil
	case "Zaire":
		return s.ProvinceQuota.Zaire, nil
	}
	err = fmt.Errorf("A província %s é inválida", province)
	return
}

func (s Scholarship) Export(db *gorm.DB, force bool) (url string, err error) {
	var f File
	now := time.Now()
	if err = zaplog.ZLog(db.Model(&File{}).Where("resource = ? AND resource_id = ? AND extension = ?", ResourceScholarship, s.ID, "zip").First(&f).Error); err == nil {
		//Check if we already have the last export or if report was generated in the last 10 minutes, do not generate a new one
		if !force && (strings.Contains(strings.ToLower(f.Name), "final") || f.UpdatedAt.After(now.Add(-10*time.Minute))) {
			return f.GetURL()
		}
	}

	export, err := s.export(db)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if f.ID != "" {
		f.Name, f.URL, f.Path, f.Location, f.Public = export.Name, export.URL, export.Path, export.Location, export.Public
		if err = zaplog.ZLog(db.Save(&f).Error); err != nil {
			return
		}
	} else {
		id, err2 := util.GenerateUUID()
		if err2 != nil {
			return
		}
		f = File{
			Base:       Base{ID: id},
			UserID:     AdminUser,
			UserName:   "System",
			Name:       export.Name,
			Resource:   ResourceScholarship,
			ResourceID: s.ID,
			Type:       GeneralFile,
			Extension:  "zip",
			URL:        export.URL,
			Status:     "Aprovado",
			Public:     export.Public,
			Location:   export.Location,
		}
		if err = zaplog.ZLog(db.Create(&f).Error); err != nil {
			return
		}
	}

	return f.GetURL()

}

func (s Scholarship) export(db *gorm.DB) (file File, err error) {
	var statuses = []string{StatusPending, StatusValidated, StatusApprovedApplication, StatusNeedsReview, StatusRejectedApplication, StatusAwarded, StatusCanceledApplication}
	var paths []string

	for _, status := range statuses {
		path, err2 := s.exportStatus(db, status)
		if err2 != nil {
			continue
		}
		paths = append(paths, path)
	}

	name := fmt.Sprintf("%s/%s-export.zip", os.Getenv("UPLOADS"), s.ID)
	if s.Status == StatusCompleted {
		name = fmt.Sprintf("%s/%s-export-final.zip", os.Getenv("UPLOADS"), s.ID)
	}

	if err = zaplog.ZLog(compressFiles(paths, name)); err != nil {
		return
	}

	id, err := util.GenerateUUID()
	if zaplog.ZLog(err); err != nil {
		return
	}

	public := true
	url, err := SaveFile(name, GetS3KeyFromID(id, public))
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	file = File{
		Base:       Base{ID: id},
		UserID:     AdminUser,
		UserName:   "System",
		Name:       name,
		Resource:   ResourceScholarship,
		ResourceID: s.ID,
		Type:       GeneralFile,
		Extension:  "zip",
		URL:        url,
		Status:     "Aprovado",
		Public:     public,
		Location:   "s3",
	}

	return
}

func (s Scholarship) exportStatus(db *gorm.DB, status string) (path string, err error) {

	q := db.Model(&Application{}).Where("scholarship = ?", s.ID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var applications []Application
	if err = zaplog.ZLog(q.Find(&applications).Error); err != nil {
		return
	}

	if len(applications) == 0 {
		err = errors.New("No applications with status " + status)
		return
	}

	//Make CSV / XLSX
	path = fmt.Sprintf("%s/%s-%s.csv", os.Getenv("UPLOADS"), s.ID, strings.ToLower(status))
	if s.Status == StatusCompleted {
		path = fmt.Sprintf("%s/%s-%s-final.csv", os.Getenv("UPLOADS"), s.ID, strings.ToLower(status))
	}
	if err = zaplog.ZLog(s.saveApplicationsToFile(&applications, path)); err != nil {
		return
	}
	return
}

func (s Scholarship) saveApplicationsToFile(applications *[]Application, path string) (err error) {
	file, err := os.Create(path)
	if err = zaplog.ZLog(err); err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err = zaplog.ZLog(writer.Write(applicationCSVHeaders)); err != nil {
		err = errors.New("Cannot write to file")
		return
	}

	for _, a := range *applications {
		rpe := "Normal"
		if a.RPE {
			rpe = "RPE"
		}

		var father, mother, spouse FamilyInfo
		for _, info := range a.Family {
			if info.Relation == "Pai" {
				father = info
			}
			if info.Relation == "Mãe" {
				mother = info
			}
			if info.Relation == "Cônjuge" {
				spouse = info
			}
		}

		var docs = make(map[string]string, len(a.Documents))
		for _, doc := range a.Documents {
			docs[doc.Type], _ = doc.GetURL()
		}

		value := []string{
			//Identity
			a.ID, a.Status, a.Name, a.IDNumber, a.PassportNumber, rpe, a.Score.StringFixed(2), a.Gender, a.Phone, a.Email, a.Skype,

			//Scholarship
			a.Scholarship, a.ScholarshipName, s.Level, s.Type,

			//Birth and Residence
			a.BirthDate.Format("02/01/2006"), a.BirthProvince, "Angola", a.CurrentProvince, a.CurrentCountry,

			//Bank
			a.BankName, a.BankAccountNumber, a.BankAccountOwner, a.IBAN,

			//Employment
			a.Employer, a.EmployerAddress, a.JobTitle, a.Salary.StringFixed(2),

			//Family
			strconv.Itoa(int(a.FamilyMembers)), strconv.Itoa(int(a.WorkingMembers)), strconv.Itoa(int(a.Children)),
			father.Name, father.BirthCountry, father.Profession, father.Employer, father.Salary.StringFixed(2),
			mother.Name, mother.BirthCountry, mother.Profession, mother.Employer, mother.Salary.StringFixed(2),
			spouse.Name, spouse.BirthCountry, spouse.Profession, spouse.Employer, spouse.Salary.StringFixed(2),

			//Education
			a.EducationLevel, a.EducationSchool, a.EducationSchoolDepartment, a.EducationCourse, a.EducationCountry, a.EducationProvince,
			a.GraduationDate.Format("02/01/2006"), a.EducationGrade.StringFixed(1), a.EducationEvaluation,

			//Application
			a.CourseSchool, a.CourseDepartment, a.CourseOption1, a.CourseOption2, a.CourseOption3, a.CourseCountry, a.CourseProvince, strconv.Itoa(int(a.CourseYear)), a.CourseAverage.StringFixed(1),

			//Languages
			a.Languages.French, a.Languages.Spanish, a.Languages.English, a.Languages.Italian, a.Languages.German, a.Languages.Russian, a.Languages.Mandarin, fmt.Sprintf("%s:%s", a.Languages.OtherName, a.Languages.Other),
			docs[IDScan], docs[PassportScan], docs[DiplomaScan], docs[CertificateScan], docs[CertificateValidationScan], docs[EnrollmentProof], docs[ResidenceProof],
			docs[HealthCertificate], docs[ProjectProposal], docs[PovertyDeclaration], docs[RPEProof], docs[RecommendationLetter], docs[ReturnAgreement],
			docs[WorkDeclaration], docs[IESAppraisal], docs[MilitaryDocument],
		}

		if len(value) != len(applicationCSVHeaders) {
			err = zaplog.ZLog(errors.New("Incorrect CSV header length for application " + a.ID))
			return
		}

		if err = zaplog.ZLog(writer.Write(value)); err != nil {
			err = errors.New("Cannot write to file")
			return
		}
	}
	return
}

func compressFiles(paths []string, path string) (err error) {
	z := archiver.Zip{
		// MkdirAll:               true,
		// SelectiveCompression:   true,
		// ContinueOnError:        false,
		OverwriteExisting: true,
		// ImplicitTopLevelFolder: false,
	}

	if err = zaplog.ZLog(z.Archive(paths, path)); err != nil {
		return
	}

	return
}

func RequiredDocuments(sType, level string) (result List, err error) {
	result, ok := requiredDocuments[sType+level]
	if !ok {
		err = fmt.Errorf("Invalid scholarship type (%s) or level (%s)", sType, level)
		return
	}
	return
}
