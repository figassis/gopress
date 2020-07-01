package model

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/jinzhu/gorm"
)

const (
	SuperAdminRole              AccessRole = 100
	AdminRole                   AccessRole = 105
	OperatorRole                AccessRole = 110
	SupportRole                 AccessRole = 115
	CompanyAdminRole            AccessRole = 120
	CompanyUserRole             AccessRole = 125
	CandidateRole               AccessRole = 130
	confirmationEmailExpiration            = time.Hour * 24 * 7
)

var (
	roles = map[AccessRole]string{
		SuperAdminRole:   "Administrador Sistema",
		AdminRole:        "Administrador",
		OperatorRole:     "Oeprador",
		SupportRole:      "Suporte",
		CompanyAdminRole: "Administrador Parceiro",
		CompanyUserRole:  "Utilizador Perceiro",
		CandidateRole:    "Candidato",
	}

	userStatusFlows = map[string]List{
		StatusPending:   List{StatusActive},
		StatusActive:    List{StatusInactive, StatusSuspended},
		StatusSuspended: List{StatusActive, StatusInactive},
		StatusInactive:  List{StatusActive},
	}
)

type (

	// AccessRole represents access role type
	AccessRole int

	User struct {
		Base
		Name               string `gorm:"not null"`
		Username           string `gorm:"unique; not null"`
		Password           string `json:"-" gorm:"not null"`
		Email              string `gorm:"unique;not null"`
		Phone              string
		Status             string     `gorm:"type:ENUM('Pendente','Activo','Suspenso','Inactivo');default:'Activo';not null"`
		LastLogin          *time.Time `json:",omitempty"`
		LastPasswordChange *time.Time `json:",omitempty"`
		Token              string     `json:"-"`
		Role               AccessRole `gorm:"not null"`
		Organization       string     `gorm:"not null"`
		OrganizationName   string     `gorm:"not null"`
		Unsubscribed       bool
		UnsubscribeID      string
		Applications       []Application `json:",omitempty" gorm:"foreignkey:UserID;association_foreignkey:uuid;association_autoupdate:false"`
	}

	// AuthUser represents data stored in JWT token for user
	AuthUser struct {
		ID           string
		Organization string
		Username     string
		Email        string
		Role         AccessRole
	}
)

// ChangePassword updates user's password related fields
func (u *User) ChangePassword(hash string) {
	now := time.Now()
	u.Password = hash
	u.LastPasswordChange = &now
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin(token string) {
	now := time.Now()
	u.Token = token
	u.LastLogin = &now
}

func (r AccessRole) Exists() bool {
	return roles[r] != ""
}

func (r AccessRole) Name() string {
	return roles[r] //Super Admin
}

func (r AccessRole) Tag() string {
	return strings.ToLower(strings.Replace(r.Name(), " ", "_", -1)) //super_admin
}

func (u *User) SendConfirmationEmail() (err error) {
	if u.Unsubscribed {
		return nil
	}
	token, err := util.GenerateUUID()
	if err != nil {
		return
	}

	//FIXME: make confirmaton tokens more persistent
	if err = util.CacheTTL(token, u.ID, confirmationEmailExpiration); err != nil {
		return
	}

	r := Recipient{
		Template:        TemplateAction,
		Email:           u.Email,
		Name:            u.Name,
		Message:         "Bem vindo ao INAGBE Online! Para aceder a sua conta, por favor confirme o seu email clicando no botão abaixo.",
		Subject:         "Active a sua conta no INAGBE Online",
		ButtonLink:      fmt.Sprintf("https://%s/session/confirm/?token=%s", os.Getenv("FRONTEND"), token),
		ButtonTitle:     "Activar a conta",
		UnsubscribeLink: fmt.Sprintf("https://%s/session/unsubscribe/?token=%s", os.Getenv("FRONTEND"), u.ID),
		Sent:            false,
	}

	go r.Send()
	return nil
}

func CompleteEmailConfirmation(token string, db *gorm.DB) (err error) {
	var id string
	if err = util.GetCache(token, &id); err != nil {
		return errors.New("Invalid confirmation token")
	}

	if err = db.Model(&User{}).Where("uuid = ?", id).Update("status", StatusActive).Error; err != nil {
		return
	}

	util.DeleteCache(token)

	return
}

func UnsubscribeUser(token string, db *gorm.DB) (err error) {
	if err = db.Model(&User{}).Where("unsubscribe_id = ?", token).Update("unsubscribed", true).Error; err != nil {
		return
	}

	return
}

func (u *User) SendResetEmail() (err error) {
	token, err := util.GenerateUUID()
	if err != nil {
		return
	}

	if err = util.Cache(token, u.ID); err != nil {
		return
	}

	r := Recipient{
		Template:        TemplateAction,
		Email:           u.Email,
		Name:            u.Name,
		Message:         "Recebemos um pedido para restaurar a sua senha. Se não fez este pedido, por favor ignore este email.",
		Subject:         "Restaure a sua senha no INAGBE Online",
		ButtonLink:      fmt.Sprintf("https://%s/reset/%s", os.Getenv("FRONTEND"), token),
		ButtonTitle:     "Activar a conta",
		UnsubscribeLink: fmt.Sprintf("https://%s/unsubscribe/%s", os.Getenv("FRONTEND"), u.ID),
		Sent:            false,
	}

	go r.Send()
	return nil
}

func (a User) AllowedStatuses(newStatus string) bool {
	allowed, ok := userStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}
