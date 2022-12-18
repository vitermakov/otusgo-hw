package model

import (
	"time"
)

// DateAlign способ выравнивания даты
type DateAlign int

const (
	DateAlignNone DateAlign = iota
	DateAlignDay
	DateAlignWeek
	DateAlignMonth
	DateAlignError
)

func (da DateAlign) IsValid() bool {
	return da > DateAlignNone && da < DateAlignError
}

// DateRange промежуток дат.
type DateRange struct {
	DateStart time.Time
	Duration  time.Duration
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
func DateRgnOnDay(date time.Time) DateRange {
	return DateRange{
		DateStart: alignDate(date, DateAlignDay),
		Duration:  time.Hour * 24,
	}
}
func DateRgnOnWeek(date time.Time) DateRange {
	return DateRange{
		DateStart: alignDate(date, DateAlignWeek),
		Duration:  time.Hour * 24 * 7,
	}
}
func DateRgnOnMonth(date time.Time) DateRange {
	date = alignDate(date, DateAlignMonth)
	return DateRgnFromDates(date, date.AddDate(0, 1, 0))
}
func alignDate(date time.Time, align DateAlign) time.Time {
	switch align {
	case DateAlignDay:
	case DateAlignWeek:
	case DateAlignMonth:
	}
	return date
}
