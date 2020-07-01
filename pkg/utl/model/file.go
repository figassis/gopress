package model

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/figassis/goinagbe/pkg/utl/util"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	"github.com/jinzhu/gorm"
)

type (
	File struct {
		Base
		UserID     string
		UserName   string
		Name       string
		Resource   string `gorm:"type:ENUM('Bolsa','Organização','Artigo','Candidatura','Utilizador');default:'Candidatura';not null"`
		ResourceID string
		Type       string `gorm:"default:'Documento';not null"`
		Extension  string `gorm:"type:ENUM('bmp','jpg','png','pdf','doc','docx','xls','xlsx','csv','webp','json','txt','zip');default:'pdf';not null"`
		URL        string
		Data       []byte `json:"-" gorm:"-"`
		Status     string `gorm:"type:ENUM('Pendente','Aprovado','Rejeitado');default:'Aprovado';not null"`
		Comment    string
		Public     bool
		Path       string `json:"-" gorm:"-"`
		Issued     *time.Time
		Expires    *time.Time
		Location   string `json:"-" gorm:"type:ENUM('s3','local');default:'s3';not null"`
	}

	Upload struct {
		Base
		Path    string
		Public  bool
		Status  string `gorm:"type:ENUM('Pendente','Concluída','A Processar','Erro');default:'Pendente';not null"`
		Comment string
	}
)

func (f *File) Save(db *gorm.DB) (err error) {
	if f.UserID == "" || f.Path == "" || f.Resource == "" || f.ResourceID == "" || !DocumentTypes.Contains(f.Type) {
		return errors.New("Invalid file")
	}

	if f.ID == "" {
		if f.ID, err = util.GenerateUUID(); err != nil {
			return
		}
	}

	if f.URL, err = SaveFile(f.Path, f.ID); err != nil {
		return
	}

	return db.Create(f).Error
}

func (f *File) LoadData() (err error) {
	if f.Location == "s3" {
		key := strings.TrimPrefix(f.URL, fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv(BUCKET)))

		f.Data, err = Get(key)
		return
	}

	file, err := os.Open(f.URL)
	if err != nil {
		return
	}

	defer file.Close()

	f.Data, err = ioutil.ReadFile(f.URL)
	return
}

func GetS3URLFromID(id string, public bool) (url string) {
	return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", os.Getenv(BUCKET), GetS3KeyFromID(id, public))
}

func GetS3KeyFromID(id string, public bool) (url string) {
	if public {
		return fmt.Sprintf("%s/%s", os.Getenv(BUCKET_PUBLIC_PREFIX), id)
	}
	return fmt.Sprintf("%s/%s", os.Getenv(BUCKET_PRIVATE_PREFIX), id)
}

func SaveFile(filename string, key string) (url string, err error) {
	bucket := os.Getenv(BUCKET)
	region := os.Getenv(AWS_REGION)

	file, err := os.Open(filename)
	if err = zaplog.ZLog(err); err != nil {
		return "", errors.New(fmt.Sprintf("Failed to open %s", filename))
	}

	defer file.Close()

	fileData, err := ioutil.ReadFile(filename)
	if err = zaplog.ZLog(err); err != nil {
		return "", errors.New("Could not read file")
	}

	svc := s3manager.NewUploader(session.New(aws.NewConfig().WithRegion(region)))

	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(http.DetectContentType(fileData)),
	})

	if err = zaplog.ZLog(err); err != nil {
		return "", errors.New("Could not store file")
	}

	url = result.Location
	os.Remove(filename)
	return
}

func Get(key string) (data []byte, err error) {
	bucket := os.Getenv(BUCKET)
	region := os.Getenv(AWS_REGION)

	svc := s3manager.NewDownloader(session.New(aws.NewConfig().WithRegion(region)))

	results, err := svc.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return
	}

	defer results.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, results.Body); err != nil {
		return
	}

	return buf.Bytes(), nil
}

//GetPresignedURL returns a presigned url that expires in 24 hours
func (f File) GetURL() (url string, err error) {
	if f.Public {
		return f.URL, nil
	}

	bucket := os.Getenv(BUCKET)
	region := os.Getenv(AWS_REGION)

	svc := s3manager.NewDownloader(session.New(aws.NewConfig().WithRegion(region)))

	req, _ := svc.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(GetS3KeyFromID(f.ID, false)),
	})

	return req.Presign(time.Hour * 24)
}

func DeleteFiles(files *[]File) (err error) {
	bucket := os.Getenv(BUCKET)
	region := os.Getenv(AWS_REGION)

	svc := s3manager.NewBatchDelete(session.New(aws.NewConfig().WithRegion(region)))
	var objects = make([]s3manager.BatchDeleteObject, len(*files))
	for i, f := range *files {
		objects[i] = s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(GetS3KeyFromID(f.ID, false)),
			},
		}
	}

	if err = svc.Delete(aws.BackgroundContext(), &s3manager.DeleteObjectsIterator{Objects: objects}); err != nil {
		return
	}

	return nil
}
