package timechecks

import (
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/go-errlib"
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
	return &TimeChecks{config: cfg}
}

func (t *TimeChecks) IsNeedForUpdateDb(updateDatetime *models.UpdateDatetime) (bool, error) {
	latestUpdateDatetime, err := time.Parse(
		time.RFC3339,
		updateDatetime.UpdateDatetime,
	)
	if err != nil {
		return false, errlib.Wrap(err, "could not parse update time from db")
	}

	todayUpdateDatetime, err := t.DayUpdateDatetime(dayToday)
	if err != nil {
		return false, errlib.Wrap(err, "could not get today update datetime")
	}

	yesterdayUpdateDatetime, err := t.DayUpdateDatetime(dayYesterday)
	if err != nil {
		return false, errlib.Wrap(err, "could not get yesterday update datetime")
	}

	currentDatetime := time.Now()

	isNeedUpdate := !(latestUpdateDatetime.After(todayUpdateDatetime) ||
		((latestUpdateDatetime.After(yesterdayUpdateDatetime) &&
			latestUpdateDatetime.Before(todayUpdateDatetime)) &&
			currentDatetime.Before(todayUpdateDatetime)))

	return isNeedUpdate, nil
}

func (t *TimeChecks) TimeToNextUpdate() (time.Duration, error) {
	currentDatetime := time.Now()
	day := dayToday

	var (
		todayUpdateDatetime time.Time
		nextUpdateDatetime  time.Time
		timeToNextUpdate    time.Duration
		err                 error
	)

	todayUpdateDatetime, err = t.DayUpdateDatetime(dayToday)
	if err != nil {
		return timeToNextUpdate, errlib.Wrap(err, "could not get today update datetime")
	}

	if currentDatetime.After(todayUpdateDatetime) {
		day = dayTomorrow
	}

	if nextUpdateDatetime, err = t.DayUpdateDatetime(day); err != nil {
		return timeToNextUpdate, errlib.Wrap(err, "could not get next update datetime")
	}

	timeToNextUpdate = time.Since(nextUpdateDatetime).Abs()

	return timeToNextUpdate, nil
}

func (t *TimeChecks) DayUpdateDatetime(todayOffset int) (time.Time, error) {
	updateTime, err := time.Parse(
		time.TimeOnly,
		t.config.TimeWhenNeedToUpdateCurrency,
	)
	if err != nil {
		return updateTime, errlib.Wrap(err, "could not parse update time from config")
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

	return todayUpdateDatetime, nil
}
