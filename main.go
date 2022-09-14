package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

// Model
type model struct {
	// bubble models
	stopwatch stopwatch.Model
	textinput textinput.Model
	keymap    keymap
	help      help.Model

	quitting bool
	logger   log.Logger
	// target sentence for player to type
	referenceSentence string
}

// Key bindings
type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

// Styles
var hitStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0BF48B"))

var missStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#F12746")).
	ColorWhitespace(false)

// render help page
func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

// first method to run when model is created
func (m model) Init() tea.Cmd {
	// init bubble models
	batch := tea.Batch(
		m.stopwatch.Init(),
		textinput.Blink,
	)
	return batch
}

// update model's state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// if the message is a keystroke
	case tea.KeyMsg:
		switch {
		// quit the game
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit

		// reset the timer
		case key.Matches(msg, m.keymap.reset):
			return m, m.stopwatch.Reset()

		// toggle timer status
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()
		}
		// TODO update reference sentence style based on key input
		// i.e. one color for a correct input, another for an incorrect input
	}

	// update bubbles
	var stopwatchcmd tea.Cmd
	m.stopwatch, stopwatchcmd = m.stopwatch.Update(msg)
	var textinputcmd tea.Cmd
	m.textinput, textinputcmd = m.textinput.Update(msg)

	batch := tea.Batch(
		stopwatchcmd,
		textinputcmd,
	)
	return m, batch
}

func (m model) View() string {
	s := m.stopwatch.View() + "\n"
	if !m.quitting {
		// render timer
		s = "Elapsed: " + s

		// render reference sentence
		s += "\n" + m.referenceSentence + "\n"

		// render typed sentence
		typed := []rune(m.textinput.Value())
		ref := []rune(m.referenceSentence)
		var styled string

		for idx, r := range typed {
			if r == ref[idx] {
				styled += hitStyle.Render(string(r))
			} else {
				styled += missStyle.Render(string(r))
			}
		}

		s += styled

		// render help
		s += m.helpView()
	}
	return s
}

// run program
func main() {
	// create and define textinput
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	// ti.CharLimit = 156
	// ti.Width = 20

	f, _ := os.OpenFile("TypeTester.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	errLogger := log.New(f, "", log.Lshortfile)

	m := model{
		// create stopwatch
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
		// set textinput field
		textinput: ti,
		// create and define keymap
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("ctrl+r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		// init help
		help: help.NewModel(),
		// TODO replace with function to generate random sentence
		referenceSentence: "jupiter coffee tiddleywinks clock funky helpless",
		logger:            *errLogger,
	}

	m.keymap.start.SetEnabled(false)

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
