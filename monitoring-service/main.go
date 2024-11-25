package main

import (
	"fmt"
	"log"
	"monitoring/db"
	"net/url"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// Update error style with more prominent styling
var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000")).
	Background(lipgloss.Color("#FFE5E5")).
	Padding(0, 1).
	MarginTop(1).
	Bold(true)

// Update success style with more obscure styling
var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#008000")).
	Background(lipgloss.Color("#D0FFD0")).
	Padding(0, 1).
	MarginTop(1).
	Bold(true)

type (
	errMsg error
)

type model struct {
	textInputURL       textinput.Model
	textInputFrequency textinput.Model
	err                error
	validationError    string
	successMessage     string // Add this
	url                string
	frequency          int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter url"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	ti2 := textinput.New()
	ti2.Placeholder = "Enter frequency"
	ti2.CharLimit = 156
	ti2.Width = 20

	return model{
		textInputURL:       ti,
		textInputFrequency: ti2,
		err:                nil,
	}
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInputURL.Focused() {
				if !isValidURL(m.textInputURL.Value()) {
					m.textInputURL.SetValue("")
					m.validationError = "Invalid URL format"
					return m, nil
				}
				m.url = m.textInputURL.Value()
				m.validationError = ""
				m.textInputURL.Blur()
				m.textInputFrequency.Focus()

			} else if m.textInputFrequency.Focused() {
				freqValue, err := strconv.Atoi(m.textInputFrequency.Value())
				if err != nil {
					m.textInputFrequency.SetValue("")
					m.validationError = "Invalid frequency format"
					return m, nil
				}
				m.frequency = freqValue
				m.validationError = ""

				err = myDB.CreateURL(m.url, m.frequency)
				if err != nil {
					m.validationError = "Error adding new URL"
					m.successMessage = ""
					return m, nil
				}
				m.successMessage = "URL successfully saved!"
				m.validationError = ""

				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInputURL, _ = m.textInputURL.Update(msg)
	m.textInputFrequency, cmd = m.textInputFrequency.Update(msg)
	return m, cmd
}

// Update View() function
func (m model) View() string {
	var messageBox string
	if m.validationError != "" {
		messageBox = "\n" + errorStyle.Render(m.validationError)
	} else if m.successMessage != "" {
		messageBox = "\n" + successStyle.Render(m.successMessage)
	}

	return fmt.Sprintf(
		"Enter URL and frequency to monitor\n\n%s%s\n\n%s\n\n%s",
		m.textInputURL.View(),
		messageBox,
		m.textInputFrequency.View(),
		"(esc to quit)",
	) + "\n"
}

var myDB *db.Mysql

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	var dbErr error
	myDB, dbErr = db.NewMysql(host, user, password, port, database)
	if dbErr != nil {
		panic(dbErr)
	}
	// err = myDB.CreateURL("http://www.google.com", 1)
	// if err != nil {
	// 	panic(err)
	// }

	// err = myDB.UpdateURLFrequency("http://www.google.com", 5)
	// if err != nil {
	// 	panic(err)
	// }

	// err = myDB.AddURLHealthCheck(1, 500, 200, 1)
	// if err != nil {
	// 	panic(err)
	// }

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
