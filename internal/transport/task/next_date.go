package task

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	dateParam   = "date"
	nowParam    = "now"
	repeatParam = "repeat"
)

var weekDay = map[int]string{
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
	7: "Sunday",
}

var month = map[int]string{
	1:  "January",
	2:  "February",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

func (h *Handlers) CalculateNextDate(c echo.Context) error {
	date, err := time.Parse("20060102", c.FormValue(dateParam))
	if err != nil {
		slog.Info("failed to parse date", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "date format must be: 20060102"})
	}

	now, err := time.Parse("20060102", c.FormValue(nowParam))
	if err != nil {
		slog.Info("failed to parse date", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "now format must be: 20060102"})
	}

	repeat := c.FormValue(repeatParam)
	if !repeatFormat(repeat) {
		slog.Info("repeat format is wrong")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong repeat format"})
	}

	nextDate := calcNextDate(now, date, repeat)

	return c.String(http.StatusOK, nextDate.Format("20060102"))
}

func calcNextDate(now, date time.Time, repeat string) time.Time {
	var next time.Time
	switch strings.Split(repeat, " ")[0] {
	case "d":
		days, _ := strconv.Atoi(strings.Split(repeat, " ")[1])
		next = date.AddDate(0, 0, days)
		for next.Before(now) {
			next = next.AddDate(0, 0, days)
		}
	case "y":
		next = date.AddDate(1, 0, 0)
		for next.Before(now) {
			next = next.AddDate(1, 0, 0)
		}
	case "w":
		next = calcWeekNextDate(repeat, now)
	case "m":
		if len(strings.Split(repeat, " ")) == 2 {
			next = calcMonthNextDate(repeat, now, date)
		}
		if len(strings.Split(repeat, " ")) == 3 {
			next = calcMonthDayNextDate(repeat, now, date)
		}
	}

	return next
}

func calcWeekNextDate(repeat string, now time.Time) time.Time {
	var nextDate time.Time
	daysOfWeek := strings.Split(repeat, " ")[1]
	days := strings.Split(daysOfWeek, ",")
	dayNumbers := make([]int, 0, len(days))
	for _, day := range days {
		dayNumber, _ := strconv.Atoi(day)
		dayNumbers = append(dayNumbers, dayNumber)
	}
	done := false
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		for _, dayNumber := range dayNumbers {
			if nextDate.Weekday().String() == weekDay[dayNumber] && nextDate.After(now) {
				done = true
				break
			}
		}
		if done {
			return nextDate
		}
	}
}

func calcMonthNextDate(repeat string, now, date time.Time) time.Time {
	var nextDate time.Time
	days := strings.Split(strings.Split(repeat, " ")[1], ",")
	dayNumbers := make([]int, 0, len(days))
	for _, day := range days {
		dayNumber, _ := strconv.Atoi(day)
		dayNumbers = append(dayNumbers, dayNumber)
	}

	done := false
	nextDate = date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		for _, dayNumber := range dayNumbers {
			if dayNumber < 0 {
				dayNumber = getDayOfMonth(nextDate, dayNumber)
			}
			if nextDate.Day() == dayNumber && nextDate.After(now) {
				done = true
				break
			}
		}
		if done {
			return nextDate
		}
	}
}

func calcMonthDayNextDate(repeat string, now, date time.Time) time.Time {
	var nextDate time.Time
	days := strings.Split(strings.Split(repeat, " ")[1], ",")
	dayNumbers := make([]int, 0, len(days))
	for _, day := range days {
		dayNumber, _ := strconv.Atoi(day)
		dayNumbers = append(dayNumbers, dayNumber)
	}

	months := strings.Split(strings.Split(repeat, " ")[2], ",")
	monthsNumbers := make([]int, 0, len(months))
	for _, m := range months {
		monthNumber, _ := strconv.Atoi(m)
		monthsNumbers = append(monthsNumbers, monthNumber)
	}

	monthMap := make(map[int]string)
	for _, m := range monthsNumbers {
		monthMap[m] = month[m]
	}

	done := false
	nextDate = date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		for _, dayNumber := range dayNumbers {
			if dayNumber < 0 {
				dayNumber = getDayOfMonth(nextDate, dayNumber)
			}
			m := int(nextDate.Month())
			var mOk bool
			_, mOk = monthMap[m]
			if nextDate.Day() == dayNumber && nextDate.After(now) && mOk {
				done = true
				break
			}
		}
		if done {
			return nextDate
		}
	}
}

func getDayOfMonth(nextDate time.Time, day int) int {
	date := time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	date = date.AddDate(0, 0, day)
	return date.Day()
}

func repeatFormat(repeat string) bool {
	switch strings.Split(repeat, " ")[0] {
	case "d":
		if len(strings.Split(repeat, " ")) != 2 {
			return false
		}
		digit, err := strconv.Atoi(strings.Split(repeat, " ")[1])
		if err != nil {
			return false
		}
		if digit < 1 || digit > 400 {
			return false
		}
	case "y":
		if len(strings.Split(repeat, " ")) != 1 {
			return false
		}
	case "w":
		if len(strings.Split(repeat, " ")) != 2 {
			return false
		}
		ds := strings.Split(repeat, " ")[1]
		digits := strings.Split(ds, ",")
		for _, d := range digits {
			digit, err := strconv.Atoi(d)
			if err != nil {
				return false
			}
			if digit < 1 || digit > 7 {
				return false
			}
		}
	case "m":
		if len(strings.Split(repeat, " ")) > 3 {
			return false
		}

		if len(strings.Split(repeat, " ")) == 2 {
			digits := strings.Split(strings.Split(repeat, " ")[1], ",")
			for _, d := range digits {
				digit, err := strconv.Atoi(d)
				if err != nil {
					return false
				}
				if digit < -2 || digit > 31 {
					return false
				}
			}
		}

		if len(strings.Split(repeat, " ")) == 3 {
			days := strings.Split(strings.Split(repeat, " ")[1], ",")
			for _, d := range days {
				digit, err := strconv.Atoi(d)
				if err != nil {
					return false
				}
				if digit < -2 || digit > 31 {
					return false
				}
			}
			months := strings.Split(strings.Split(repeat, " ")[2], ",")
			for _, m := range months {
				digit, err := strconv.Atoi(m)
				if err != nil {
					return false
				}
				if digit < 1 || digit > 12 {
					return false
				}
			}
		}
	default:
		return false
	}

	return true
}
