package model

import (
	"gorm.io/datatypes"
)

type Fact struct {
	ID                string         `gorm:"index" json:"id"`
	PeriodStart       string         `gorm:"not null" json:"period_start"`
	PeriodEnd         string         `gorm:"not null" json:"period_end"`
	PeriodKey         string         `gorm:"not null" json:"period_key"`
	IndicatorToMoId   uint32         `gorm:"not null;uint" json:"indicator_to_mo_id"`
	IndicatorToFactId uint32         `gorm:"not null;uint" json:"indicator_to_fact_id"`
	Value             uint8          `gorm:"type:smallint;not null" json:"value"`
	FactTime          string         `gorm:"not null" json:"fact_time"`
	IsPlan            uint8          `gorm:"type:smallint;not null" json:"is_plan"`
	Supertags         datatypes.JSON `gorm:"serializer:json;" json:"supertags"`
	AuthUserId        uint32         `gorm:"not null;uint" json:"auth_user_id"`
	Comment           datatypes.JSON `gorm:"serializer:json;" json:"comment"`
}
