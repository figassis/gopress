package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	PostCategories  = []string{CategoryGeneral, CategoryNews, CategoryAnnouncement}
	postStatusFlows = map[string]List{
		StatusDraft:     List{StatusPublished, StatusTrash},
		StatusPublished: List{StatusDraft, StatusTrash},
		StatusTrash:     List{StatusDraft},
	}
)

type (
	Post struct {
		Base
		Author     string `gorm:"not null"`
		AuthorName string `gorm:"not null"`
		Category   string `gorm:"not null;default:'Notícia'"`
		Tags       List   `gorm:"type:varchar(256)" sql:"type:varchar(256)"`
		Title      string `gorm:"not null"`
		Slug       string `gorm:"not null"`
		Template   string `gorm:"default:'post';not null"`
		Content    string `sql:"type:longtext"`
		Image      string
		html       string
		Excerpt    string `gorm:"not null"`
		Status     string `gorm:"type:ENUM('Rascunho','Publicado','Lixeira');default:'Rascunho';not null"`
	}
)

func (p Post) GetHtml() (html string, err error) {
	data, err := ioutil.ReadFile(os.Getenv("ASSETS") + "/html/" + p.ID)
	if err != nil {
		if err = p.generateHtml(); err != nil {
			return
		}
		return p.html, nil
	}

	return string(data), nil
}

func (p Post) generateHtml() (err error) {
	p.html, err = wysiwygToHTML(p.Content)
	if err != nil {
		return
	}

	return ioutil.WriteFile(os.Getenv("ASSETS")+"/html/"+p.ID, []byte(p.html), 0644)
}

func wysiwygToHTML(content string) (html string, err error) {
	//Convert all images and files to CDN
	public := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", os.Getenv("BUCKET"), os.Getenv("BUCKET_PUBLIC_PREFIX"))

	if os.Getenv("CDN") != "" {
		html = strings.ReplaceAll(html, public, strings.TrimSuffix(os.Getenv("CDN"), "/"))
	}
	return
}

func (a Post) AllowedStatuses(newStatus string) bool {
	allowed, ok := postStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}
