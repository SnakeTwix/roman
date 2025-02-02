package model

import "time"

type Birthday struct {
	DiscordId      uint `gorm:"primaryKey"`
	UpdatedAt      time.Time
	BirthdayDate   uint `gorm:"not null"`
	BirthdayYear   *uint
	ShouldBeListed bool
	ShouldBePinged bool
}
