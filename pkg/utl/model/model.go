package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
)

// Base contains common fields for all tables
type (
	Base struct {
		IntID     int64      `gorm:"column:id;type:bigint(20) unsigned auto_increment;unique;not null;index:idx_id_uuid" json:"-"`
		ID        string     `gorm:"column:uuid;primary_key;index:idx_id_uuid" json:"ID,omitempty"`
		CreatedAt time.Time  `json:"created,omitempty" gorm:"not null"`
		UpdatedAt time.Time  `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
		DeletedAt *time.Time `gorm:"index" json:"-"`
	}

	// ListQuery holds company/location data used for list db queries
	ListQuery struct {
		Query string
		ID    string
	}

	List []string
)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeSave() error {
	if b.ID == "" {
		uuid, err := util.GenerateUUID()
		if err != nil {
			return err
		}
		b.ID = uuid
	}
	return nil
}

func (l List) Contains(key string) bool {
	for _, value := range l {
		if value == key {
			return true
		}
	}
	return false
}

func (j List) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}

	result := strings.Join([]string(j), "|")
	if strings.Contains(result, "||") {
		return nil, errors.New("input contains invalid character |")
	}
	return result, nil
}

func (j *List) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}

	switch s := value.(type) {
	case []string:
		*j = s
	case string:
		*j = strings.Split(s, "|")
	case []byte:
		*j = strings.Split(string(s), "|")
	default:
		return fmt.Errorf("Invalid Scan Source")
	}

	return nil
}

func (m List) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(m))
}
func (m *List) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("null pointer exception")
	}

	var tempResult []string
	if err := json.Unmarshal(data, &tempResult); err != nil {
		return err
	}
	result := List(tempResult)
	m = &result
	return nil
}

func (j List) IsNull() bool {
	if len(j) == 0 {
		return true
	}

	for _, str := range []string(j) {
		if str != "" && str != "|" {
			return false
		}
	}
	return true
}

func (j List) Equals(j1 List) bool {
	return strings.Join(j, "|") == strings.Join(j1, "|")
}

func AutoMigrate(db *gorm.DB) (err error) {
	// if os.Getenv("RESET") == "true" {
	// 	db.DropTable(&User{})
	// 	db.DropTable(&Organization{})
	// 	db.DropTable(&File{})
	// 	db.DropTable(&Scholarship{})
	// 	db.DropTable(&Application{})
	// 	db.DropTable(&ProvinceQuota{})
	// 	db.DropTable(&Statistic{})
	// 	db.DropTable(&Post{})
	// 	db.DropTable(&Notification{})
	// 	db.DropTable(&Recipient{})
	// 	db.DropTable(&City{})
	// 	db.DropTable(&Course{})
	// 	db.DropTable(&Languages{})
	// 	db.DropTable(&FamilyInfo{})
	// 	db.DropTable(&Appointment{})
	// 	db.DropTable(&Upload{})
	// 	db.DropTable(&Bounce{})
	// 	db.DropTable(&CourseDomain{})
	// 	db.DropTable(&CourseCluster{})
	// }

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Scholarship{})
	db.AutoMigrate(&Application{})
	db.AutoMigrate(&ProvinceQuota{})
	db.AutoMigrate(&Statistic{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Notification{})
	db.AutoMigrate(&Recipient{})
	db.AutoMigrate(&City{})
	db.AutoMigrate(&Course{})
	db.AutoMigrate(&Languages{})
	db.AutoMigrate(&FamilyInfo{})
	db.AutoMigrate(&Appointment{})
	db.AutoMigrate(&Upload{})
	db.AutoMigrate(&Bounce{})
	db.AutoMigrate(&CourseDomain{})
	db.AutoMigrate(&CourseCluster{})
	return
}

func Initialize(db *gorm.DB) (err error) {
	assetsDir = fmt.Sprintf("%s/%s/assets", os.Getenv("DATADIR"), "goinagbe")

	if err = zaplog.ZLog(initializeCourseDomains(db)); err != nil {
		return
	}
	// zaplog.ZLog(GenerateFakeData(10, db))

	if os.Getenv("RESET") == "true" {
		if err = zaplog.ZLog(ImportPosts(db)); err != nil {
			return
		}

		if err = zaplog.ZLog(ImportUsers(db)); err != nil {
			return
		}

		scholarships, err2 := ImportScholarships(db)
		if err = zaplog.ZLog(err2); err != nil {
			return
		}

		if err = zaplog.ZLog(ImportApplications(db, scholarships)); err != nil {
			return
		}
	}

	return
}
