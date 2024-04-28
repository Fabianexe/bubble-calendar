package calendar

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (c *Calendar) weekView() string {
	day := (c.width-4)/8 - 1
	aktualeHeight := c.height - 10
	hour := aktualeHeight / 25

	if day < lipgloss.Width("|2006-01-02|") || hour < 1 {
		return "To small to display"
	}

	currentWeekDay := c.ViewedDay.Weekday()
	if c.WeekStartMonday {
		currentWeekDay = (7 + currentWeekDay - 1) % 7
	}
	from := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 0, 0, 0, 0, c.ViewedDay.Location()).AddDate(0, 0, int(currentWeekDay)*-1)
	to := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 23, 59, 59, 0, c.ViewedDay.Location()).AddDate(0, 0, 6-int(currentWeekDay))

	events := c.filterEvents(from, to)

	sections := make([]string, 0, 5)

	sections = append(sections, lipgloss.PlaceHorizontal(c.width, lipgloss.Center, fmt.Sprintf("Week %s to %s", from.Format("2006-01-02"), to.Format("2006-01-02"))))
	verticalParts := make([]string, 0, 18)
	timeStrings := make([]string, 0, 25)
	timeStrings = append(timeStrings, lipgloss.PlaceHorizontal(4, lipgloss.Center, "Time"))
	for i := range 24 {
		timeStrings = append(timeStrings, lipgloss.PlaceVertical(hour, lipgloss.Top, lipgloss.PlaceHorizontal(4, lipgloss.Center, fmt.Sprintf("%02d", i))))
	}
	timeStrings = append(timeStrings, lipgloss.PlaceHorizontal(4, lipgloss.Center, "24"))

	verticalParts = append(verticalParts, lipgloss.JoinVertical(lipgloss.Left, timeStrings...))

	spacer := "|\n" + strings.Repeat("_\n"+strings.Repeat("|\n", hour-1), 24) + "_"
	verticalParts = append(verticalParts, spacer)

	weekBlocks, title := c.createWeekBlock(events, hour, day, from)
	for i := range 7 {
		dayString := lipgloss.PlaceHorizontal(day, lipgloss.Center, from.AddDate(0, 0, i).Format("2006-01-02"))
		if i == int(currentWeekDay) {
			dayString = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Render(dayString)
		}
		verticalParts = append(verticalParts, dayString+"\n\n"+weekBlocks[i])
		verticalParts = append(verticalParts, spacer)
	}
	verticalParts = append(verticalParts, c.createWeekTitle(title))

	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, verticalParts...))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (c *Calendar) createWeekBlock(events []Event, hour, day int, start time.Time) ([7]string, map[string]int) {
	weekBlocks := [7]string{}
	if len(events) == 0 {
		return weekBlocks, nil
	}

	titleMap := make(map[string]int, len(events))
	colorSpinner := 1
	eventPosition := 0
	for i := range 7 {
		blocks := make([]string, 0, 100)
		beginOfDay := start.AddDate(0, 0, i)
		endOfDay := time.Date(beginOfDay.Year(), beginOfDay.Month(), beginOfDay.Day(), 23, 59, 59, 0, beginOfDay.Location())
		breakLoop := false
		last := beginOfDay
		alreadyUsedLine := 0

		for {
			if eventPosition >= len(events) {
				break
			}
			event := events[eventPosition]
			start := event.Start()
			if endOfDay.Before(start) {
				break
			}

			if start.Before(last) {
				start = last
			}
			end := event.End()
			if end.After(endOfDay) {
				end = endOfDay
				// last event of the day
				breakLoop = true
			}

			currentColor := colorSpinner
			if foundColor, ok := titleMap[event.Title()]; !ok {
				titleMap[event.Title()] = colorSpinner
				colorSpinner++
			} else {
				currentColor = foundColor
			}

			contentSize := int(float64(hour) * end.Sub(start).Hours())
			contetntString := fmt.Sprintf("%02d", currentColor)
			topSpace := int(float64(hour) * start.Sub(beginOfDay).Hours())

			content := lipgloss.PlaceHorizontal(day, lipgloss.Center, lipgloss.PlaceVertical(contentSize, lipgloss.Center, contetntString))
			block := c.Colors[currentColor%len(c.Colors)].Copy().PaddingTop(topSpace - alreadyUsedLine).ColorWhitespace(false).Render(content)
			alreadyUsedLine = topSpace + contentSize

			blocks = append(blocks, block)
			if breakLoop {
				break
			}
			last = event.End()
			eventPosition++
		}
		weekBlocks[i] = lipgloss.JoinVertical(lipgloss.Left, blocks...)
	}

	return weekBlocks, titleMap
}

func (c *Calendar) createWeekTitle(titleMqp map[string]int) string {
	maxlength := 0
	ordertitle := make([]string, 0, len(titleMqp))
	for title := range titleMqp {
		if lipgloss.Width(title) > maxlength {
			maxlength = lipgloss.Width(title)
		}
		ordertitle = append(ordertitle, title)
	}

	slices.SortFunc(ordertitle, func(i, j string) int {
		if titleMqp[i] < titleMqp[j] {
			return -1
		}
		if titleMqp[i] > titleMqp[j] {
			return 1
		}

		return 0
	})

	titleBlocks := make([]string, 0, len(titleMqp))
	titleBlocks = append(titleBlocks, "") // one line free
	for _, title := range ordertitle {
		color := titleMqp[title]
		titleBlocks = append(titleBlocks,
			c.Colors[color%len(c.Colors)].Copy().PaddingLeft(2).ColorWhitespace(false).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, fmt.Sprintf("%02d", color))),
			c.Colors[color%len(c.Colors)].Copy().PaddingLeft(2).ColorWhitespace(false).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, title)),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBlocks...)
}
