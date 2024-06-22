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
				response: "Every day of month",
			},
			"number": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "On day of month %s",
			},
			"separator": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1]),([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "On day of month %s and %s",
			},
			"range": {
				regexp:   `^([1-9]|[1-2][0-9]|[3][0-1])-([1-9]|[1-2][0-9]|[3][0-1])$`,
				response: "On every day of month from %s to %s",
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
}

func NewMonth(value string) *Month {
	month := &Month{
		value: value,
		patterns: map[string]pattern{
			"*": {
				regexp:   "*",
				response: "Every month",
			},
			"number": {
				regexp:   `^([1-9]|1[0-2])$`,
				response: "in %s",
			},
			"separator": {
				regexp:   `^([1-9]|1[0-2]),([1-9]|1[0-2])$`,
				response: "in %s and %s",
			},
			"range": {
				regexp:   `^([1-9]|1[0-2])-([1-9]|1[0-2])$`,
				response: "from %s to %s",
			},
		},
	}
	return month
}

func (m *Month) Validate() error {
	resp, err := validate(m.patterns, m.value)
	if err != nil {
		return err
	}
	m.response = resp
	return nil
}

func (m *Month) PrettyFormat() string {
	months := []string{
		"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December",
	}

	for x := 1; x <= 12; x++ {
		pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(fmt.Sprintf("%d", x)))
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		if re.MatchString(m.response) {
			m.response = re.ReplaceAllString(m.response, months[x-1])
		}
	}

	return m.response
}

func main() {

	// param := os.Args[1]
	// fmt.Println(param)

	cron := []string{"5", "10,20", "31", "1-12"}

	minute := NewMinute(cron[0])
	hour := NewHour(cron[1])
	dayOfMonth := NewDayOfMonth(cron[2])
	month := NewMonth(cron[3])

	crontab := []Cron{minute, hour, dayOfMonth, month}

	for _, c := range crontab {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		fmt.Println(c.PrettyFormat())
	}

}
