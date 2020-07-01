package sql

import (
	"fmt"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

// NewUser returns a new user database instance
func New() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id string) (user *model.User, err error) {
	user = new(model.User)
	err = db.Where("uuid = ?", id).First(user).Error
	err = zaplog.ZLog(err)
	return
}

func (u *User) GetPublicData(db *gorm.DB) (data *model.Public, err error) {
	var company model.Organization
	var sponsors []model.Organization
	var schools []model.Organization
	var domains []model.CourseDomain
	var scholarships []model.Scholarship
	var posts []model.Post

	var TotalApplications, AwardedScholarships, Candidates, CourseStats decimal.Decimal

	if err = db.Where("type = ?", model.OrgMain).First(&company).Error; err != nil {
		return
	}

	if err = db.Where("type = ? AND status = ?", model.OrgSponsor, model.StatusActive).Order("id DESC").Limit(100).Find(&sponsors).Error; err != nil {
		return
	}

	sponsors = append(sponsors, company)

	if err = db.Where("type = ? AND status = ?", model.OrgSchool, model.StatusActive).Order("id DESC").Limit(300).Find(&schools).Error; err != nil {
		return
	}

	if err = db.Model(&model.CourseDomain{}).Find(&domains).Error; err != nil {
		return
	}

	if err = db.Where("status = ?", model.StatusPublished).Order("created_at DESC").Limit(20).Find(&posts).Error; err != nil {
		return
	}

	if err = db.Where("status = ? AND end > ?", model.StatusOpen, time.Now()).Order("id DESC").Limit(100).Find(&scholarships).Error; err != nil {
		return
	}

	if err = db.Model(&model.Application{}).Where("status = ?", model.StatusPending).Count(&TotalApplications).Error; err != nil {
		return
	}

	if err = db.Model(&model.Application{}).Where("status = ?", model.StatusAwarded).Count(&AwardedScholarships).Error; err != nil {
		return
	}

	if err = db.Model(&model.User{}).Where("status = ? AND role = ?", model.StatusActive, model.CandidateRole).Count(&Candidates).Error; err != nil {
		return
	}

	if err = db.Model(&model.Course{}).Count(&CourseStats).Error; err != nil {
		return
	}

	err = zaplog.ZLog(err)

	data = &model.Public{
		Company:             company,
		Sponsors:            sponsors,
		Schools:             schools,
		Scholarships:        scholarships,
		CourseDomains:       domains,
		TotalApplications:   TotalApplications,
		AwardedScholarships: AwardedScholarships,
		TotalScholarships:   decimal.NewFromFloat(float64(len(scholarships))),
		Candidates:          Candidates,
		TotalSponsors:       decimal.NewFromFloat(float64(len(sponsors))),
		Courses:             CourseStats,
		Posts:               posts,
	}

	fmt.Printf("%d Sponsors, %d Schools, %d Scholarships, %d Domains, %d Posts\n", len(data.Sponsors), len(data.Schools), len(data.Scholarships), len(data.CourseDomains), len(data.Posts))
	return
}

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db *gorm.DB, uname string) (user *model.User, err error) {
	user = new(model.User)
	err = db.Where("username = ?", uname).First(user).Error
	err = zaplog.ZLog(err)
	return
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db *gorm.DB, token string) (user *model.User, err error) {
	user = new(model.User)
	err = db.Where("token = ?", token).First(user).Error
	err = zaplog.ZLog(err)
	return
}

// Update updates user's info
func (u *User) Update(db *gorm.DB, user *model.User) error {
	now := time.Now()
	// data := map[string]interface{}{"status": user.Status, "role": user.Role}
	update := model.User{Status: user.Status, Role: user.Role}

	if user.Password != "" {
		update.Password, update.LastPasswordChange = user.Password, &now
	}

	if user.Token != "" {
		update.Token = user.Token
	}

	return zaplog.ZLog(db.Model(&user).Updates(update).Error)
}

func (u *User) Signup(db *gorm.DB, user *model.User) error {

	user.Status = model.StatusPending
	if err := db.Create(user).Error; err != nil {
		return zaplog.ZLog(err)
	}

	return zaplog.ZLog(user.SendConfirmationEmail())
}

func (u *User) Unsubscribe(db *gorm.DB, token string) (err error) {
	return zaplog.ZLog(model.UnsubscribeUser(token, db))
}

func (u *User) ConfirmEmail(db *gorm.DB, token string) (err error) {
	return zaplog.ZLog(model.CompleteEmailConfirmation(token, db))
}

func (u *User) Bounce(db *gorm.DB, n model.BounceNotification) (err error) {
	go model.HandleBouncedEmail(&n, db)
	return nil
}
