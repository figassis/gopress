package model

import (
	"log"
	"os"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
	cron "github.com/robfig/cron/v3"
)

const (
	queueSpeed       = time.Millisecond * 200
	scholarshipSpeed = time.Second * 30
)

var (
	schedule = cron.New()
)

func StartQueue(d *gorm.DB) {
	db := d.Set("role", "batch")
	db.SetLogger(util.NullLogger{})

	if db.Callback().Create().Get("validations:validate") != nil {
		db.Callback().Create().Before("gorm:before_create").Remove("validations:validate")
	}

	if db.Callback().Update().Get("validations:validate") != nil {
		db.Callback().Update().Before("gorm:before_update").Remove("validations:validate")
	}

	db.SetLogger(gorm.Logger{log.New(os.Stdout, "\r\n", 0)})
	schedule.Start()
	go notificationQueue(db)
	go uploadQueue(db)
	go scholarshipFlow(db)
}

func notificationQueue(db *gorm.DB) {
	for {
		time.Sleep(queueSpeed)
		zaplog.ZLog(processNotifications(db))
	}
}

func uploadQueue(db *gorm.DB) {
	for {
		time.Sleep(queueSpeed)
		zaplog.ZLog(processUploads(db))
	}
}

func scholarshipFlow(db *gorm.DB) {
	for {
		time.Sleep(scholarshipSpeed)
		zaplog.ZLog(processScholarships(db))
	}
}

func processScholarships(db *gorm.DB) (err error) {
	now := time.Now()
	if err = db.Model(&Scholarship{}).Where("status = ? AND start <= ? AND end > ?", StatusPublished, now, now).Update("status", StatusOpen).Error; err != nil {
		return
	}
	if err = db.Model(&Scholarship{}).Where("status = ? AND end <= ?", StatusOpen, now).Update("status", StatusClosed).Error; err != nil {
		return
	}
	return
}

func processNotifications(db *gorm.DB) (err error) {

	var notifications []Notification
	var limit = 20

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Model(&Notification{}).Where("status = ?", StatusActive).Order("created_at asc").Limit(limit).Find(&notifications).Error; err != nil {
		return
	}

	ids := make([]int64, len(notifications))
	for _, n := range notifications {
		ids = append(ids, n.IntID)
	}

	if err = tx.Model(&Notification{}).Where("id IN (?)", ids).Update("status", StatusProcessing).Error; err != nil {
		return
	}

	tx.Commit()

	for _, n := range notifications {
		zaplog.ZLog(n.Process(db))
	}
	return
}

func processUploads(db *gorm.DB) (err error) {

	var uploads []Upload
	var limit = 20

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Model(&Upload{}).Where("status = ?", StatusPending).Order("created_at asc").Limit(limit).Find(&uploads).Error; err != nil {
		return
	}

	ids := []int64{}
	for _, n := range uploads {
		ids = append(ids, n.IntID)
	}

	if err = tx.Model(&Upload{}).Where("id IN (?)", ids).Update("status", StatusProcessing).Error; err != nil {
		return
	}

	tx.Commit()

	for _, u := range uploads {
		url, err2 := SaveFile(u.Path, GetS3KeyFromID(u.ID, u.Public))
		if err2 = zaplog.ZLog(err2); err2 != nil {
			zaplog.ZLog(db.Model(&Upload{}).Where("uuid = ?", u.ID).Update(Upload{Status: StatusError, Comment: err2.Error()}).Error)
			continue
		}

		if err2 = db.Model(&File{}).Where("uuid = ?", u.ID).Update("url", url).Error; err2 != nil {
			zaplog.ZLog(err)
		}

		zaplog.ZLog(db.Model(&Upload{}).Where("uuid = ?", u.ID).Update("status", StatusCompleted).Error)
	}
	return
}

func Schedule(cronSchedule string, cmd func()) (err error) {
	_, err = schedule.AddFunc(cronSchedule, cmd)
	return
}
