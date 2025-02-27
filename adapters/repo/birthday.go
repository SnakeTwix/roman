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

func (b *BirthdayRepo) GetBirthdaysFromDate(date uint, page uint, maxAmount uint) ([]port.Birthday, util.RomanError) {
	amountOfBirthdays, romanErr := b.GetTotalBirthdayCount()
	if romanErr != nil {
		return nil, util.NewError("[BirthdayRepo GetBirthdayFromDate]", romanErr)
	}

	if uint(amountOfBirthdays) < maxAmount {
		maxAmount = uint(amountOfBirthdays)
	}

	portBirthdays := make([]port.Birthday, 0, maxAmount)
	modelBirthdays := make([]model.Birthday, 0)

	err := b.db.Offset(int(maxAmount*page)).Limit(int(maxAmount)).Order("birthday_date ASC, discord_id ASC").Where("birthday_date >= ?", date).Find(&modelBirthdays).Error
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

	firstFetchedAmount := len(modelBirthdays)
	if !(firstFetchedAmount < int(maxAmount)) {
		return portBirthdays, nil
	}

	// Second fetch for start of the year, if we didn't get enough birthdays
	var amountOfBirthdaysAfterDate int64
	err = b.db.Model(&model.Birthday{}).Where("birthday_date >= ?", date).Count(&amountOfBirthdaysAfterDate).Error
	if err != nil {
		return nil, util.NewErrorWithDisplay("[BirthdayRepo GetBirthdayFromDate]", err, "Error when counting amount of birthdays")
	}

	modelBirthdays = modelBirthdays[:0]

	err = b.db.Offset(int(page*maxAmount)-int(amountOfBirthdaysAfterDate)+firstFetchedAmount).Limit(int(maxAmount)-firstFetchedAmount).Order("birthday_date ASC, discord_id ASC").Where("birthday_date < ?", date).Find(&modelBirthdays).Error
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

func (b *BirthdayRepo) GetTotalBirthdayCount() (int, util.RomanError) {
	var totalBirthdayCount int64
	err := b.db.Model(model.Birthday{}).Count(&totalBirthdayCount).Error
	if err != nil {
		return 0, util.NewErrorWithDisplay("[BirthdayRepo GetTotalBirthdayCount]", err, "Error when retrieving birthday count")
	}

	return int(totalBirthdayCount), nil
}
