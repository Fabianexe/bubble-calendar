package calendar

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode uint8

const (
	DayMode mode = iota
	WeekMode
)

const (
	keyWeek = iota
	keyDay
	keyPrev
	keyNext
)

var (
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)
)

type Event interface {
	Start() time.Time
	End() time.Time
	Title() string
}

type Calendar struct {
	CurrentMode      mode
	AllowedModes     []mode
	Events           []Event
	Title            string
	ViewedDay        time.Time
	Colors           []lipgloss.Style
	ExternalKeybinds []key.Binding
	WeekStartMonday  bool

	internalKeybinds []key.Binding
	width            int
	height           int
	help             help.Model
}

func New(title string, width, height int, events ...Event) *Calendar {
	return &Calendar{
		CurrentMode:  DayMode,
		AllowedModes: []mode{DayMode, WeekMode},
		ViewedDay:    time.Now(),
		Colors: []lipgloss.Style{
			lipgloss.NewStyle().Background(lipgloss.Color("#ff0000")).Foreground(lipgloss.Color("#000000")),
			lipgloss.NewStyle().Background(lipgloss.Color("#00ff00")).Foreground(lipgloss.Color("#000000")),
			lipgloss.NewStyle().Background(lipgloss.Color("#0000ff")).Foreground(lipgloss.Color("#ffffff")),
			lipgloss.NewStyle().Background(lipgloss.Color("#ffff00")).Foreground(lipgloss.Color("#000000")),
			lipgloss.NewStyle().Background(lipgloss.Color("#ff00ff")).Foreground(lipgloss.Color("#000000")),
			lipgloss.NewStyle().Background(lipgloss.Color("#00ffff")).Foreground(lipgloss.Color("#000000")),
		},
		Events: events,
		Title:  title,
		width:  width,
		height: height,
		internalKeybinds: []key.Binding{
			keyWeek: key.NewBinding(
				key.WithKeys("up", "w"),
				key.WithHelp("↑/w", "Week View"),
			),
			keyDay: key.NewBinding(
				key.WithKeys("down", "d"),
				key.WithHelp("↓/d", "Day View"),
			),
			keyPrev: key.NewBinding(
				key.WithKeys("left", "a"),
				key.WithHelp("←/a", "Previous Day"),
			),
			keyNext: key.NewBinding(
				key.WithKeys("right", "s"),
				key.WithHelp("→/s", "Next Day"),
			),
		},
		help:            help.New(),
		WeekStartMonday: true,
	}
}

func (c *Calendar) Init() tea.Cmd {
	return nil
}

func (c *Calendar) Update(msg tea.Msg) (*Calendar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.internalKeybinds[keyWeek]):
			c.CurrentMode = WeekMode
		case key.Matches(msg, c.internalKeybinds[keyDay]):
			c.CurrentMode = DayMode
		case key.Matches(msg, c.internalKeybinds[keyPrev]):
			c.ViewedDay = c.ViewedDay.AddDate(0, 0, -1)
		case key.Matches(msg, c.internalKeybinds[keyNext]):
			c.ViewedDay = c.ViewedDay.AddDate(0, 0, 1)
		}

	case tea.WindowSizeMsg:
		c.width, c.height = msg.Width, msg.Height
	}

	return c, nil
}

func (c *Calendar) View() string {
	var sections []string
	sections = append(sections, titleStyle.Render(c.Title))

	switch c.CurrentMode {
	case DayMode:
		sections = append(sections, c.dayView())
	case WeekMode:
		sections = append(sections, c.weekView())
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	helpView := c.help.ShortHelpView(append(c.internalKeybinds, c.ExternalKeybinds...))
	height := c.height - strings.Count(content, "\n") - strings.Count(helpView, "\n") - 1

	if height <= 0 {
		return content
	}

	return content + strings.Repeat("\n", height) + helpView
}

func (c *Calendar) filterEvents(start, end time.Time) []Event {
	filtered := make([]Event, 0, len(c.Events))
	for _, e := range c.Events {
		startOverlap := start.Before(e.Start()) && end.After(e.Start())
		endOverlap := start.Before(e.End()) && end.After(e.End())
		if startOverlap || endOverlap {
			filtered = append(filtered, e)
		}
	}

	return filtered
}
