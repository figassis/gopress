package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type (
	UserImport struct {
		ID          string
		Email       string `json:"User Email"`
		FirstName   string `json:"First Name"`
		LastName    string `json:"Last Name"`
		Username    string
		DisplayName string `json:"Display Name"`
		Role        string `json:"User Role"`
		Password    string `json:"User Pass"`
	}

	ApplicationImport struct {
		Name                      string `json:"Nome"`
		Scholarship               string `json:"Bolsa"`
		Email                     string `json:"Email"`
		Phone                     string `json:"Telefone"`
		Score                     string `json:"Pontuação"`
		Status                    string `json:"Estado"`
		ApplicationDate           string `json:"Data"`
		CourseYear                string `json:"Ano Lectivo"`
		ApplicationSchoolContact  string `json:"Contacto da IES"`
		CourseOption1             string `json:"Opção de Curso 1"`
		CourseOption2             string `json:"Opção de Curso 2"`
		CourseOption3             string `json:"Opção de Curso 3"`
		ApplicationSchool         string `json:"IES"`
		ApplicationDepartment     string `json:"UO"`
		ApplicationLevel          string `json:"Nível da Candidatura"`
		BankAccount               string `json:"Conta Bancária"`
		BankName                  string `json:"Banco"`
		IDNumber                  string `json:"BI"`
		BirthProvince             string `json:"Naturalidade"`
		BirthDate                 string `json:"Data de Nascimento"` //1992-02-12 00:00:00
		ResidenceCountry          string `json:"País de Residência"`
		ResidenceProvince         string `json:"Residência"`
		ResidenceProvinceAbroad   string `json:"Residência no Exterior"`
		Certificate               string `json:"Certificado"`
		CertificateValidation     string `json:"Certificado Homologado"`
		ApplicationSchoolContact1 string `json:"Contacto da IES__1"`
		ApplicationCountry        string `json:"País da IES"`
		EducationCourse           string `json:"Curso Concluído"`
		Employer                  string `json:"Empregador"`
		HealthDocument            string `json:"Atestado Médico"`
		IBAN                      string `json:"IBAN"`
		IDExpiration              string `json:"Exp. BI"`
		IDScan                    string `json:"Cópia do BI"`
		IESAppraisal              string `json:"Parecer da IES"`
		MilitaryDocument          string `json:"Documento Militar"`
		Passport                  string `json:"Passaporte"`
		PassportExpiration        string `json:"Exp. Passaporte"`
		PassportScan              string `json:"Cópia do Passaporte"`
		PovertyProof              string `json:"Atestado de pobreza"`
		Project                   string `json:"Projecto de Investigação Científica"`
		RecommendationLetter      string `json:"Carta de Recomendação"`
		ResidenceProof            string `json:"Atestado de Residência"`
		ReturnAgreement           string `json:"Compromisso de Retorno"`
		RPE                       string `json:"Regime de Protecção Especial"`
		Gender                    string `json:"Sexo"`
		Skype                     string `json:"Skype"`
		ID                        string `json:"ID"`
		MotherSalary              string `json:"Salário da Mãe"`
		MotherProfession          string `json:"Profissão da Mãe"`
		MotherName                string `json:"Nome Completo da Mãe"`
		MotherCountry             string `json:"Nacionalidade da Mãe"`
		MotherEmployer            string `json:"Empregador da Mãe"`
		FatherSalary              string `json:"Salário do Pai"`
		FatherProfession          string `json:"Profissão do Pai"`
		FatherName                string `json:"Nome Completo do Pai"`
		FatherCountry             string `json:"Nacionalidade do Pai"`
		FatherEmployer            string `json:"Empregador do Pai"`
		EducationEvaluation       string `json:"Avaliação"`
		Educationprovince         string `json:"Província do Grau anterior"`
		EducationCountry          string `json:"País do Grau anterior"`
		EducationAverage          string `json:"Média do Grau anterior"`
		EducationSchool           string `json:"Instituição do Grau anterior"`
		EducationSchoolDepartment string `json:"Departamento / UO  da IES Anterior"`
		EducationSchoolContact    string `json:"Contacto da IES Anterior"`
		EducationLevel            string `json:"Nível do Curso Concluído"`
		ApplicationProvince       string `json:"Província da IES"`
		HasEmployerPermission     string `json:"Tem autorização do empregador?"`
		Salary                    string `json:"Salário"`
		EmpployerAddress          string `json:"Endereço do Empregador"`
		EmployerContact           string `json:"Contacto do Empregador"`
		JobTitle                  string `json:"Cargo / Função"`
		BirthCity                 string `json:"Município de Nascimento"`
		AccountOwner              string `json:"Titular da Conta"`
		ResidenceCity             string `json:"Cidade/Município de Residência"`
		SpouseSalary              string `json:"Salário do Cônjuge"`
		SpouseCountry             string `json:"Nacionalidade do Cônjuge"`
		SpouseName                string `json:"Nome Completo do Cônjuge"`
		SpouseProfession          string `json:"Profissão do Cônjuge"`
		SpouseEmployer            string `json:"Empregador do Cônjuge"`
		German                    string `json:"Alemão"`
		Spanish                   string `json:"Espanhol"`
		French                    string `json:"Francês"`
		English                   string `json:"Inglês"`
		Italian                   string `json:"Italiano"`
		Mandarin                  string `json:"Mandarin"`
		Russian                   string `json:"Russo"`
		ScoreAge                  string `json:"Pontuação (Idade)"`
		ScoreGrade                string `json:"Pontuação (Média)"`
		ScoreJustification        string `json:"Pontuação (Obs)"`
		EmployerAuthorization     string `json:"Tem autorização do empregador?__1"`
		FamilyMembers             string `json:"Membros da Família"`
		WorkingMembers            string `json:"Membros Trabalhadores"`
		Children                  string `json:"Menores de Idade"`
		GraduationYear            string `json:"Ano de Conclusão"`
		CourseAverage             string `json:"Média"`
	}

	ScholarshipImport struct {
		ID            string
		Title         string
		Content       string
		Excerpt       string
		Date          string //"2020-01-10"
		Start         string `json:"Início"`       //20200113
		End           string `json:"Encerramento"` //20200331
		Category      string `json:"Categoria"`
		Sponsor       string `json:"Patrocinador"`
		Document      string `json:"Edital"`
		Available     string `json:"Vagas"`
		MaxAge        string `json:"Idade Máxima"`
		MinGrade      string `json:"Média Mínima"`
		RPEQuota      string `json:"Quota para Regime de Protecção Especial"`
		PriorityQuota string `json:"Quota para Cursos Prioritários"`
		ProvinceQuota string `json:"Quotas por Província"`
	}
)

func ImportPosts(db *gorm.DB) (err error) {
	data := fmt.Sprintf("%s/import/posts.json", assetsDir)
	jsonFile, err := os.Open(data)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var admin User
	if err = zaplog.ZLog(db.Model(&User{}).Where("email = ?", os.Getenv("ADMIN_EMAIL")).First(&admin).Error); err != nil {
		return
	}

	var posts []Post
	if err = zaplog.ZLog(json.Unmarshal(bytes, &posts)); err != nil {
		return
	}

	for _, p := range posts {
		post := Post{
			Author:     admin.ID,
			AuthorName: admin.Name,
			Category:   CategoryNews,
			Title:      p.Title,
			Slug:       p.Slug,
			Content:    p.Content,
			Image:      p.Image,
			Excerpt:    p.Excerpt,
			Status:     StatusPublished,
		}

		if err = zaplog.ZLog(db.Where(post).FirstOrCreate(&post).Error); err != nil {
			continue
		}
	}

	return
}

func ImportUsers(db *gorm.DB) (err error) {
	fmt.Println("Importing users")
	data := fmt.Sprintf("%s/import/users.json", assetsDir)
	jsonFile, err := os.Open(data)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var users []UserImport
	if err = zaplog.ZLog(json.Unmarshal(bytes, &users)); err != nil {
		return
	}

	roles := map[string]AccessRole{"administrator": 100, "admin": 105, "operator": 110, "support": 115, "candidate": 130}

	var org Organization
	if err = zaplog.ZLog(db.Model(&Organization{}).Where("type = ?", OrgMain).First(&org).Error); err != nil {
		return
	}

	for i, user := range users {
		if user.Email == "" {
			continue
		}
		// zaplog.ZLog(fmt.Sprintf("Importing %s", user.Email))

		var count int
		if err := zaplog.ZLog(db.Model(&User{}).Where("email = ?", user.Email).Count(&count).Error); err == nil && count > 0 {
			continue
		}

		role, ok := roles[user.Role]
		if !ok {
			continue
		}
		ids, err := util.GenerateUUIDS(2)
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		u := User{
			Base:             Base{ID: ids[0]},
			Username:         user.Email,
			Password:         user.Password,
			Email:            user.Email,
			Status:           StatusActive,
			Role:             role,
			Organization:     org.ID,
			OrganizationName: org.Name,
			UnsubscribeID:    ids[1],
		}

		name := strings.TrimSpace(user.FirstName + " " + user.LastName)
		if name != "" {
			u.Name = strings.Title(strings.ToLower(name))
		} else if user.DisplayName != "" {
			u.Name = strings.Title(strings.ToLower(user.DisplayName))
		} else {
			u.Name = user.Email
		}

		if err := zaplog.ZLog(db.Create(&u).Error); err != nil {
			continue
		}

		if i%100 == 0 {
			zaplog.ZLog(fmt.Sprintf("Imported %d / %d users", i, len(users)))
		}

	}

	return
}

func ImportScholarships(db *gorm.DB) (scholarshipIds map[string]Scholarship, err error) {
	fmt.Println("Importing scholarships")
	data := fmt.Sprintf("%s/import/scholarships.json", assetsDir)
	jsonFile, err := os.Open(data)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var items []ScholarshipImport
	if err = zaplog.ZLog(json.Unmarshal(bytes, &items)); err != nil {
		return
	}

	var org Organization
	if err = zaplog.ZLog(db.Model(&Organization{}).Where("type = ?", OrgMain).First(&org).Error); err != nil {
		return
	}

	scholarshipIds = make(map[string]Scholarship)

	for _, item := range items {
		ids, err := util.GenerateUUIDS(2)
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		sponsor := org
		if item.Sponsor != org.Name {
			country, province, city := "Angola", "Luanda", "Luanda"
			if strings.Contains(item.Title, "China") {
				country, province, city = "China", "Beijing", "Beijing"
			}
			if strings.Contains(item.Title, "Portugal") {
				country, province, city = "Portugal", "Lisboa", "Lisboa"
			}
			sponsor = Organization{}
			if err2 := db.Model(&Organization{}).Where("name = ?", item.Sponsor).First(&sponsor).Error; err2 != nil {
				sponsor = Organization{Base: Base{ID: ids[1]}, Name: item.Sponsor, Status: StatusActive, Type: OrgSponsor, Country: country, Province: province, City: city, Phone: gofakeit.Phone(), Email: gofakeit.Email()}

				sponsor.Logo = fmt.Sprintf("assets/images/%s.png", strings.ToLower(sponsor.Country))
				if err = zaplog.ZLog(db.Create(&sponsor).Error); err != nil {
					continue
				}
			}
		}

		maxAge, _ := strconv.Atoi(item.MaxAge)
		minGrade, _ := strconv.Atoi(item.MinGrade)
		var status = StatusClosed
		var scholarshipLevel, scholarshipType string

		if strings.Contains(item.Title, "Licenciatura") {
			scholarshipLevel = Grad
		}
		if strings.Contains(item.Title, "Especialidade") {
			scholarshipLevel = Specialization
		}
		if strings.Contains(item.Title, "Mestrado") {
			scholarshipLevel = Masters
		}
		if strings.Contains(item.Title, "Doutoramento") {
			scholarshipLevel = Doctorate
		}

		if strings.Contains(item.Category, "Interna") {
			scholarshipType = InternalScholarship
		}
		if strings.Contains(item.Category, "Externa") {
			scholarshipType = ExternalScholarship
			if strings.Contains(item.Category, "Mérito") {
				scholarshipType = MeritScholarship
				status = StatusOpen
			}
		}

		start, err := time.Parse("20060102", item.Start) //YYYYMMDD
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		end, err := time.Parse("20060102", item.End) //YYYYMMDD
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		available, _ := strconv.Atoi(item.Available)

		requiredDocs, err := RequiredDocuments(scholarshipType, scholarshipLevel)
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		var u Scholarship
		if err := db.Model(&Scholarship{}).Where("name = ? AND sponsor_name = ? AND start = ? AND end = ?", item.Title, sponsor.Name, start, end).First(&u).Error; err != nil {
			u = Scholarship{
				Base:              Base{ID: ids[0]},
				Name:              item.Title,
				SponsorName:       sponsor.Name,
				Sponsor:           sponsor.ID,
				Content:           item.Content,
				Available:         int64(available),
				RequiredDocuments: requiredDocs,
				Start:             start,
				End:               end,
				Type:              scholarshipType,
				Level:             scholarshipLevel,
				Status:            status,
				MinGrade:          int64(minGrade),
				MaxAge:            int64(maxAge),
			}

			if err := zaplog.ZLog(db.Create(&u).Error); err != nil {
				continue
			}
		}

		scholarshipIds[item.ID] = u
	}

	return
}

func ImportApplications(db *gorm.DB, scholarshipIds map[string]Scholarship) (err error) {
	fmt.Println("Importing applications")
	if len(scholarshipIds) == 0 {
		return errors.New("No scholarships")
	}

	data := fmt.Sprintf("%s/import/applications.json", assetsDir)
	jsonFile, err := os.Open(data)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var items []ApplicationImport
	if err = zaplog.ZLog(json.Unmarshal(bytes, &items)); err != nil {
		return
	}

	for _, item := range items {
		var user User
		if err = zaplog.ZLog(db.Model(&User{}).Where("email = ?", item.Email).First(&user).Error); err != nil {
			fmt.Println(item.Email)
			continue
		}

		ids, err := util.GenerateUUIDS(2)
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		scholarship, ok := scholarshipIds[item.Scholarship]
		if !ok {
			continue
		}

		var count int
		if err = zaplog.ZLog(db.Model(&Application{}).Where("email = ? AND scholarship_name = ?", item.Email, scholarship.Name).Count(&count).Error); err == nil && count > 0 {
			continue
		}

		//1992-02-12 00:00:00
		birthDate, err2 := time.Parse("2006-01-02 03:04:05", item.BirthDate)
		if err2 != nil {
			birthDate, err2 = time.Parse("2006-01-02 03:04:05", Epoch)
			zaplog.ZLog(err2)
			// continue
		}
		gradYear, err := time.Parse("2006", item.GraduationYear)
		if err = zaplog.ZLog(err); err != nil {
			continue
		}

		score, _ := decimal.NewFromString(item.Score)
		scoreAge, _ := decimal.NewFromString(item.ScoreAge)
		scoreGrade, _ := decimal.NewFromString(item.ScoreGrade)
		familyMembers, _ := strconv.Atoi(item.FamilyMembers)
		workingMembers, _ := strconv.Atoi(item.WorkingMembers)
		educationGrade, _ := decimal.NewFromString(item.EducationAverage)
		children, _ := strconv.Atoi(item.Children)
		courseYear, _ := strconv.Atoi(item.CourseYear)
		average, _ := decimal.NewFromString(item.CourseAverage)
		salary, _ := decimal.NewFromString(item.Salary)
		status := StatusPending
		if item.Status == "" {
			status = item.Status
		}

		contact := item.ApplicationSchoolContact
		if contact == "" {
			contact = item.ApplicationSchoolContact1
		}

		var family []FamilyInfo

		if item.FatherName != "" {
			id, err := util.GenerateUUID()
			if err = zaplog.ZLog(err); err != nil {
				continue
			}
			familySalary, _ := decimal.NewFromString(item.FatherSalary)
			family = append(family, FamilyInfo{
				Base:         Base{ID: id},
				Application:  ids[0],
				Relation:     RelationFather,
				Name:         item.FatherName,
				BirthCountry: item.FatherCountry,
				Profession:   item.FatherProfession,
				Employer:     item.FatherEmployer,
				Salary:       familySalary,
			})
		}
		if item.MotherName != "" {
			id, err := util.GenerateUUID()
			if err = zaplog.ZLog(err); err != nil {
				continue
			}
			familySalary, _ := decimal.NewFromString(item.MotherSalary)
			family = append(family, FamilyInfo{
				Base:         Base{ID: id},
				Application:  ids[0],
				Relation:     RelationMother,
				Name:         item.MotherName,
				BirthCountry: item.MotherCountry,
				Profession:   item.MotherProfession,
				Employer:     item.MotherEmployer,
				Salary:       familySalary,
			})
		}
		if item.SpouseName != "" {
			id, err := util.GenerateUUID()
			if err = zaplog.ZLog(err); err != nil {
				continue
			}
			familySalary, _ := decimal.NewFromString(item.SpouseSalary)
			family = append(family, FamilyInfo{
				Base:         Base{ID: id},
				Application:  ids[0],
				Relation:     RelationSpouse,
				Name:         item.SpouseName,
				BirthCountry: item.SpouseCountry,
				Profession:   item.SpouseProfession,
				Employer:     item.SpouseEmployer,
				Salary:       familySalary,
			})
		}

		var docs []File
		if item.Certificate != "" {
			docs = append(docs, applicationFileFromURL(item.Certificate, ids[0], user.ID, user.Name, CertificateScan))
		}
		if item.CertificateValidation != "" {
			docs = append(docs, applicationFileFromURL(item.CertificateValidation, ids[0], user.ID, user.Name, CertificateValidationScan))
		}
		if item.HealthDocument != "" {
			docs = append(docs, applicationFileFromURL(item.HealthDocument, ids[0], user.ID, user.Name, HealthCertificate))
		}
		if item.IDScan != "" {
			docs = append(docs, applicationFileFromURL(item.IDScan, ids[0], user.ID, user.Name, IDScan))
		}
		if item.IESAppraisal != "" {
			docs = append(docs, applicationFileFromURL(item.IESAppraisal, ids[0], user.ID, user.Name, IESAppraisal))
		}
		if item.MilitaryDocument != "" {
			docs = append(docs, applicationFileFromURL(item.MilitaryDocument, ids[0], user.ID, user.Name, MilitaryDocument))
		}
		if item.PassportScan != "" {
			docs = append(docs, applicationFileFromURL(item.PassportScan, ids[0], user.ID, user.Name, PassportScan))
		}
		if item.PovertyProof != "" {
			docs = append(docs, applicationFileFromURL(item.PovertyProof, ids[0], user.ID, user.Name, PovertyDeclaration))
		}
		if item.Project != "" {
			docs = append(docs, applicationFileFromURL(item.Project, ids[0], user.ID, user.Name, ProjectProposal))
		}
		if item.RecommendationLetter != "" {
			docs = append(docs, applicationFileFromURL(item.RecommendationLetter, ids[0], user.ID, user.Name, RecommendationLetter))
		}
		if item.ResidenceProof != "" {
			docs = append(docs, applicationFileFromURL(item.ResidenceProof, ids[0], user.ID, user.Name, ResidenceProof))
		}
		if item.ReturnAgreement != "" {
			docs = append(docs, applicationFileFromURL(item.ReturnAgreement, ids[0], user.ID, user.Name, ReturnAgreement))
		}
		if item.RPE != "" {
			docs = append(docs, applicationFileFromURL(item.RPE, ids[0], user.ID, user.Name, RPEProof))
		}

		var languages Languages
		if item.English != "" || item.French != "" || item.Spanish != "" || item.Italian != "" || item.Mandarin != "" || item.German != "" || item.Russian != "" {
			languageID, err := util.GenerateUUID()
			if err = zaplog.ZLog(err); err != nil {
				continue
			}

			languages = Languages{
				Base:        Base{ID: languageID},
				Application: ids[0],
				French:      item.French,
				Spanish:     item.Spanish,
				English:     item.English,
				Italian:     item.Italian,
				German:      item.German,
				Russian:     item.Russian,
				Mandarin:    item.Mandarin,
			}
		}

		var educationLevel = HighSchool
		switch scholarship.Level {
		case Masters, Specialization:
			educationLevel = GradOther
		case Doctorate:
			educationLevel = PostGrad
		}

		province := item.ResidenceProvince
		if item.ResidenceProvince == "" && item.ResidenceProvinceAbroad != "" {
			province = item.ResidenceProvinceAbroad
		}

		u := Application{
			Base:               Base{ID: ids[0]},
			UserID:             user.ID,
			IDNumber:           item.IDNumber,
			PassportNumber:     item.Passport,
			Scholarship:        scholarship.ID,
			ScholarshipName:    scholarship.Name,
			Name:               item.Name,
			BirthDate:          birthDate,
			Gender:             item.Gender,
			BirthProvince:      item.BirthProvince,
			BirthCity:          item.BirthCity,
			CurrentCountry:     item.ResidenceCountry,
			CurrentProvince:    province,
			CurrentCity:        item.ResidenceCity,
			Phone:              item.Phone,
			Email:              item.Email,
			Skype:              item.Skype,
			BankName:           item.BankName,
			BankAccountNumber:  item.BankAccount,
			BankAccountOwner:   item.AccountOwner,
			IBAN:               item.IBAN,
			RPE:                item.RPE != "",
			Score:              score,
			ScoreGrades:        scoreGrade,
			ScoreAge:           scoreAge,
			ScoreJustification: item.ScoreJustification,
			FamilyMembers:      int64(familyMembers),
			Children:           int64(children),
			WorkingMembers:     int64(workingMembers),

			EducationLevel:            educationLevel,
			EducationSchool:           item.EducationSchool,
			EducationSchoolDepartment: item.EducationSchoolDepartment,
			EducationSchoolContact:    item.EducationSchoolContact,
			EducationCourse:           item.EducationCourse,
			EducationCountry:          item.EducationCountry,
			EducationProvince:         item.Educationprovince,
			GraduationDate:            gradYear,
			EducationGrade:            educationGrade,
			EducationEvaluation:       item.EducationEvaluation,

			ApplicationLevel: item.ApplicationLevel,
			CourseSchool:     item.ApplicationSchool,
			CourseDepartment: item.ApplicationDepartment,
			CourseOption1:    item.CourseOption1,
			CourseOption2:    item.CourseOption2,
			CourseOption3:    item.CourseOption3,
			CourseCountry:    item.ApplicationCountry,
			CourseProvince:   item.ApplicationProvince,
			CourseContact:    contact,
			CourseYear:       int64(courseYear),
			CourseAverage:    average,
			Family:           family,
			Employer:         item.Employer,
			EmployerAddress:  item.EmpployerAddress,
			EmployerContact:  item.EmployerContact,
			JobTitle:         item.JobTitle,
			Salary:           salary,
			Status:           status,
			Languages:        languages,
			Documents:        docs,
		}

		if u.EducationCountry == "" {
			u.EducationCountry = "Angola"
		}

		if u.CurrentProvince == "" && u.CurrentCity != "" {
			u.CurrentProvince = u.CurrentCity
			// if _, ok := Angola.Cities[u.CurrentCity]; ok {
			// 	u.CurrentProvince = u.CurrentCity
			// }
		}

		if u.Status == "" {
			u.Status = StatusPending
		}

		if err := zaplog.ZLog(db.Create(&u).Error); err != nil {
			continue
		}

		for _, doc := range docs {
			zaplog.ZLog(db.Create(&doc).Error)
		}

		if err = zaplog.ZLog(db.Exec("UPDATE scholarships SET total_applications = total_applications + 1 where uuid = ?", scholarship.ID).Error); err != nil {
			continue
		}
	}

	return
}

func applicationFileFromURL(url, applicationID, userID, userName, docName string) File {
	fileID, _ := util.GenerateUUID()
	return File{
		Base:       Base{ID: fileID},
		UserID:     userID,
		UserName:   userName,
		Name:       docName,
		Resource:   ResourceApplication,
		ResourceID: applicationID,
		Type:       docName,
		Extension:  "pdf",
		URL:        url,
		Status:     StatusApproved,
		Public:     true,
		Location:   "s3",
	}
}
