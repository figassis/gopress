package model

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/ses"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
)

const (
	sesConcurrency = 4
)

var (
	sender                                           *Sender
	messageTemplate, actionTemplate, contactTemplate *template.Template
)

type (
	Contact struct {
		Name     string
		Email    string
		Honeypot string
		Subject  string
		Message  string
	}
	Sender struct {
		ses          *ses.Email
		dailyUsage   float64
		dailyQuota   float64
		maxRate      float64
		sendInterval time.Duration
	}

	Notification struct {
		Base
		Type           string `gorm:"type:ENUM('Utilizador','Função','Candidatura','Marcação');default:'Utilizador'"`
		Targets        List   `gorm:"type:varchar(256)" sql:"type:varchar(256)"` //One or more IDs or roles, but not too many. For more, use a query
		Query          string
		Title          string
		Message        string
		Template       string `gorm:"type:ENUM('Mensagem','Acção');default:'Mensagem'"`
		Start          *time.Time
		Completed      bool
		RecipientCount int64 `gorm:"type:integer"`
		Sent           int64 `gorm:"type:integer"`
		Opened         int64 `gorm:"type:integer"`
		Bounced        int64 `gorm:"type:integer"`
		Unsubscribed   int64 `gorm:"type:integer"`
		Success        float64
		Status         string      `gorm:"type:ENUM('Activo','Concluído','A Processar');default:'Activo';not null"`
		Recipients     []Recipient `gorm:"foreignkey:Notification;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
	}

	Recipient struct {
		Base
		Notification    string
		Template        string `gorm:"type:ENUM('Mensagem','Acção');default:'Mensagem'"`
		UserID          string
		Email           string
		Name            string
		Message         string
		Subject         string
		ButtonLink      string
		ButtonTitle     string
		UnsubscribeLink string
		Domain          string
		Sent            bool
	}

	Bounce struct {
		Base
		Email string `gorm:"unique;not null"`
	}

	BounceNotification struct {
		NotificationType string
		Bounce           struct {
			BounceType        string
			BounceSubType     string
			BouncedRecipients []struct {
				EmailAddress string
			}
			Timestamp   time.Time
			FeedbackID  string
			RemoteMtaIP string
		}
		Complaint struct {
			ComplainedRecipients []struct {
				EmailAddress string
			}
			Timestamp  time.Time
			FeedbackID string
		}
		Mail struct {
			Timestamp        time.Time
			MessageID        string
			Source           string
			SourceArn        string
			SourceIP         string
			SendingAccountID string
			Destination      []string
			HeadersTruncated bool
			Headers          []struct {
				Name  string
				Value string
			}
			CommonHeaders struct {
				From      []string
				Date      string
				To        []string
				MessageID string
				Subject   string
			}
		}
	}
)

func NewMailer() (err error) {
	sender = &Sender{
		ses: ses.NewEmail(os.Getenv(AWS_ACCESS_KEY_ID), os.Getenv(AWS_SECRET_ACCESS_KEY), os.Getenv(AWS_REGION)),
	}
	sender.ses.SetupProfile("default", os.Getenv(SES_SENDER), []string{os.Getenv(SES_SENDER)}, "", "", "")

	messageTemplate, err = template.New("Message").Parse(MessageTemplate)
	if err != nil {
		return
	}

	actionTemplate, err = template.New("Action").Parse(ActionTemplate)
	if err != nil {
		return
	}

	contactTemplate, err = template.New("Contact").Parse(ContactTemplate)
	if err != nil {
		return
	}

	//Update limits on minute 5 of every hour
	if err = Schedule("5 * * * *", limits); err != nil {
		return zaplog.ZLog(err)
	}

	return
}

func (r Recipient) Send() (err error) {
	if sender == nil {
		return zaplog.ZLog("Email sender is not initialized")
	}

	if actionTemplate == nil || messageTemplate == nil || contactTemplate == nil {
		return zaplog.ZLog("Email templates are not initialized")
	}

	if r.Sent {
		return zaplog.ZLog("Email has already been sent")
	}

	if r.Email == "" {
		return zaplog.ZLog("Email address is required")
	}

	if r.Message == "" {
		return zaplog.ZLog("Email message is required")
	}

	if r.Subject == "" {
		return zaplog.ZLog("Email subject is required")
	}

	r.Domain = os.Getenv("FRONTEND")

	var buffer bytes.Buffer

	switch r.Template {
	case TemplateAction:
		if r.ButtonLink == "" || r.ButtonTitle == "" {
			return zaplog.ZLog("Missing Button link or title")
		}
		err = actionTemplate.Execute(&buffer, r)
	case TemplateMessage:
		err = messageTemplate.Execute(&buffer, r)
	case TemplateContact:
		err = contactTemplate.Execute(&buffer, r)
		r.Email = os.Getenv("CONTACT_EMAIL")
	default:
		return zaplog.ZLog("Invalid email template")
	}

	if err != nil {
		return zaplog.ZLog(err)
	}

	if err = sender.ses.Send("default", []string{r.Email}, []string{}, []string{}, r.Subject, buffer.String(), "", "UTF-8"); err != nil {
		return zaplog.ZLog(err)
	}

	return
}

func SendEmail(email, subject, message string) (err error) {
	return zaplog.ZLog(sender.ses.Send("default", []string{email}, []string{}, []string{}, subject, message, message, "UTF-8"))
}

func limits() {
	var err error
	sender.dailyUsage, sender.dailyQuota, sender.maxRate, sender.sendInterval, err = sender.ses.Limits()
	if err != nil {
		zaplog.ZLog(err)
	}
	return
}

func (n Notification) Process(db *gorm.DB) (err error) {
	var recipients []Recipient

	if err = db.Model(&Recipient{}).Where("notification = ? AND sent = ?", n.ID, false).Order("created_at ASC").Limit(100).Find(&recipients).Error; err != nil {
		return
	}

	chunks := funk.Chunk(recipients, sesConcurrency)
	shards, ok := chunks.([][]Recipient)
	if !ok {
		return zaplog.ZLog(errors.New("Invalid shards"))
	}

	for _, shard := range shards {
		tempShard := shard
		go processEmails(&tempShard, db)
	}

	return
}

func processEmails(recipients *[]Recipient, db *gorm.DB) (err error) {
	var emails []string
	var bounces []Bounce
	var unsub []User

	for _, r := range *recipients {
		emails = append(emails, r.Email)
	}

	if err = db.Model(&Bounce{}).Where("email IN (?)", emails).Find(&bounces).Error; err != nil {
		return
	}

	if err = db.Model(&User{}).Where("email IN (?) AND unsubscribed = ?", emails, true).Find(&unsub).Error; err != nil {
		return
	}

	var bounceMap = make(map[string]bool, len(bounces))
	for _, b := range bounces {
		bounceMap[b.Email] = true
	}
	for _, b := range unsub {
		bounceMap[b.Email] = true
	}

	var success []int64
	for _, r := range *recipients {

		//Do not send to bounced emails
		if bounceMap[r.Email] {
			success = append(success, r.IntID)
			continue
		}

		if err2 := zaplog.ZLog(r.Send()); err2 == nil {
			success = append(success, r.IntID)
		}
		time.Sleep(sender.sendInterval)
	}

	if err = db.Model(&Recipient{}).Where("id IN (?)", success).Update("sent", true).Error; err != nil {
		zaplog.ZLog(err)
	}
	return
}

func HandleBouncedEmail(n *BounceNotification, db *gorm.DB) (err error) {
	switch n.NotificationType {
	case "Bounce":
		for _, email := range n.Bounce.BouncedRecipients {
			id, _ := util.GenerateUUID()
			zaplog.ZLog(db.Create(&Bounce{Base: Base{ID: id}, Email: email.EmailAddress}).Error)
		}
	case "Complaint":
		for _, email := range n.Complaint.ComplainedRecipients {
			id, _ := util.GenerateUUID()
			zaplog.ZLog(db.Create(&Bounce{Base: Base{ID: id}, Email: email.EmailAddress}).Error)
		}
	}

	return
}

func SendContactEmail(contact Contact) (err error) {
	_, err = util.ValidateEmail(contact.Email)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	r := Recipient{
		Template: TemplateContact,
		Email:    contact.Email,
		Name:     contact.Name,
		Message:  contact.Message,
		Subject:  contact.Subject,
		Sent:     false,
	}

	go r.Send()
	return
}
