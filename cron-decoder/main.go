package main

import (
	"fmt"
	"regexp"
)

type Cron interface {
	Validate() error
	PrettyFormat() string
}

type pattern struct {
	regexp   string
	response string
}

func validate(p map[string]pattern, value string) (string, error) {
	var response string
	if value == "*" {
		response = p["*"].response
		return response, nil
	}

	pattern := p["number"].regexp
	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		return "", err
	}
	if match {
		response = fmt.Sprintf(p["number"].response, value)
		return response, nil
	}

	pattern = p["separator"].regexp
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(value)
	if matches != nil {
		response = fmt.Sprintf(p["separator"].response, matches[1], matches[2])
		return response, nil
	}

	pattern = p["range"].regexp
	re = regexp.MustCompile(pattern)

	matches = re.FindStringSubmatch(value)
	if matches != nil {
		response = fmt.Sprintf(p["range"].response, matches[1], matches[2])
		return response, nil
	}

	return "", fmt.Errorf("invalid value: %s", value)
}

type Minute struct {
	value    string
	response string
	patterns map[string]pattern
}

func NewMinute(value string) *Minute {
	minutePatterns := map[string]pattern{
		"*": {
			regexp:   "*",
			response: "Every minute",
		},
		"number": {
			regexp:   `^([0-9]|1[0-9]|2[0-3])$`,
			response: "At %s minute",
		},

		"separator": {
			regexp:   `^([0-9]|[1-5][0-9]),([0-9]|[1-5][0-9])$`,
			response: "At %s minute and %s",
		},
		"range": {
			regexp:   `^([0-9]|[1-5][0-9])-([0-9]|[1-5][0-9])$`,
			response: "At every minute from %s to %s",
		},
	}
	minute := &Minute{
		value:    value,
		patterns: minutePatterns,
	}

	return minute
}

func (m *Minute) Validate() error {
	resp, err := validate(m.patterns, m.value)
	if err != nil {
		return err
	}

	m.response = resp
	return nil
}

type Hour struct {
	value    string
	response string
	patterns map[string]pattern
}

func NewHour(value string) *Hour {
	hour := &Hour{
		value: value,
		patterns: map[string]pattern{
			"*": {
				regexp:   "*",
				response: "Every hour",
			},
			"number": {
				regexp:   `^([0-9]|1[0-9]|2[0-3])$`,
				response: "At %s hour",
			},
			"separator": {
				regexp:   `^([0-9]|[1][0-9]|[2][0-3]),([0-9]|[1][0-9]|[2][0-3])$`,
				response: "past hour %s and %s",
			},
			"range": {
				regexp:   `^([0-9]|[1][0-9]|[2][0-3])-([0-9]|[1][0-9]|[2][0-3])$`,
				response: "past every hour from %s to %s",
			},
		},
	}

	return hour
}

func (h *Hour) Validate() error {
	resp, err := validate(h.patterns, h.value)
	if err != nil {
		return err
	}
	h.response = resp
	return nil
}

func (h *Hour) PrettyFormat() string {
	return h.response
}

func (m *Minute) PrettyFormat() string {
	return m.response
}

type DayOfMonth struct {
	value    string
	response string
	patterns map[string]pattern
}

func NewDayOfMonth(value string) *DayOfMonth {
	month := &DayOfMonth{
		value: value,
		patterns: map[string]pattern{
			"*": {
				regexp:   "*",
				response: "every day of month",
			},
			"number": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "on day of month %s",
			},
			"separator": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1]),([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "on day of month %s and %s",
			},
			"range": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1])-([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "on every day of month from %s to %s",
			},
		},
	}
	return month
}

func (d *DayOfMonth) Validate() error {
	resp, err := validate(d.patterns, d.value)
	if err != nil {
		return err
	}
	d.response = resp
	return nil
}

func (d *DayOfMonth) PrettyFormat() string {
	return d.response
}

type Month struct {
	value    string
	response string
	patterns map[string]pattern
	months   map[string]string
}

func NewMonth(value string) *Month {
	month := &Month{
		value: value,
		patterns: map[string]pattern{
			"*": {
				regexp:   "*",
				response: "every month",
			},
			"number": {
				regexp:   `^([1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)$`,
				response: "in %s",
			},
			"separator": {
				regexp:   `^([1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC),([1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)$`,
				response: "in %s and %s",
			},
			"range": {
				regexp:   `^([1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)-([1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)$`,
				response: "from %s to %s",
			},
		},
		months: map[string]string{
			"1":   "January",
			"2":   "February",
			"3":   "March",
			"4":   "April",
			"5":   "May",
			"6":   "June",
			"7":   "July",
			"8":   "August",
			"9":   "September",
			"10":  "October",
			"11":  "November",
			"12":  "December",
			"JAN": "January",
			"FEB": "February",
			"MAR": "March",
			"APR": "April",
			"MAY": "May",
			"JUN": "June",
			"JUL": "July",
			"AUG": "August",
			"SEP": "September",
			"OCT": "October",
			"NOV": "November",
			"DEC": "December",
		},
	}
	return month
}

func (m *Month) Validate() error {
	resp, err := validate2(m.value, m.patterns, m.months)
	if err != nil {
		return err
	}

	m.response = resp
	return nil
}

func (m *Month) PrettyFormat() string {
	return m.response
}

type DayOfWeek struct {
	value    string
	response string
	patterns map[string]pattern
	days     map[string]string
}

func NewDayOfWeek(value string) *DayOfWeek {
	dayOfWeek := &DayOfWeek{
		value: value,
		patterns: map[string]pattern{
			"*": {
				regexp:   "*",
				response: "every day of week",
			},
			"number": {
				regexp:   `^([0-6]|MON|TUE|WED|THU|FRI|SAT|SUN)$`,
				response: "on day of week %s",
			},
			"separator": {
				regexp:   `^([0-6]|MON|TUE|WED|THU|FRI|SAT|SUN),([0-6]|MON|TUE|WED|THU|FRI|SAT|SUN)$`,
				response: "on %s and %s",
			},
			"range": {
				regexp:   `^([0-6]|MON|TUE|WED|THU|FRI|SAT|SUN)-([0-6]|MON|TUE|WED|THU|FRI|SAT|SUN)$`,
				response: "from %s to %s",
			},
		},
		days: map[string]string{
			"0":   "Sunday",
			"1":   "Monday",
			"2":   "Tuesday",
			"3":   "Wednesday",
			"4":   "Thursday",
			"5":   "Friday",
			"6":   "Saturday",
			"MON": "Monday",
			"TUE": "Tuesday",
			"WED": "Wednesday",
			"THU": "Thursday",
			"FRI": "Friday",
			"SAT": "Saturday",
			"SUN": "Sunday",
		},
	}
	return dayOfWeek
}

func (d *DayOfWeek) Validate() error {
	resp, err := validate2(d.value, d.patterns, d.days)
	if err != nil {
		return err
	}

	d.response = resp
	return nil
}

func (d *DayOfWeek) PrettyFormat() string {
	return d.response
}

func validate2(value string, patterns map[string]pattern, days map[string]string) (string, error) {
	var response string
	if value == "*" {
		response = patterns["*"].response
		return response, nil
	}

	reg := patterns["number"].regexp
	match, err := regexp.MatchString(reg, value)
	if err != nil {
		return "", err
	}
	if match {
		response = fmt.Sprintf(patterns["number"].response, days[value])
		return response, nil
	}

	pattern := patterns["separator"].regexp
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(value)
	if matches != nil {
		response = fmt.Sprintf(patterns["separator"].response, days[matches[1]], days[matches[2]])
		return response, nil
	}

	pattern = patterns["range"].regexp
	re = regexp.MustCompile(pattern)

	matches = re.FindStringSubmatch(value)
	if matches != nil {
		response = fmt.Sprintf(patterns["range"].response, days[matches[1]], days[matches[2]])
		return response, nil
	}

	return "", fmt.Errorf("invalid value: %s", value)

}

func main() {

	// param := os.Args[1]
	// fmt.Println(param)

	cron := []string{"5", "10,20", "31", "JAN,JUL", "0,1"}

	minute := NewMinute(cron[0])
	hour := NewHour(cron[1])
	dayOfMonth := NewDayOfMonth(cron[2])
	month := NewMonth(cron[3])
	dayOfWeek := NewDayOfWeek(cron[4])

	crontab := []Cron{minute, hour, dayOfMonth, month, dayOfWeek}

	for _, c := range crontab {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		fmt.Println(c.PrettyFormat())
	}

}
