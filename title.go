package calendar

import (
	"fmt"
	"time"
)

type title struct {
	name     string
	color    int
	duration time.Duration
}

type titleMap map[string]*title

func (t titleMap) get(name string) (*title, bool) {
	ele, ok := t[name]
	if !ok {
		ele = &title{name: name}
		t[name] = ele
	}

	return ele, ok
}

func (t titleMap) list() []*title {
	list := make([]*title, 0, len(t))
	for _, ele := range t {
		list = append(list, ele)
	}

	return list
}

func (t *title) durationString() string {
	d := t.duration.Round(time.Minute)
	if d == 0 {
		return "0m"
	}

	s := ""
	if d.Hours() > 24 {
		s += fmt.Sprintf("%dd", int(d.Hours()/24))
		d -= 24 * time.Hour * time.Duration(int(d.Hours()/24))
	}

	if d.Hours() > 0 {
		s += fmt.Sprintf("%dh", int(d.Hours()))
		d -= time.Hour * time.Duration(int(d.Hours()))
	}

	if d.Minutes() > 0 {
		s += fmt.Sprintf("%dm", int(d.Minutes()))
	}

	return s

}
