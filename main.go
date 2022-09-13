package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// bubble models
	stopwatch stopwatch.Model
	textinput textinput.Model
	keymap    keymap
	help      help.Model

	quitting bool
	// target sentence for player to type
	ReferenceSentence string
	// letters player has typed
	Typed string
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

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
	return m.stopwatch.Init()
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

	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := m.stopwatch.View() + "\n"
	if !m.quitting {
		// render timer
		s = "Elapsed: " + s

		// render reference sentence
		s += "\n" + m.ReferenceSentence + "\n"

		m.textinput.View()

		// render help
		s += m.helpView()
	}
	return s
}

// run program
func main() {
	// init and define textinput
	ti := textinput.New()

	m := model{
		// init stopwatch
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
		// set textinput field
		textinput: ti,
		// init and define keymap
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
		ReferenceSentence: "jupiter coffee tiddleywinks clock funky helpless",
	}

	m.keymap.start.SetEnabled(false)

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
