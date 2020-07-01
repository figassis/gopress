package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/figassis/goinagbe/pkg/utl/orm"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var (
	environment                    = []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "BUCKET", "BUCKET_PRIVATE_PREFIX", "BUCKET_PUBLIC_PREFIX", "AWS_REGION", "SES_SENDER", "FQDN", "COMPANY", "ADMIN_EMAIL", "ADMIN_PASSWORD", "ENVIRONMENT", "JWT_SECRET", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "APP_PORT", "FRONTEND", "CONTACT_EMAIL", "DATADIR"}
	homeDir, assetsDir, uploadsDir string
)

const (
	version = "0.0.5"
)

func Version() string {
	return version
}

// Load returns Configuration struct
func Load() (cfg *Configuration, err error) {
	log := zaplog.New(version)
	if err = zaplog.ZLog(validateEnvironment()); err != nil {
		return
	}

	homeDir = fmt.Sprintf("%s/%s", os.Getenv("DATADIR"), "goinagbe")
	assetsDir = fmt.Sprintf("%s/%s/assets", os.Getenv("DATADIR"), "goinagbe")
	uploadsDir = fmt.Sprintf("%s/%s/uploads", os.Getenv("DATADIR"), "goinagbe")

	cfg = &Configuration{
		Version: version,
		Assets:  assetsDir,
		Uploads: uploadsDir,
		Log:     log,
		Server: &Server{
			Port:         os.Getenv("APP_PORT"),
			Debug:        os.Getenv("ENVIRONMENT") != "production",
			ReadTimeout:  30, //seconds
			WriteTimeout: 30, //seconds
		},
		AWS: &AWS{
			Key:           os.Getenv("AWS_ACCESS_KEY_ID"),
			Secret:        os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Bucket:        os.Getenv("BUCKET"),
			PublicPrefix:  os.Getenv("BUCKET_PUBLIC_PREFIX"),
			PrivatePrefix: os.Getenv("BUCKET_PRIVATE_PREFIX"),
			Region:        os.Getenv("AWS_REGION"),
			Sender:        os.Getenv("SES_SENDER"),
		},
		DB: &Database{
			Host:       os.Getenv("DB_HOST"),
			Port:       os.Getenv("DB_PORT"),
			User:       os.Getenv("DB_USER"),
			Password:   os.Getenv("DB_PASSWORD"),
			Database:   os.Getenv("DB_NAME"),
			LogQueries: os.Getenv("ENVIRONMENT") != "production",
		},
		JWT: &JWT{
			Secret:           os.Getenv("JWT_SECRET"),
			Duration:         60,     //1 hour
			RefreshDuration:  43200,  //1 month
			MaxRefresh:       129600, //3 months
			SigningAlgorithm: "HS256",
		},
		App: &Application{
			FQDN:           os.Getenv("FQDN"),
			Frontend:       os.Getenv("FRONTEND"),
			Company:        os.Getenv("COMPANY"),
			AdminEmail:     os.Getenv("ADMIN_EMAIL"),
			AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
			Environment:    os.Getenv("ENVIRONMENT"),
			MinPasswordStr: 1,
			SwaggerUIPath:  "assets/swaggerui",
		},
	}

	if err = zaplog.ZLog(EnsureDirs([]string{homeDir, assetsDir, uploadsDir}, 0700)); err != nil {
		return
	}

	if err = zaplog.ZLog(os.Setenv("UPLOADS", uploadsDir)); err != nil {
		return
	}

	if err = zaplog.ZLog(os.Setenv("ASSETS", assetsDir)); err != nil {
		return
	}

	cfg.DB.Db, err = orm.New(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database, cfg.Log.GetLogger())
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	//Initialize Cache
	if err = zaplog.ZLog(util.New()); err != nil {
		return
	}

	return
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Version string
	Assets  string
	Uploads string
	Log     *zaplog.Log
	Server  *Server
	DB      *Database
	JWT     *JWT
	App     *Application
	AWS     *AWS
}

type AWS struct {
	Key           string
	Secret        string
	Bucket        string
	PublicPrefix  string
	PrivatePrefix string
	Region        string
	Sender        string
}

// Database holds data necessery for database configuration
type Database struct {
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	LogQueries bool
	Db         *gorm.DB
}

// Server holds data necessery for server configuration
type Server struct {
	Port              string
	Debug             bool
	ReadTimeout       int
	WriteTimeout      int
	RequestsPerSecond int
}

// JWT holds data necessery for JWT configuration
type JWT struct {
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}

// Application holds application configuration details
type Application struct {
	FQDN           string
	Frontend       string
	Company        string
	AdminEmail     string
	AdminPassword  string
	Environment    string
	MinPasswordStr int
	SwaggerUIPath  string
}

func validateEnvironment() error {
	//If we have an .env fle, use that, otherwise, use existing environment
	if err := CheckPath(".env"); err == nil {
		if err = zaplog.ZLog(godotenv.Load()); err != nil {
			return zaplog.ZLog("Could not load env file")
		}
	}

	for _, env := range environment {
		if os.Getenv(env) == "" {
			return fmt.Errorf("Environment variable %s is required", env)
		}
	}
	return nil
}

func EnsureDir(path string, perms os.FileMode) (err error) {
	if err = CheckPath(path); err == nil {
		return
	}

	return zaplog.ZLog(makeDir(path, perms))

}

func EnsureDirs(paths []string, perms os.FileMode) (err error) {
	for _, path := range paths {
		if err = EnsureDir(path, perms); err != nil {
			return
		}
	}
	return
}

func CheckPath(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s does not exist", path))
	}
	return nil
}

func makeDir(dir string, perms os.FileMode) (err error) {
	if err = os.MkdirAll(dir, perms); err != nil {
		return zaplog.ZLog(errors.New("Could not create directory" + dir))
	}

	return
}
