package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
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

	// state fields
	referenceSentence string
	quitting          bool
	finished          bool
	wordbank          []string

	// helpers
	logger log.Logger
}

// Key bindings
type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

// Messages
type logMsg struct{ value string }
type startMsg struct{}
type stopMsg struct{}

// Commands
func stopGame() tea.Msg {
	return stopMsg{}
}

func startGame() tea.Msg {
	return startMsg{}
}

func logCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return logMsg{s}
	}
}

// Styles
var hitStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0BF48B"))

var missStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#F12746")).
	ColorWhitespace(false)

var unwrittenStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#828282")).
	ColorWhitespace(false)

// Helper Functions

// Creates a random sentence from a provided list of words
func createNewSentence(length int, words []string) string {
	s := []string{}
	for i := 0; i < length; i++ {
		word := words[rand.Intn(len(words))]
		s = append(s, word)
	}
	return strings.Join(s, " ")
}

// Methods

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
		startGame,
	)
	return batch
}

// update model's state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// log incoming message
	case logMsg:
		m.logger.Println("From message " + msg.value)
		return m, nil
	// if the message is a keystroke
	case tea.KeyMsg:
		switch {
		// quit the game
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit

		// reset the timer
		case key.Matches(msg, m.keymap.reset):
			return m, startGame

		// toggle timer status
		// TODO change this to pause the game
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()

		default:
			// check completion status
			if val := m.textinput.Value() + msg.String(); len(val) > 0 && val == m.referenceSentence {
				m.textinput.Reset()
				m.finished = true
				return m, stopGame
			}
		}
	// Stop the game when it's finished
	case stopMsg:
		return m, m.stopwatch.Stop()
	// Continue playing with a new game
	case startMsg:
		m.textinput.Reset()
		m.referenceSentence = createNewSentence(10, m.wordbank)
		m.textinput.CharLimit = len(m.referenceSentence)
		m.finished = false
		cmd := tea.Batch(
			m.stopwatch.Reset(),
			m.stopwatch.Start(),
			logCmd(fmt.Sprint(m.stopwatch.Interval)),
		)
		return m, cmd
	}

	// update bubbles
	var stopwatchcmd tea.Cmd
	m.stopwatch, stopwatchcmd = m.stopwatch.Update(msg)
	var textinputcmd tea.Cmd
	m.textinput, textinputcmd = m.textinput.Update(msg)

	batch := tea.Batch(
		stopwatchcmd,
		textinputcmd,
		cmd,
	)
	return m, batch
}

func (m model) View() string {
	if m.finished {
		s := "Good job! Your final time was: " + m.stopwatch.View()

		s += "\n\nPress Enter to restart"
		return s
	}

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

		// render remainer of sentence
		s += unwrittenStyle.Render(m.referenceSentence[len(typed):])

		// render help
		s += m.helpView()
	}
	return s
}

// run program
func main() {
	// create and define textinput
	ti := textinput.New()
	ti.Placeholder = "Start Typing!"
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
				key.WithKeys("enter"),
				key.WithHelp("enter", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		// init help
		help: help.NewModel(),
		// TODO fetch wordbank from file
		wordbank: []string{"jupiter", "coffee", "read", "book", "river", "create", "diaspora", "easy", "rumble", "tiddleywink", "funky"},
		finished: false,
		logger:   *errLogger,
	}

	m.keymap.start.SetEnabled(false)

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
