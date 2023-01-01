package model

import (
	"errors"
	"time"
)

var ErrorUnknownRangType = errors.New("unknown log range type")

// RangeKind тип промежутка даты.
type RangeKind int

const (
	RangeTypeNone RangeKind = iota
	RangeTypeDay
	RangeTypeWeek
	RangeTypeMonth
	RangeTypeError
)

func (rk RangeKind) Valid() bool {
	return rk > RangeTypeNone && rk < RangeTypeError
}

func ParseRangeType(sType string) (RangeKind, error) {
	switch sType {
	case "day":
		return RangeTypeDay, nil
	case "week":
		return RangeTypeWeek, nil
	case "month":
		return RangeTypeMonth, nil
	}
	return RangeTypeNone, ErrorUnknownRangType
}

// DateRange промежуток дат.
type DateRange struct {
	DateStart time.Time
	Duration  time.Duration
}

func (dr DateRange) Valid() bool {
	return !dr.DateStart.IsZero() && dr.Duration.Abs() > 0
}

func (dr DateRange) GetFrom() time.Time {
	return dr.DateStart
}

func (dr DateRange) GetTo() time.Time {
	return dr.DateStart.Add(dr.Duration)
}

func DateRgnFromDates(from time.Time, to time.Time) DateRange {
	if from.After(to) {
		from, to = to, from
	}
	return DateRange{
		DateStart: from,
		Duration:  to.Sub(from),
	}
}

func DateRgnOn(kind RangeKind, date time.Time) DateRange {
	alignDate(kind, &date)
	switch kind { //nolint:exhaustive // по дефолту DateRange{}
	case RangeTypeDay:
		return DateRange{DateStart: date, Duration: time.Hour * 24}
	case RangeTypeWeek:
		return DateRange{DateStart: date, Duration: time.Hour * 24 * 7}
	case RangeTypeMonth:
		return DateRgnFromDates(date, date.AddDate(0, 1, 0))
	}
	return DateRange{}
}

// alignDate выравнивание даты по началу дня, недели (по понедельнику), месяцу (первому числу).
func alignDate(kind RangeKind, date *time.Time) {
	y, m, d, l := date.Year(), date.Month(), date.Day(), date.Location()
	switch kind { //nolint:exhaustive // другие варианты ошибочны
	case RangeTypeWeek:
		wd := int(date.Weekday()) - 1
		if wd < 0 {
			wd = 6
		}
		*date = time.Date(y, m, d-wd, 0, 0, 0, 0, l)
	case RangeTypeMonth:
		*date = time.Date(y, m, 1, 0, 0, 0, 0, l)
	default:
		*date = time.Date(y, m, d, 0, 0, 0, 0, l)
	}
}
