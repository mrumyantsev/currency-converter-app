package timechecks

import (
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"
)

const (
	dayYesterday = -1
	dayToday     = 0
	dayTomorrow  = 1
)

type TimeChecks struct {
	config *config.Config
}

func New(cfg *config.Config) *TimeChecks {
	return &TimeChecks{
		config: cfg,
	}
}

func (t *TimeChecks) IsNeedForUpdateDb(updateDatetime *models.UpdateDatetime) (bool, error) {
	latestUpdateDatetime, err := time.Parse(
		time.RFC3339, updateDatetime.UpdateDatetime)
	if err != nil {
		return false, e.Wrap("could not parse update time from db", err)
	}

	todayUpdateDatetime, err := t.GetDayUpdateDatetime(dayToday)
	if err != nil {
		return false, e.Wrap("could not get today update datetime", err)
	}

	yesterdayUpdateDatetime, err := t.GetDayUpdateDatetime(dayYesterday)
	if err != nil {
		return false, e.Wrap("could not get yesterday update datetime", err)
	}

	currentDatetime := time.Now()

	return !(latestUpdateDatetime.After(*todayUpdateDatetime) ||
		((latestUpdateDatetime.After(*yesterdayUpdateDatetime) &&
			latestUpdateDatetime.Before(*todayUpdateDatetime)) &&
			currentDatetime.Before(*todayUpdateDatetime))), nil
}

func (t *TimeChecks) GetTimeToNextUpdate() (*time.Duration, error) {
	var (
		currentDatetime     time.Time = time.Now()
		todayUpdateDatetime *time.Time
		nextUpdateDatetime  *time.Time
		timeToNextUpdate    time.Duration
		day                 int = dayToday
		err                 error
	)

	todayUpdateDatetime, err = t.GetDayUpdateDatetime(dayToday)
	if err != nil {
		return nil, e.Wrap("could not get today update datetime", err)
	}

	if currentDatetime.After(*todayUpdateDatetime) {
		day = dayTomorrow
	}

	nextUpdateDatetime, err = t.GetDayUpdateDatetime(day)
	if err != nil {
		return nil, e.Wrap("could not get next update datetime", err)
	}

	timeToNextUpdate = time.Since(*nextUpdateDatetime).Abs()

	return &timeToNextUpdate, nil
}

func (t *TimeChecks) GetDayUpdateDatetime(todayOffset int) (*time.Time, error) {
	updateTime, err := time.Parse(time.TimeOnly, t.config.TimeWhenNeedToUpdateCurrency)
	if err != nil {
		return nil, e.Wrap("could not parse update time from config", err)
	}

	todayYear, todayMonth, todayDay := time.Now().Date()

	todayUpdateDatetime := time.Date(
		todayYear,
		todayMonth,
		todayDay+todayOffset,
		updateTime.Hour(),
		updateTime.Minute(),
		updateTime.Second(),
		0, // drop nanoseconds
		time.Now().Location(),
	)

	return &todayUpdateDatetime, nil
}
