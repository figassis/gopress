package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brianvoe/gofakeit"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
)

type (
	Course struct {
		Base
		Name        string `gorm:"unique_index:school_course"`
		School      string `gorm:"unique_index:school_course"`
		SchoolName  string
		Department  string `gorm:"unique_index:school_course"`
		Domain      string `gorm:"not null"`
		DomainName  string `gorm:"not null"`
		Cluster     string `gorm:"not null"`
		ClusterName string `gorm:"not null"`
		Type        string `gorm:"type:ENUM('Normal','Prioritário');default:'Normal';not null"`
		Level       string `gorm:"type:ENUM('Ensino Médio','Licenciatura','Mestrado','Doutoramento','Especialidade');default:'Licenciatura';unique_index:school_course;not null"`
	}

	CourseDomain struct {
		Base
		Name     string          `gorm:"unique"`
		Clusters []CourseCluster `gorm:"foreignkey:Domain;association_foreignkey:uuid;association_autoupdate:false"`
		Courses  []Course        `gorm:"foreignkey:Domain;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
	}

	CourseCluster struct {
		Base
		Name       string `gorm:"unique_index:domain_cluster"`
		Domain     string `gorm:"unique_index:domain_cluster"`
		DomainName string
		Courses    []Course `gorm:"foreignkey:Cluster;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
	}

	MeritImport struct {
		Domain  string
		Cluster string
		Course  string
	}

	CourseImport struct {
		Institution string
		Department  string
		Level       string
		Course      string
	}
)

func initializeCourseDomains(db *gorm.DB) (err error) {
	importCourses := false
	domain := CourseDomain{}
	if err = db.Model(&CourseDomain{}).Where("name = ?", "Geral").First(&domain).Error; err != nil {
		importCourses = true
		domain = CourseDomain{Name: "Geral"}
		if err = zaplog.ZLog(db.Create(&domain).Error); err != nil {
			return
		}
	}

	cluster := CourseCluster{}
	if err = db.Model(&CourseCluster{}).Where("name = ? AND domain = ?", "Geral", domain.ID).First(&cluster).Error; err != nil {
		importCourses = true
		cluster = CourseCluster{Name: "Geral", Domain: domain.ID, DomainName: domain.Name}
		if err = zaplog.ZLog(db.Create(&cluster).Error); err != nil {
			return
		}
	}

	if importCourses {
		if err = zaplog.ZLog(ImportMeritCourses(db)); err != nil {
			return
		}
		if err = zaplog.ZLog(ImportCourses(db)); err != nil {
			return
		}
	}

	return
}

func GetDefaultDomainAndCluster(db *gorm.DB) (domain CourseDomain, cluster CourseCluster, err error) {
	if err = zaplog.ZLog(db.Model(&CourseDomain{}).Where("name = ?", "Geral").First(&domain).Error); err != nil {
		return
	}

	if err = zaplog.ZLog(db.Model(&CourseCluster{}).Where("name = ? AND domain = ?", "Geral", domain.ID).First(&cluster).Error); err != nil {
		return
	}
	return
}

func ImportMeritCourses(db *gorm.DB) (err error) {
	merit := fmt.Sprintf("%s/merit.json", assetsDir)
	jsonFile, err := os.Open(merit)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var courses []MeritImport
	var domains, clusters = make(map[string]string), make(map[string]string)
	if err = zaplog.ZLog(json.Unmarshal(bytes, &courses)); err != nil {
		return
	}

	for _, course := range courses {
		domain := CourseDomain{Name: course.Domain}
		if err = zaplog.ZLog(db.Where(domain).FirstOrCreate(&domain).Error); err != nil {
			continue
		}

		domains[domain.Name] = domain.ID
	}

	for _, course := range courses {
		domainID, ok := domains[course.Domain]
		if !ok {
			continue
		}

		cluster := CourseCluster{Name: course.Cluster, Domain: domainID, DomainName: course.Domain}
		if err = zaplog.ZLog(db.Where(cluster).FirstOrCreate(&cluster).Error); err != nil {
			continue
		}

		clusters[fmt.Sprintf("%s:%s", course.Domain, course.Cluster)] = cluster.ID
	}

	for _, course := range courses {
		domainID, ok := domains[course.Domain]
		if !ok {
			continue
		}

		clusterID, ok := clusters[fmt.Sprintf("%s:%s", course.Domain, course.Cluster)]
		if !ok {
			continue
		}

		course1 := Course{Name: course.Course, Domain: domainID, DomainName: course.Domain, Cluster: clusterID, ClusterName: course.Cluster, Type: PriorityCourse, Level: Masters}
		course2 := Course{Name: course.Course, Domain: domainID, DomainName: course.Domain, Cluster: clusterID, ClusterName: course.Cluster, Type: PriorityCourse, Level: Doctorate}
		course3 := Course{Name: course.Course, Domain: domainID, DomainName: course.Domain, Cluster: clusterID, ClusterName: course.Cluster, Type: PriorityCourse, Level: Specialization}
		if err = zaplog.ZLog(db.Where(course1).FirstOrCreate(&course1).Error); err != nil {
			continue
		}
		if err = zaplog.ZLog(db.Where(course2).FirstOrCreate(&course2).Error); err != nil {
			continue
		}
		if err = zaplog.ZLog(db.Where(course3).FirstOrCreate(&course3).Error); err != nil {
			continue
		}
	}

	return
}

func ImportCourses(db *gorm.DB) (err error) {
	data := fmt.Sprintf("%s/courses.json", assetsDir)
	jsonFile, err := os.Open(data)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	var courses []CourseImport
	var schools = make(map[string]string)
	if err = zaplog.ZLog(json.Unmarshal(bytes, &courses)); err != nil {
		return
	}

	domain, cluster, err := GetDefaultDomainAndCluster(db)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	for _, course := range courses {
		school := Organization{
			Name:     course.Institution,
			Status:   StatusActive,
			Type:     OrgSchool,
			Country:  "Angola",
			Province: "Luanda",
			City:     "Luanda",
			Phone:    gofakeit.Phone(),
			Email:    gofakeit.Email(),
			Logo:     "assets/images/thumb-1.jpg",
		}

		if err = zaplog.ZLog(db.Where(school).FirstOrCreate(&school).Error); err != nil {
			continue
		}

		schools[course.Institution] = school.ID
	}

	for _, course := range courses {
		schoolID, ok := schools[course.Institution]
		if !ok {
			continue
		}

		if course.Level == "Pós Graduação" {
			for _, level := range []string{Masters, Doctorate, Specialization} {
				course1 := Course{Name: course.Course, Domain: domain.ID, DomainName: domain.Name, Cluster: cluster.ID, ClusterName: cluster.Name, Type: NormalCourse, Level: level, School: schoolID, SchoolName: course.Institution, Department: course.Department}
				if err = zaplog.ZLog(db.Where(course1).FirstOrCreate(&course1).Error); err != nil {
					continue
				}
			}
		} else {
			course1 := Course{Name: course.Course, Domain: domain.ID, DomainName: domain.Name, Cluster: cluster.ID, ClusterName: cluster.Name, Type: NormalCourse, Level: course.Level, School: schoolID, SchoolName: course.Institution, Department: course.Department}
			if err = zaplog.ZLog(db.Where(course1).FirstOrCreate(&course1).Error); err != nil {
				continue
			}
		}

	}

	return
}
