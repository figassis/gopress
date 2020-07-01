package api

import (
	"crypto/sha1"

	"github.com/figassis/goinagbe/pkg/api/auth"
	authlog "github.com/figassis/goinagbe/pkg/api/auth/logging"
	authtransport "github.com/figassis/goinagbe/pkg/api/auth/transport"

	"github.com/figassis/goinagbe/pkg/api/password"
	passlog "github.com/figassis/goinagbe/pkg/api/password/logging"
	passtransport "github.com/figassis/goinagbe/pkg/api/password/transport"

	"github.com/figassis/goinagbe/pkg/api/user"
	userlog "github.com/figassis/goinagbe/pkg/api/user/logging"
	usertransport "github.com/figassis/goinagbe/pkg/api/user/transport"

	"github.com/figassis/goinagbe/pkg/api/organization"
	organizationlog "github.com/figassis/goinagbe/pkg/api/organization/logging"
	organizationtransport "github.com/figassis/goinagbe/pkg/api/organization/transport"

	"github.com/figassis/goinagbe/pkg/api/application"
	applicationlog "github.com/figassis/goinagbe/pkg/api/application/logging"
	applicationtransport "github.com/figassis/goinagbe/pkg/api/application/transport"

	"github.com/figassis/goinagbe/pkg/api/appointment"
	appointmentlog "github.com/figassis/goinagbe/pkg/api/appointment/logging"
	appointmenttransport "github.com/figassis/goinagbe/pkg/api/appointment/transport"

	"github.com/figassis/goinagbe/pkg/api/scholarship"
	scholarshiplog "github.com/figassis/goinagbe/pkg/api/scholarship/logging"
	scholarshiptransport "github.com/figassis/goinagbe/pkg/api/scholarship/transport"

	"github.com/figassis/goinagbe/pkg/api/course"
	courselog "github.com/figassis/goinagbe/pkg/api/course/logging"
	coursetransport "github.com/figassis/goinagbe/pkg/api/course/transport"

	"github.com/figassis/goinagbe/pkg/api/media"
	medialog "github.com/figassis/goinagbe/pkg/api/media/logging"
	mediatransport "github.com/figassis/goinagbe/pkg/api/media/transport"

	"github.com/figassis/goinagbe/pkg/api/post"
	postlog "github.com/figassis/goinagbe/pkg/api/post/logging"
	posttransport "github.com/figassis/goinagbe/pkg/api/post/transport"

	"github.com/figassis/goinagbe/pkg/utl/config"
	"github.com/figassis/goinagbe/pkg/utl/middleware/jwt"
	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/rbac"
	"github.com/figassis/goinagbe/pkg/utl/secure"
	"github.com/figassis/goinagbe/pkg/utl/server"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
)

// Start starts the API service
func Start(cfg *config.Configuration) (err error) {
	log := cfg.Log

	if err = zaplog.ZLog(model.AutoMigrate(cfg.DB.Db)); err != nil {
		return
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.New()
	jwt := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration, cfg.JWT.RefreshDuration)
	if err = zaplog.ZLog(model.NewMailer()); err != nil {
		return
	}

	model.StartQueue(cfg.DB.Db)
	e := server.New()
	// e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authtransport.NewHTTP(authlog.New(auth.Initialize(cfg.DB.Db, jwt, sec, rbac), log), e, jwt.MWFunc())

	v1 := e.Group("/v1")
	v1.Use(jwt.MWFunc())

	usertransport.NewHTTP(userlog.New(user.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	passtransport.NewHTTP(passlog.New(password.Initialize(cfg.DB.Db, rbac, sec), log), v1)
	organizationtransport.NewHTTP(organizationlog.New(organization.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	applicationtransport.NewHTTP(applicationlog.New(application.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	scholarshiptransport.NewHTTP(scholarshiplog.New(scholarship.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	coursetransport.NewHTTP(courselog.New(course.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	posttransport.NewHTTP(postlog.New(post.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	appointmenttransport.NewHTTP(appointmentlog.New(appointment.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)
	mediatransport.NewHTTP(medialog.New(media.Initialize(cfg.DB.Db, cfg.App, rbac, sec), log), v1)

	if err = zaplog.ZLog(model.Initialize(cfg.DB.Db)); err != nil {
		return
	}

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return
}
