package service

import (
	"errors"
	"roman/port"
	"roman/util"
	"time"
)

type BirthdayService struct {
	birthdayRepo port.BirthdayRepo
}

func NewBirthdayService(repo port.BirthdayRepo) *BirthdayService {
	return &BirthdayService{birthdayRepo: repo}
}

func (b *BirthdayService) SetBd(discordId uint, date uint, year uint) util.RomanError {
	currentTime := time.Now()
	if year != 0 && currentTime.Year()-int(year) < 13 {
		return util.NewErrorWithDisplay("[BirthdayService SetBd]", errors.New("big age"), "Слишком маленький :baby:")
	}

	if year != 0 && currentTime.Year()-int(year) > 70 {
		return util.NewErrorWithDisplay("[BirthdayService SetBd]", errors.New("big age"), "Слишком долго живешь... :older_adult:")
	}

	return util.NewError("[BirthdayService SetBd]", b.birthdayRepo.SetBd(discordId, date, year))
}

func (b *BirthdayService) GetBirthdaysFromDate(date uint, maxAmount uint) ([]port.Birthday, util.RomanError) {
	birthdays, err := b.birthdayRepo.GetBirthdaysFromDate(date, maxAmount)
	return birthdays, err
}
