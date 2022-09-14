package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"

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
	referenceSentence string
	// tracks which letters were correctly typed
	hits []bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

var hitStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0BF48B"))

var missStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#F12746"))

var blankStyle = lipgloss.NewStyle()

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

		switch msg.Type {
		case tea.KeyRunes, tea.KeySpace:

			typed := m.textinput.Value()
			if len(typed) > len(m.hits) &&
				msg.Runes[0] == rune(m.referenceSentence[len(typed)]) {

				m.hits = append(m.hits, true)
			} else {
				m.hits = append(m.hits, false)
			}
		case tea.KeyBackspace:
			if len(m.hits) > 0 {
				m.hits = m.hits[0 : len(m.hits)-1]
			}
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
		hits := []int{}
		misses := []int{}

		for idx, val := range m.hits {
			if val {
				hits = append(hits, idx)
			} else {
				misses = append(misses, idx)
			}
		}

		refStr := lipgloss.StyleRunes(m.referenceSentence, hits, hitStyle, blankStyle)
		refStr = lipgloss.StyleRunes(refStr, misses, missStyle, blankStyle)

		// render reference sentence
		s += "\n" + refStr + "\n"

		// render text box
		s += m.textinput.View()

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
		hits:              []bool{},
	}

	m.keymap.start.SetEnabled(false)

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
