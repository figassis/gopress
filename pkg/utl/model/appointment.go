package model

import (
	"time"
)

var (
	//'Pendente','Activo','Concluído','Cancelado');default:'Pendente';not null"`
	appointmentStatusFlows = map[string]List{
		StatusPending:              List{StatusActive, StatusCanceled},
		StatusActive:               List{StatusConcludedAppointment, StatusCanceled},
		StatusConcludedAppointment: List{StatusCanceled},
		StatusCanceled:             List{},
	}
)

// User represents user domain model
type (
	Appointment struct {
		Base
		User          string
		UserName      string
		Resource      string `gorm:"type:ENUM('Utilizador','Candidatura');default:'Candidatura';not null"`
		ResourceID    string
		Address       string
		Date          time.Time
		Admin         string
		AdminName     string
		Comments      string
		Status        string `gorm:"type:ENUM('Pendente','Activo','Concluído','Cancelado');default:'Pendente';not null"`
		ContactName   string
		ContactNumber string
		ContactEmail  string
		Message       string
	}
)

func (a Appointment) AllowedStatuses(newStatus string) bool {
	allowed, ok := appointmentStatusFlows[a.Status]
	if !ok {
		return false
	}
	return allowed.Contains(newStatus)
}
