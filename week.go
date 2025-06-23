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

func (c *Calendar) createWeekBlock(events []Event, hour, day int, start time.Time) ([7]string, []*title) {
	weekBlocks := [7]string{}
	if len(events) == 0 {
		return weekBlocks, nil
	}

	tMap := make(titleMap, len(events))
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

			foundTitle, ok := tMap.get(event.Title())
			if !ok {
				foundTitle.color = colorSpinner
				colorSpinner++
			}

			foundTitle.duration += end.Sub(start)

			contentSize := int(float64(hour) * end.Sub(start).Hours())
			contetntString := fmt.Sprintf("%02d", foundTitle.color)
			topSpace := int(float64(hour) * start.Sub(beginOfDay).Hours())

			content := lipgloss.PlaceHorizontal(day, lipgloss.Center, lipgloss.PlaceVertical(contentSize, lipgloss.Center, contetntString))
			block := c.Colors[foundTitle.color%len(c.Colors)].MarginTop(topSpace - alreadyUsedLine).Render(content)
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

	return weekBlocks, tMap.list()
}

func (c *Calendar) createWeekTitle(titleList []*title) string {
	maxlength := 0
	for _, t := range titleList {
		if lipgloss.Width(t.name) > maxlength {
			maxlength = lipgloss.Width(t.name)
		}
		if c.ShowDuration {
			duration := t.durationString()
			if lipgloss.Width(duration) > maxlength {
				maxlength = lipgloss.Width(duration)
			}
		}
	}

	slices.SortFunc(titleList, func(i, j *title) int {
		if i.color < j.color {
			return -1
		}
		if i.color > j.color {
			return 1
		}

		return 0
	})

	titleBlocks := make([]string, 0, len(titleList))
	titleBlocks = append(titleBlocks, "") // one line free
	for _, t := range titleList {
		color := t.color
		titleBlocks = append(titleBlocks,
			c.Colors[color%len(c.Colors)].MarginLeft(2).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, fmt.Sprintf("%02d", color))),
			c.Colors[color%len(c.Colors)].MarginLeft(2).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, t.name)),
		)
		if c.ShowDuration {
			titleBlocks = append(titleBlocks,
				c.Colors[color%len(c.Colors)].MarginLeft(2).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, t.durationString())),
			)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBlocks...)
}
