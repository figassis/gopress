package model

var (
	orgStatusFlows = map[string]List{
		StatusActive:    List{StatusInactive, StatusSuspended},
		StatusSuspended: List{StatusActive, StatusInactive},
		StatusInactive:  List{StatusActive},
	}
)

// Company represents company model
type (
	Organization struct {
		Base
		Name      string `gorm:"not null"`
		Status    string `gorm:"type:ENUM('Activo','Suspenso','Inactivo');default:'Activo';not null"`
		Type      string `gorm:"type:ENUM('Escola','Patrocinador','Principal');default:'Escola';not null"`
		Country   string `gorm:"not null"`
		Province  string `gorm:"not null"`
		City      string `gorm:"not null"`
		Phone     string `gorm:"not null;unique"`
		Email     string `gorm:"not null;unique"`
		Website   string
		Facebook  string
		Twitter   string
		Linkedin  string
		Address   string   `gorm:"not null;default:''"`
		Users     []User   `json:",omitempty" gorm:"foreignkey:Organization;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
		Courses   []Course `json:",omitempty" gorm:"foreignkey:School;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
		Documents []File   `json:",omitempty" gorm:"foreignkey:ResourceID;association_foreignkey:uuid;PRELOAD:false;association_autoupdate:false"`
		Logo      string
	}
)

func (a Organization) AllowedStatuses(newStatus string) bool {
	allowed, ok := orgStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}
