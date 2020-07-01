package orm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/jinzhu/gorm"
	"github.com/qor/validations"
	"go.uber.org/zap"

	// DB adapter

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var hostKey string

// New creates new database connection to a mysql database
func New(user, pass, host, port, database string, zlog *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, database))
	if err != nil {
		return nil, err
	}

	// Set max open connections
	// db.LogMode(true)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(100)
	db.DB().SetConnMaxLifetime(time.Hour * 8)
	db.SetLogger(util.NullLogger{})
	db.Callback().Create().Before("gorm:before_create").Register("generate_uuid", generateUUID)
	validations.RegisterCallbacks(db)
	db.SetLogger(gorm.Logger{log.New(os.Stdout, "\r\n", 0)})
	db = db.Set("gorm:auto_preload", true).BlockGlobalUpdate(true)

	hostKey, _ = util.GenerateUUID()
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func HostKey() string {
	return hostKey
}

func generateUUID(scope *gorm.Scope) {
	if !scope.HasError() {
		id, err := util.GenerateUUID()
		if scope.Err(err) != nil {
			return
		}

		primaryField := scope.PrimaryField()
		if primaryField.IsBlank {
			primaryField.Set(id)
		}
	}
}
