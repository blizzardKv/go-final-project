package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d.After(n)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}

	parts := strings.SplitN(repeat, " ", 2)
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
			return "", fmt.Errorf("d rule requires interval")
		}
		interval, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("invalid d interval: %s", parts[1])
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
			return "", fmt.Errorf("y rule takes no parameters")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
			return "", fmt.Errorf("w rule requires weekday list")
		}
		var weekdays [8]bool
		for _, ds := range strings.Split(strings.TrimSpace(parts[1]), ",") {
			d, err := strconv.Atoi(strings.TrimSpace(ds))
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("invalid weekday: %s", ds)
			}
			weekdays[d] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] && afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
			return "", fmt.Errorf("m rule requires day list")
		}
		subParts := strings.Fields(parts[1])
		if len(subParts) == 0 || len(subParts) > 2 {
			return "", fmt.Errorf("invalid m rule format")
		}
		var days []int
		for _, ds := range strings.Split(subParts[0], ",") {
			d, err := strconv.Atoi(strings.TrimSpace(ds))
			if err != nil || d == 0 || d < -2 || d > 31 {
				return "", fmt.Errorf("invalid day of month: %s", ds)
			}
			days = append(days, d)
		}
		var months [13]bool
		allMonths := true
		if len(subParts) == 2 {
			allMonths = false
			for _, ms := range strings.Split(subParts[1], ",") {
				m, err := strconv.Atoi(strings.TrimSpace(ms))
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("invalid month: %s", ms)
				}
				months[m] = true
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if !afterNow(date, now) {
				continue
			}
			mon := int(date.Month())
			if !allMonths && !months[mon] {
				continue
			}
			day := date.Day()
			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			for _, d := range days {
				if (d > 0 && d == day) ||
					(d == -1 && day == lastDay) ||
					(d == -2 && day == lastDay-1) {
					return date.Format(DateFormat), nil
				}
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format: %s", repeat)
	}

	return date.Format(DateFormat), nil
}
