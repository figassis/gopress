package model

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

var (
	testPassword = "password"
)

func GenerateFakeData(number int, db *gorm.DB) (err error) {
	organizationIDs := util.ShortIDS(number)

	for _, id := range organizationIDs {
		orgnaization := Organization{
			Base:     Base{ID: id},
			Name:     gofakeit.Company(),
			Status:   StatusActive,
			Type:     OrgSponsor,
			Country:  gofakeit.Country(),
			Province: gofakeit.State(),
			City:     gofakeit.City(),
			Logo:     "assets/images/thumb-1.jpg",
		}

		if err = zaplog.ZLog(db.Create(&orgnaization).Error); err != nil {
			return
		}
	}

	userIDs := util.ShortIDS(number)

	for _, id := range userIDs {
		var org string
		for _, o := range organizationIDs {
			org = o
			break
		}

		user := User{
			Base:             Base{ID: id},
			Name:             gofakeit.Name(),
			Username:         gofakeit.Email(),
			Password:         util.Hash(testPassword),
			Email:            gofakeit.Email(),
			Phone:            gofakeit.Phone(),
			Status:           StatusActive,
			Role:             AdminRole,
			Organization:     org,
			OrganizationName: gofakeit.Name(),
			UnsubscribeID:    util.ShortID(),
			Unsubscribed:     true,
		}

		if err = zaplog.ZLog(db.Create(&user).Error); err != nil {
			return
		}
	}

	// courseIDs := util.ShortIDS(number)

	// for _, id := range courseIDs {
	// 	var org string
	// 	for _, o := range organizationIDs {
	// 		org = o
	// 		break
	// 	}

	// 	course := Course{
	// 		Base:       Base{ID: id},
	// 		Name:       gofakeit.Sentence(3),
	// 		Domain:     gofakeit.Sentence(3),
	// 		Cluster:    gofakeit.Sentence(3),
	// 		Type:       NormalCourse,
	// 		School:     org,
	// 		Department: gofakeit.Company(),
	// 		SchoolName: gofakeit.Company(),
	// 		Level:      Masters,
	// 	}

	// 	if err = zaplog.ZLog(db.Create(&course).Error); err != nil {
	// 		return
	// 	}
	// }

	scholarshipIDs := util.ShortIDS(number)

	for _, id := range scholarshipIDs {
		var org string
		for _, o := range organizationIDs {
			org = o
			break
		}

		scholarship := Scholarship{
			Base:              Base{ID: id},
			Name:              gofakeit.Sentence(3),
			Sponsor:           org,
			SponsorName:       gofakeit.Name(),
			Start:             time.Now(),
			End:               time.Now().Add(time.Hour * 24 * 30),
			Available:         int64(gofakeit.Float64Range(10, 10000)),
			MaxAge:            45,
			MinGrade:          10,
			RPEQuota:          decimal.NewFromFloat(gofakeit.Float64Range(10, 1000)),
			PriorityQuota:     decimal.NewFromFloat(gofakeit.Float64Range(10, 1000)),
			ProvinceQuota:     ProvinceQuota{},
			Type:              InternalScholarship,
			Level:             Masters,
			Status:            StatusOpen,
			RequiredDocuments: []string{IDScan, CertificateScan},
		}

		if err = zaplog.ZLog(db.Create(&scholarship).Error); err != nil {
			return
		}
	}

	applicationIDs := util.ShortIDS(number)

	for _, id := range applicationIDs {
		var s, user string
		for _, o := range scholarshipIDs {
			s = o
			break
		}

		for _, u := range userIDs {
			user = u
			break
		}

		application := Application{
			Base:               Base{ID: id},
			UserID:             user,
			Scholarship:        s,
			ScholarshipName:    gofakeit.Name(),
			Name:               gofakeit.Name(),
			BirthDate:          time.Now().AddDate(-20, -2, -10),
			Gender:             GenderFemale,
			BirthProvince:      "Luanda",
			BirthCity:          "Luanda",
			CurrentCountry:     "Angola",
			CurrentProvince:    "Luanda",
			CurrentCity:        "Luanda",
			Phone:              gofakeit.Phone(),
			Email:              gofakeit.Email(),
			Skype:              gofakeit.Email(),
			BankName:           gofakeit.Company(),
			BankAccountNumber:  gofakeit.Numerify("###-###-#####"),
			BankAccountOwner:   gofakeit.Name(),
			IBAN:               gofakeit.Numerify("###-###-#####"),
			RPE:                gofakeit.Bool(),
			Operator:           gofakeit.UUID(),
			OperatorName:       gofakeit.Name(),
			Score:              decimal.NewFromFloat(gofakeit.Float64Range(0, 100)),
			ScoreGrades:        decimal.NewFromFloat(gofakeit.Float64Range(0, 100)),
			ScoreAge:           decimal.NewFromFloat(gofakeit.Float64Range(0, 100)),
			ScoreJustification: gofakeit.Sentence(8),
			FamilyMembers:      int64(gofakeit.Number(0, 10)),
			Children:           int64(gofakeit.Number(0, 3)),
			WorkingMembers:     int64(gofakeit.Number(0, 4)),

			EducationLevel:            HighSchool,
			EducationSchool:           gofakeit.Company(),
			EducationSchoolDepartment: gofakeit.Company(),
			EducationSchoolContact:    gofakeit.Phone(),
			EducationCourse:           gofakeit.Sentence(3),
			EducationCountry:          "Angola",
			EducationProvince:         "Luanda",
			GraduationDate:            time.Now().AddDate(-5, 0, 0),
			EducationGrade:            decimal.NewFromFloat(gofakeit.Float64Range(0, 20)),
			EducationEvaluation:       GradeNA,

			ApplicationLevel: Grad,
			CourseSchool:     gofakeit.Company(),
			CourseDepartment: gofakeit.Company(),
			CourseOption1:    gofakeit.Sentence(3),
			CourseOption2:    gofakeit.Sentence(3),
			CourseOption3:    gofakeit.Sentence(3),
			CourseCountry:    "Angola",
			CourseProvince:   "Luanda",
			CourseContact:    gofakeit.Phone(),
			CourseYear:       int64(gofakeit.Number(0, 4)),
			CourseAverage:    decimal.NewFromFloat(gofakeit.Float64Range(0, 20)),

			Employer:                 gofakeit.Company(),
			EmployerAddress:          gofakeit.Address().Address,
			EmployerContact:          gofakeit.Phone(),
			JobTitle:                 gofakeit.JobTitle(),
			Salary:                   decimal.NewFromFloat(gofakeit.Float64Range(10000, 1000000)),
			HasEmployerAuthorization: gofakeit.Bool(),
			Status:                   StatusPending,
			// Languages:                Languages{Application: id, English: LanguageGood, French: LanguageWeak},
		}

		if err = zaplog.ZLog(db.Create(&application).Error); err != nil {
			return
		}

		lang := Languages{Application: id, English: LanguageGood, French: LanguageWeak}
		if err = zaplog.ZLog(db.Create(&lang).Error); err != nil {
			return
		}

		family := FamilyInfo{
			Base:         Base{ID: util.ShortID()},
			Application:  application.ID,
			Relation:     RelationSpouse,
			Name:         gofakeit.Name(),
			BirthCountry: "Angola",
			Profession:   gofakeit.JobTitle(),
			Employer:     gofakeit.Company(),
			Salary:       decimal.NewFromFloat(gofakeit.Float64Range(10000, 1000000)),
		}

		if err = zaplog.ZLog(db.Create(&family).Error); err != nil {
			return
		}
	}

	postIDs := util.ShortIDS(number)

	for _, id := range postIDs {
		var author string
		for _, a := range userIDs {
			author = a
			break
		}
		post := Post{
			Base:       Base{ID: id},
			Author:     author,
			AuthorName: gofakeit.Name(),
			Category:   CategoryNews,
			Tags:       []string{"this", "is", "a", "tag"},
			Title:      gofakeit.Sentence(5),
			Slug:       strings.ToLower(strings.ReplaceAll(gofakeit.Sentence(5), " ", "-")),
			Content:    gofakeit.Paragraph(5, 10, 10, ". "),
			Excerpt:    gofakeit.Sentence(10),
			Status:     StatusPublished,
		}
		if err = zaplog.ZLog(db.Create(&post).Error); err != nil {
			return
		}
	}

	return
}
