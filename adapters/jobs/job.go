package jobs

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"roman/port"
	"roman/util"
)

// TODO: Add a list of registered jobs somewhere, so they can be tracked

type Job interface {
	Register() util.RomanError
	Execute() util.RomanError
}

type Jobs struct {
	jobs []Job
}

func Init(configService port.ConfigService, scheduler gocron.Scheduler, client bot.Client, birthdayRepo port.BirthdayRepo) *Jobs {
	return &Jobs{
		jobs: []Job{
			NewJobHappyBirthdayGreeting(configService, scheduler, client, birthdayRepo),
		},
	}
}

func (j *Jobs) GetAll() []Job {
	return j.jobs
}
