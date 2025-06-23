package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	calendar "github.com/Fabianexe/bubble-calendar"
)

type Event struct {
	start time.Time
	end   time.Time
	title string
}

func (e Event) Start() time.Time {
	return e.start
}

func (e Event) End() time.Time {
	return e.end
}

func (e Event) Title() string {
	return e.title
}

var _ calendar.Event = Event{}

type Model struct {
	Calendar *calendar.Calendar
}

func (m *Model) Init() tea.Cmd {
	return m.Calendar.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Calendar, cmd = m.Calendar.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return m.Calendar.View()
}

func main() {
	events := []calendar.Event{
		Event{
			start: time.Now().Add(-73 * time.Hour),
			end:   time.Now().Add(-25 * time.Hour),
			title: "Two Days",
		},
		Event{
			start: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Add(-30 * time.Minute),
			end:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Add(30 * time.Minute),
			title: "In the middle of the night",
		},
		Event{
			start: time.Now().Add(-2 * time.Hour),
			end:   time.Now().Add(-time.Hour),
			title: "Surrounding Hour",
		},
		Event{
			start: time.Now().Add(-time.Hour),
			end:   time.Now(),
			title: "Last Hour",
		},
		Event{
			start: time.Now(),
			end:   time.Now().Add(time.Hour),
			title: "Now",
		},
		Event{
			start: time.Now().Add(time.Hour),
			end:   time.Now().Add(2 * time.Hour),
			title: "Next Hour",
		},
		Event{
			start: time.Now().Add(2 * time.Hour),
			end:   time.Now().Add(3 * time.Hour),
			title: "Surrounding Hour",
		},
	}

	m := &Model{
		Calendar: calendar.New("Example Calendar", 25, 10, events...),
	}

	m.Calendar.ExternalKeybinds = append(m.Calendar.ExternalKeybinds, key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp(" ctrl+c/q", "Quit"),
	))

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
