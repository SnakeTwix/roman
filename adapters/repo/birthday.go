package repo

import (
	"gorm.io/gorm"
	"roman/adapters/repo/model"
	"roman/port"
	"roman/util"
)

type BirthdayRepo struct {
	db *gorm.DB
}

func NewBirthdayRepo(db *gorm.DB) *BirthdayRepo {
	return &BirthdayRepo{db: db}
}

func (b *BirthdayRepo) SetBd(discordId uint, date uint, year uint) util.RomanError {
	bd := model.Birthday{
		DiscordId:    discordId,
		BirthdayDate: date,
		BirthdayYear: &year,

		// The defaults are this, TODO: probably should let users decide these
		ShouldBeListed: true,
		ShouldBePinged: true,
	}

	if year == 0 {
		bd.BirthdayYear = nil
	}

	return util.NewErrorWithDisplay("[BirthdayRepo SetBd]", b.db.Save(&bd).Error, "Failed to write birthday to db")
}

func (b *BirthdayRepo) GetBirthdaysFromDate(date uint, maxAmount uint) ([]port.Birthday, util.RomanError) {
	var amountOfBirthdays int64
	err := b.db.Model(&model.Birthday{}).Count(&amountOfBirthdays).Error
	if err != nil {
		return nil, util.NewErrorWithDisplay("[BirthdayRepo GetBirthdayFromDate]", err, "Error when counting amount of birthdays")
	}

	if uint(amountOfBirthdays) < maxAmount {
		maxAmount = uint(amountOfBirthdays)
	}

	portBirthdays := make([]port.Birthday, 0, maxAmount)
	modelBirthdays := make([]model.Birthday, 0)

	err = b.db.Limit(int(maxAmount)).Order("birthday_date ASC, discord_id ASC").Where("birthday_date >= ?", date).Find(&modelBirthdays).Error
	if err != nil {
		return nil, util.NewErrorWithDisplay("[BirthdayRepo GetBirthdayFromDate]", err, "Error when retrieving birthdays")
	}

	for _, birthday := range modelBirthdays {
		portBirthday := port.Birthday{
			DiscordId: birthday.DiscordId,
			Date:      birthday.BirthdayDate,
		}

		if birthday.BirthdayYear != nil {
			portBirthday.BirthYear = *birthday.BirthdayYear
		}

		portBirthdays = append(portBirthdays, portBirthday)
	}

	if !(len(modelBirthdays) < int(maxAmount)) {
		return portBirthdays, nil
	}

	modelBirthdays = modelBirthdays[:0]
	err = b.db.Limit(int(maxAmount)-len(modelBirthdays)).Order("birthday_date ASC, discord_id ASC").Where("birthday_date < ?", date).Find(&modelBirthdays).Error
	if err != nil {
		return nil, util.NewErrorWithDisplay("[BirthdayRepo GetBirthdayFromDate]", err, "Error when retrieving birthdays")
	}

	for _, birthday := range modelBirthdays {
		portBirthday := port.Birthday{
			DiscordId: birthday.DiscordId,
			Date:      birthday.BirthdayDate,
		}

		if birthday.BirthdayYear != nil {
			portBirthday.BirthYear = *birthday.BirthdayYear
		}

		portBirthdays = append(portBirthdays, portBirthday)
	}

	return portBirthdays, nil
}

func (b *BirthdayRepo) GetBirthdaysOnDate(date uint) ([]port.Birthday, util.RomanError) {
	portBirthdays := make([]port.Birthday, 0)
	modelBirthdays := make([]model.Birthday, 0)

	err := b.db.Order("discord_id ASC").Where("birthday_date = ?", date).Find(&modelBirthdays).Error
	if err != nil {
		return nil, util.NewErrorWithDisplay("[BirthdayRepo GetBirthdayOnDate]", err, "Error when retrieving birthdays")
	}

	for _, birthday := range modelBirthdays {
		portBirthday := port.Birthday{
			DiscordId: birthday.DiscordId,
			Date:      birthday.BirthdayDate,
		}

		if birthday.BirthdayYear != nil {
			portBirthday.BirthYear = *birthday.BirthdayYear
		}

		portBirthdays = append(portBirthdays, portBirthday)
	}

	return portBirthdays, nil
}
