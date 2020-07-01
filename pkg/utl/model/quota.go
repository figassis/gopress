package model

import (
	"github.com/shopspring/decimal"
)

type (
	ProvinceQuota struct {
		Base
		Scholarship   string          `gorm:"unique"`
		Bengo         decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Benguela      decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Bie           decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Cabinda       decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		CuandoCubango decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		CuanzaNorte   decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		CuanzaSul     decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Cunene        decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Huambo        decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Huila         decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Luanda        decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		LundaNorte    decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		LundaSul      decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Malanje       decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Moxico        decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Namibe        decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Uige          decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
		Zaire         decimal.Decimal `json:",omitempty" sql:"type:decimal(5,2)"`
	}

	Statistic struct {
		Base
		Name       string          `gorm:"not null"`
		Resource   string          `gorm:"type:ENUM('Bolsa','Artigo');default:'Bolsa';not null"`
		ResourceID string          `gorm:"not null"`
		Type       string          `gorm:"type:ENUM('Global','Província');default:'Global';not null"`
		Resources  decimal.Decimal `sql:"type:decimal(9,3)" gorm:"not null"`
		Value      decimal.Decimal `sql:"type:decimal(9,3)" gorm:"not null"`
		Percentage decimal.Decimal `sql:"type:decimal(5,2)" gorm:"not null;"`
	}

	Public struct {
		Company             Organization
		Sponsors            []Organization
		Schools             []Organization
		Scholarships        []Scholarship
		CourseDomains       []CourseDomain
		Posts               []Post
		TotalApplications   decimal.Decimal
		AwardedScholarships decimal.Decimal
		TotalScholarships   decimal.Decimal
		Candidates          decimal.Decimal
		TotalSponsors       decimal.Decimal
		Courses             decimal.Decimal
	}
)

func (p ProvinceQuota) Valid() bool {
	total := p.Bengo.Add(p.Benguela).Add(p.Bie).Add(p.Cabinda).
		Add(p.CuandoCubango).Add(p.CuanzaNorte).Add(p.CuanzaSul).Add(p.Cunene).
		Add(p.Huambo).Add(p.Huila).Add(p.Luanda).Add(p.LundaNorte).Add(p.LundaSul).
		Add(p.Malanje).Add(p.Moxico).Add(p.Namibe).Add(p.Uige).Add(p.Zaire)

	return total.Equal(decimal.NewFromFloat(100)) || total.Equal(decimal.NewFromFloat(0))
}

func (p ProvinceQuota) Empty() bool {
	var temp ProvinceQuota
	return p == temp
}
