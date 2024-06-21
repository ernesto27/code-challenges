package main

import (
	"fmt"
	"regexp"
)

type Cron interface {
	Validate() error
	PrettyFormat() string
}

type Minute struct {
	value    string
	response string
}

func (m *Minute) Validate() error {
	if m.value == "*" {
		m.response = "Every minute"
		return nil
	}

	pattern := `^([0-9]|[1-5][0-9])$`
	match, err := regexp.MatchString(pattern, m.value)
	if err != nil {
		return err
	}
	if match {
		m.response = fmt.Sprintf("At %s minute", m.value)
		return nil
	}

	pattern = `^([0-9]|[1-5][0-9]),([0-9]|[1-5][0-9])$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(m.value)
	if matches != nil {
		fmt.Println(matches)
		m.response = fmt.Sprintf("At %s minute and %s", matches[1], matches[2])
		return nil
	}

	pattern = `^([0-9]|[1-5][0-9])-([0-9]|[1-5][0-9])$`
	re = regexp.MustCompile(pattern)

	matches = re.FindStringSubmatch(m.value)
	if matches != nil {
		fmt.Println(matches)
		m.response = fmt.Sprintf("At every minute from %s to %s", matches[1], matches[2])
		return nil
	}

	return fmt.Errorf("invalid value for minute: %s", m.value)
}

type Hour struct{}

func (h *Hour) Validate() error {
	return nil
}

func (h *Hour) PrettyFormat() string {
	return ""
}

func (m *Minute) PrettyFormat() string {
	return m.response
}

func main() {

	var minute Cron

	minute = &Minute{value: "13-440"}
	if err := minute.Validate(); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(minute.PrettyFormat())

}
