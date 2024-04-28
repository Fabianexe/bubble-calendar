package calendar

import "C"
import (
	"fmt"
	"slices"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (c *Calendar) dayView() string {
	sections := make([]string, 0, 5)
	if c.width <= 52 {
		return "To small to display"
	}
	hour := c.width / 26
	sections = append(sections, lipgloss.PlaceHorizontal(c.width, lipgloss.Center, fmt.Sprintf("Day %s:", c.ViewedDay.Format("2006-01-02"))))

	times := make([]string, 25)
	pipes := make([]string, 25)
	for i := range 25 {
		times = append(times, lipgloss.PlaceHorizontal(hour, lipgloss.Center, fmt.Sprintf("%02d", i)))
		pipes = append(pipes, lipgloss.PlaceHorizontal(hour, lipgloss.Center, "|"))
	}
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left, times...))
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left, pipes...))

	events := c.filterEvents(time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 0, 0, 0, 0, c.ViewedDay.Location()), time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 23, 59, 59, 0, c.ViewedDay.Location()))
	slices.SortFunc(events, func(i, j Event) int {
		if i.Start().Before(j.Start()) {
			return -1
		}
		if i.Start().After(j.Start()) {
			return 1
		}

		return 0
	})
	blocks, title := c.createDayBlock(events, hour)
	sections = append(sections, blocks)

	sections = append(sections, c.createDayTitle(title, hour))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (c *Calendar) createDayBlock(events []Event, hour int) (string, map[string]int) {
	if len(events) == 0 {
		return "", nil
	}
	blocks := make([]string, 0, 100)
	// left padding
	blocks = append(blocks, lipgloss.PlaceHorizontal(hour/2, lipgloss.Center, ""))
	beginOfDay := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 0, 0, 0, 0, c.ViewedDay.Location())
	endOfDay := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 23, 59, 59, 0, c.ViewedDay.Location())
	colorSpinner := 1
	titleMap := make(map[string]int, len(events))
	last := beginOfDay
	alreadyUsedSpaqce := 0
	for _, event := range events {
		start := event.Start()
		if start.Before(last) {
			start = last
		}
		end := event.End()
		if end.After(endOfDay) {
			end = endOfDay
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
		if lipgloss.Width(contetntString) > contentSize {
			contetntString = ""
		}
		leftSpace := int(float64(hour) * start.Sub(beginOfDay).Hours())

		content := lipgloss.PlaceHorizontal(contentSize, lipgloss.Center, contetntString)
		block := c.Colors[currentColor%len(c.Colors)].Copy().PaddingLeft(leftSpace - alreadyUsedSpaqce).ColorWhitespace(false).Render(content)
		alreadyUsedSpaqce += lipgloss.Width(block)
		blocks = append(blocks, block)
		last = event.End()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, blocks...), titleMap
}

func (c *Calendar) createDayTitle(titleMqp map[string]int, hour int) string {
	maxlength := 0
	ordertitle := make([]string, 0, len(titleMqp))
	for title := range titleMqp {
		if lipgloss.Width(title) > maxlength {
			maxlength = lipgloss.Width(title)
		}
		ordertitle = append(ordertitle, title)
	}
	maxlength += hour / 2

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
			c.Colors[color%len(c.Colors)].Copy().PaddingLeft(hour/2).ColorWhitespace(false).Render(fmt.Sprintf("%02d", color))+
				c.Colors[color%len(c.Colors)].Copy().PaddingLeft(hour).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, title)),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBlocks...)
}
