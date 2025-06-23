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

func (c *Calendar) createDayBlock(events []Event, hour int) (string, []*title) {
	if len(events) == 0 {
		return "", nil
	}

	blocks := make([]string, 0, 100)
	// left padding
	blocks = append(blocks, lipgloss.PlaceHorizontal(hour/2, lipgloss.Center, ""))
	beginOfDay := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 0, 0, 0, 0, c.ViewedDay.Location())
	endOfDay := time.Date(c.ViewedDay.Year(), c.ViewedDay.Month(), c.ViewedDay.Day(), 23, 59, 59, 0, c.ViewedDay.Location())
	colorSpinner := 1
	tMap := make(titleMap, len(events))
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

		foundTitle, ok := tMap.get(event.Title())
		if !ok {
			foundTitle.color = colorSpinner
			colorSpinner++
		}

		foundTitle.duration += end.Sub(start)

		contentSize := int(float64(hour) * end.Sub(start).Hours())
		contetntString := fmt.Sprintf("%02d", foundTitle.color)
		if lipgloss.Width(contetntString) > contentSize {
			contetntString = ""
		}
		leftSpace := int(float64(hour) * start.Sub(beginOfDay).Hours())

		content := lipgloss.PlaceHorizontal(contentSize, lipgloss.Center, contetntString)
		block := c.Colors[foundTitle.color%len(c.Colors)].MarginLeft(leftSpace - alreadyUsedSpaqce).Render(content)
		alreadyUsedSpaqce += lipgloss.Width(block)
		blocks = append(blocks, block)
		last = event.End()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, blocks...), tMap.list()
}

func (c *Calendar) createDayTitle(titleList []*title, hour int) string {
	maxlength := 0
	for _, t := range titleList {
		if c.ShowDuration {
			t.name += fmt.Sprintf(" (%s)", t.durationString())
		}

		if lipgloss.Width(t.name) > maxlength {
			maxlength = lipgloss.Width(t.name)
		}
	}
	maxlength += hour / 2

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
	for _, title := range titleList {
		color := title.color
		titleBlocks = append(titleBlocks,
			c.Colors[color%len(c.Colors)].MarginLeft(hour/2).Render(fmt.Sprintf("%02d", color))+
				c.Colors[color%len(c.Colors)].PaddingLeft(hour).Render(lipgloss.PlaceHorizontal(maxlength, lipgloss.Left, title.name)),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBlocks...)
}
