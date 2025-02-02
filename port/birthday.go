package port

import (
	"roman/util"
)

type Birthday struct {
	DiscordId uint
	Date      uint
	BirthYear uint
}

type BirthdayService interface {
	SetBd(discordId uint, date uint, year uint) util.RomanError
	GetBirthdaysFromDate(date uint, maxAmount uint) ([]Birthday, util.RomanError)
	GetBirthdaysOnDate(date uint) ([]Birthday, util.RomanError)
}

type BirthdayRepo interface {
	SetBd(discordId uint, date uint, year uint) util.RomanError
	GetBirthdaysFromDate(date uint, maxAmount uint) ([]Birthday, util.RomanError)
	GetBirthdaysOnDate(date uint) ([]Birthday, util.RomanError)
}
