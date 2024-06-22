package main

import (
	"testing"
)

func TestMinute(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectError bool
	}{
		{value: "*", expected: "Every minute"},
		{value: "0", expected: "At 0 minute"},
		{value: "1", expected: "At 1 minute"},
		{value: "2", expected: "At 2 minute"},
		{value: "3", expected: "At 3 minute"},
		{value: "4", expected: "At 4 minute"},
		{value: "5", expected: "At 5 minute"},
		{value: "10,20", expected: "At 10 minute and 20"},
		{value: "0-3", expected: "At every minute from 0 to 3"},
		{value: "60", expectError: true},
	}
	for _, test := range tests {
		minute := NewMinute(test.value)
		err := minute.Validate()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for value %s, but got none", test.value)
			}
			continue
		}

		result := minute.PrettyFormat()
		if result != test.expected {
			t.Errorf("Unexpected result for value %s: got %s, want %s", test.value, result, test.expected)
		}

	}
}

func TestDayOfWeek(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectError bool
	}{
		{value: "*", expected: "every day of week"},
		{value: "0", expected: "on day of week Sunday"},
		{value: "1", expected: "on day of week Monday"},
		{value: "2", expected: "on day of week Tuesday"},
		{value: "3", expected: "on day of week Wednesday"},
		{value: "4", expected: "on day of week Thursday"},
		{value: "5", expected: "on day of week Friday"},
		{value: "6", expected: "on day of week Saturday"},
		{value: "MON", expected: "on day of week Monday"},
		{value: "TUE", expected: "on day of week Tuesday"},
		{value: "WED", expected: "on day of week Wednesday"},
		{value: "THU", expected: "on day of week Thursday"},
		{value: "FRI", expected: "on day of week Friday"},
		{value: "SAT", expected: "on day of week Saturday"},
		{value: "SUN", expected: "on day of week Sunday"},
		{value: "0,1", expected: "on Sunday and Monday"},
		{value: "MON,TUE", expected: "on Monday and Tuesday"},
		{value: "0-3", expected: "from Sunday to Wednesday"},
		{value: "MON-WED", expected: "from Monday to Wednesday"},
		{value: "10", expectError: true},
	}

	for _, test := range tests {
		dayOfWeek := NewDayOfWeek(test.value)
		err := dayOfWeek.Validate()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for value %s, but got none", test.value)
			}
			continue
		}

		result := dayOfWeek.PrettyFormat()
		if result != test.expected {
			t.Errorf("Unexpected result for value %s: got %s, want %s", test.value, result, test.expected)
		}
	}
}

func TestHour(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectError bool
	}{
		{value: "*", expected: "Every hour"},
		{value: "0", expected: "At 0 hour"},
		{value: "1", expected: "At 1 hour"},
		{value: "2", expected: "At 2 hour"},
		{value: "3", expected: "At 3 hour"},
		{value: "4", expected: "At 4 hour"},
		{value: "5", expected: "At 5 hour"},
		{value: "10,20", expected: "past hour 10 and 20"},
		{value: "0-3", expected: "past every hour from 0 to 3"},
		{value: "24", expectError: true},
	}
	for _, test := range tests {
		hour := NewHour(test.value)
		err := hour.Validate()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for value %s, but got none", test.value)
			}
			continue
		}

		result := hour.PrettyFormat()
		if result != test.expected {
			t.Errorf("Unexpected result for value %s: got %s, want %s", test.value, result, test.expected)
		}
	}
}
func TestDayOfMonth(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectError bool
	}{
		{value: "*", expected: "every day of month"},
		{value: "1", expected: "on day of month 1"},
		{value: "2", expected: "on day of month 2"},
		{value: "3", expected: "on day of month 3"},
		{value: "4", expected: "on day of month 4"},
		{value: "5", expected: "on day of month 5"},
		{value: "10,20", expected: "on day of month 10 and 20"},
		{value: "1-5", expected: "on every day of month from 1 to 5"},
		{value: "31", expected: "on day of month 31"},
		{value: "32", expectError: true},
	}
	for _, test := range tests {
		dayOfMonth := NewDayOfMonth(test.value)
		err := dayOfMonth.Validate()
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for value %s, but got none", test.value)
			}
			continue
		}
		result := dayOfMonth.PrettyFormat()
		if result != test.expected {
			t.Errorf("Unexpected result for value %s: got %s, want %s", test.value, result, test.expected)
		}
	}
}
func TestMonth(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectError bool
	}{
		{value: "*", expected: "every month"},
		{value: "1", expected: "in January"},
		{value: "2", expected: "in February"},
		{value: "3", expected: "in March"},
		{value: "4", expected: "in April"},
		{value: "5", expected: "in May"},
		{value: "6", expected: "in June"},
		{value: "7", expected: "in July"},
		{value: "8", expected: "in August"},
		{value: "9", expected: "in September"},
		{value: "10", expected: "in October"},
		{value: "11", expected: "in November"},
		{value: "12", expected: "in December"},
		{value: "JAN", expected: "in January"},
		{value: "FEB", expected: "in February"},
		{value: "MAR", expected: "in March"},
		{value: "APR", expected: "in April"},
		{value: "MAY", expected: "in May"},
		{value: "JUN", expected: "in June"},
		{value: "JUL", expected: "in July"},
		{value: "AUG", expected: "in August"},
		{value: "SEP", expected: "in September"},
		{value: "OCT", expected: "in October"},
		{value: "NOV", expected: "in November"},
		{value: "DEC", expected: "in December"},
		{value: "1,2", expected: "in January and February"},
		{value: "JAN,FEB", expected: "in January and February"},
		{value: "1-3", expected: "from January to March"},
		{value: "JAN-MAR", expected: "from January to March"},
		{value: "13", expectError: true},
	}
	for _, test := range tests {
		month := NewMonth(test.value)
		err := month.Validate()
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for value %s, but got none", test.value)
			}
			continue
		}
		result := month.PrettyFormat()
		if result != test.expected {
			t.Errorf("Unexpected result for value %s: got %s, want %s", test.value, result, test.expected)
		}
	}
}
