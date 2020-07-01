package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
)

// Pagination constants
const (
	paginationDefaultLimit = 20
	paginationMaxLimit     = 100
)

// Transform checks and converts http pagination into database pagination model
func (p *Pagination) DbPagination(db *gorm.DB) (limit int64, cursor Result, previous, next string) {

	maxLimit := int64(paginationMaxLimit)
	if p.CourseQuery != nil {
		maxLimit = 1500
		// p.Limit = int(maxLimit)
	}
	limit = int64(p.Limit)
	if limit < 1 {
		limit = paginationDefaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	var results []Result
	if p.Cursor == "" {
		rows, err := db.Order("id DESC").Limit(limit + 1).Select("uuid, id").Rows()
		defer rows.Close()

		if err = zaplog.ZLog(err); err != nil {
			return
		}

		for rows.Next() {
			var result Result
			rows.Scan(&result.ID, &result.IntID)
			results = append(results, result)
		}

		if len(results) == 0 {
			return
		}

		cursor = results[0]

		if len(results) > int(limit) {
			result := results[len(results)-1]
			// fmt.Printf("Length: %d, First: %v, Last: %v, Next: %v, Current: %v\n", len(results), results[0], results[len(results)-1], result, cursor)
			next = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%d", result.ID, result.IntID)))
		}
		return
	}

	decodedCursor, err := base64.URLEncoding.DecodeString(p.Cursor)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	cursorParts := strings.Split(string(decodedCursor), ",")
	if len(cursorParts) != 2 {
		zaplog.ZLog(errors.New("Invalid cursor"))
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(cursorParts[1]), 10, 64)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	cursor = Result{ID: cursorParts[0], IntID: id}

	rows, err := db.Where("(id,uuid) <= (?,?)", id, cursorParts[0]).Order("id DESC").Limit(limit + 1).Select("uuid, id").Rows()
	if err = zaplog.ZLog(err); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var result Result
		rows.Scan(&result.ID, &result.IntID)
		results = append(results, result)
	}

	if len(results) == 0 {
		return
	}

	cursor = results[0]

	if len(results) > int(limit) {
		result := results[len(results)-1]
		next = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%d", result.ID, result.IntID)))
		// fmt.Printf("Length: %d, First: %v, Last: %v, Next: %v, Current: %v\n", len(results), results[0], results[len(results)-1], result, cursor)
	}

	//Get Previous cursor
	rows, err = db.Where("(id,uuid) > (?,?)", cursor.IntID, cursor.ID).Order("id ASC").Limit(limit + 1).Select("uuid, id").Rows()
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	results = []Result{}
	for rows.Next() {
		var result Result
		rows.Scan(&result.ID, &result.IntID)
		results = append(results, result)
	}

	if len(results) == 0 {
		return
	}

	result := results[len(results)-1]
	previous = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%d", result.ID, result.IntID)))
	// fmt.Printf("Length: %d, First: %v, Last: %v, Previous: %v, Current: %v\n", len(results), results[0], results[len(results)-1], result, cursor)
	return

}

// Pagination holds paginations data
type (
	Pagination struct {
		Limit             int           `json:",omitempty" query:"limit" validate:"max=100"`
		Page              int           `json:",omitempty" query:"page" validate:"min=0"`
		Cursor            string        `json:",omitempty" query:"cursor"`
		CourseQuery       *Course       `json:"-"`
		UserQuery         *User         `json:"-"`
		OrganizationQuery *Organization `json:"-"`
		ScholarshipQuery  *Scholarship  `json:"-"`
		ApplicationQuery  *Application  `json:"-"`
		PostQuery         *Post         `json:"-"`
		AppointmentQuery  *Appointment  `json:"-"`
		SearchQuery       string        `json:"-"`
		CacheKey          string        `json:"-"`
	}

	ListResponse struct {
		Page          int            `json:",omitempty"`
		Users         []User         `json:",omitempty"`
		Organizations []Organization `json:",omitempty"`
		Files         []File         `json:",omitempty"`
		Scholarships  []Scholarship  `json:",omitempty"`
		Applications  []Application  `json:",omitempty"`
		Statistic     []Statistic    `json:",omitempty"`
		Posts         []Post         `json:",omitempty"`
		Notifications []Notification `json:",omitempty"`
		Recipients    []Recipient    `json:",omitempty"`
		Cities        []City         `json:",omitempty"`
		Courses       []Course       `json:",omitempty"`
		Appointments  []Appointment  `json:",omitempty"`
	}

	Result struct {
		ID    string
		IntID int64
	}
)
