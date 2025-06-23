
# Bubble Calendar
<img alt="simple demo" src="example/simple.gif" width="1200" />
Bubble Calendar is a calendar component for the Bubble Tea framework.
It provides a visual representation of events throughout the day, 
with each event color-coded for easy identification.

## Features

- Visual representation of a day's events.
- Week view feature to see events for the entire week.
- Events are color-coded for easy identification.
- The calendar adjusts to fit the width of your terminal.
- Compatible with the Bubble Tea framework.
- Week start on Monday or Sunday.
- Customizable event colors.

## Installation

You can use the `go get` command:

```bash
go get github.com/Fabianexe/bubble-calendar
```

## Usage

Import the package in your Go file:

```go
import calendar "github.com/Fabianexe/bubble-calendar"
```

Create a new calendar and add events to it:

```go
c := calendar.New("Calendar Title", 25, 10, events...),
```

The events are represented by the `calendar.Event` interface:

```go
type Event interface {
	Start() time.Time
	End() time.Time
	Title() string
}
```

See the [Example code](https://github.com/Fabianexe/bubble-calendar/tree/main/example/simple.go).

## Configuration

The `Calendar` struct has several public fields and methods that you can use to configure your calendar:

### Fields

- `CurrentMode`: The current view mode of the calendar. It can be `calendar.DayMode` or `calendar.WeekMode`.
- `Events`: An array of events to be displayed on the calendar.
- `Title`: The title of the calendar.
- `ViewedDay`: The day currently being viewed in the calendar.
- `Colors`: An array of `lipgloss.Style` that define the colors used for the events in the calendar.
- `ExternalKeybinds`: An array of `key.Binding` that define the external key bindings for the calendar.
- `WeekStartMonday`: A boolean that determines whether the week starts on Monday. If set to `false`, the week starts on Sunday.
- `ShowDuration`: A boolean that determines whether to show the duration of events.

### Methods

- `SetAllowedModes(m ...mode) error`: Sets the allowed modes for the calendar.

Please refer to the source code for more detailed information about these fields and methods.

## Limitations

While Bubble Calendar provides a robust and flexible way to visualize your events, there are a few limitations to be aware of:

- **Space Constraints:**
    Especially in the week view, the space available for each day might be limited depending on the size of your terminal. 
    This could potentially lead to inaccurate event start times being displayed and in the worst case, no calendar is displayed at all.

- **No Parallel Events:**
    Currently, Bubble Calendar does not support the display of parallel events 
    If two events occur at the same time, they will be displayed sequentially. 
    This means that the start time of the second event might not accurately reflect its actual start time.

These limitations may be subject to change in future versions as we continue to improve and expand the capabilities of Bubble Calendar.
## Contributing

Contributions to Bubble Calendar are welcome. Please open an issue or submit a pull request on GitHub.

## License

Bubble Calendar is licensed under the MIT License. See the `LICENSE` file for more details.
