package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pax-newman/teatime"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

// State Enum
type state int

const (
	playing state = iota
	quitting
	erroring
	finished
	paused
)

// Model
type Model struct {
	// bubble models
	stopwatch teatime.Model
	textinput textinput.Model
	keymap    keymap
	help      help.Model

	// state fields
	state             state
	referenceSentence string
	wordbank          []string
	err               error
	errNote           string

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

func loadWordbank() []string {
	f, err := os.Open("wordbank.txt")
	if err != nil {
		log.Fatal("Error opening wordbank file: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	words := []string{}
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading wordbank file: ", err)
	}
	return words
}

// Methods

// render help page
func (m Model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

// first method to run when model is created
func (m Model) Init() tea.Cmd {
	// init bubble models
	batch := tea.Batch(
		m.stopwatch.Init(),
		textinput.Blink,
		setSeed, // currently has a chance to not execute before the game starts
		startGame,
	)
	return batch
}

// update model's state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// report error and quit game
	case errMsg:
		m.err = msg.err
		m.errNote = msg.note
		m.state = erroring
		return m, tea.Quit
	// log incoming message
	case logMsg:
		m.logger.Println("From message " + msg)
		return m, nil
	// load wordbank
	case wordbankMsg:
		m.wordbank = []string(msg)
		return m, nil
	// set random seed
	case seedMsg:
		rand.Seed(int64(msg))
		return m, nil
	// if the message is a keystroke
	case tea.KeyMsg:
		switch {
		// quit the game
		case key.Matches(msg, m.keymap.quit):
			m.state = quitting
			return m, tea.Quit

		// reset the timer
		case key.Matches(msg, m.keymap.reset):
			return m, startGame

		// toggle paused/playing state
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			if m.state == paused {
				m.state = playing
			} else {
				m.state = paused
			}
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()

		default:
			// check completion status
			if val := m.textinput.Value() + msg.String(); len(val) > 0 && val == m.referenceSentence {
				m.textinput.Reset()
				m.state = finished
				return m, stopGame
			}
		}
	// Stop the game when it's finished
	case stopMsg:
		return m, m.stopwatch.Stop()
	// Start a new game
	case startMsg:
		m.textinput.Reset()
		m.referenceSentence = createNewSentence(10, m.wordbank)
		m.textinput.CharLimit = len(m.referenceSentence)
		m.state = playing
		cmd := tea.Batch(
			m.stopwatch.Reset(),
			m.stopwatch.Start(),
		)
		return m, cmd
	}

	// update bubbles
	var stopwatchcmd tea.Cmd
	m.stopwatch, stopwatchcmd = m.stopwatch.Update(msg)

	var textinputcmd tea.Cmd
	if m.state == playing {
		m.textinput, textinputcmd = m.textinput.Update(msg)
	}

	batch := tea.Batch(
		stopwatchcmd,
		textinputcmd,
		cmd,
	)
	return m, batch
}

func (m Model) View() string {

	var s string

	switch m.state {

	case playing:
		// render timer
		timer := "Elapsed: " + m.stopwatch.View() + "\n"
		s = defaultStyle.Render(timer)

		// render reference sentence
		refSent := "\n" + m.referenceSentence + "\n"
		s += defaultStyle.Render(refSent)

		// render typed sentence
		typed := []rune(m.textinput.Value())
		ref := []rune(m.referenceSentence)
		cursorLocation := m.textinput.Cursor()
		var styledInput string

		for idx, r := range typed {
			if r == ref[idx] {
				styledInput += hitStyle.Render(string(r))
			} else {
				styledInput += missStyle.Render(string(r))
			}
		}

		// render cursor
		offset := 0
		if len(typed) < len(ref) {
			offset++
			styledInput += cursorStyle.Render(string(ref[cursorLocation]))
		}

		// render remainer of sentence
		remainder := m.referenceSentence[len(typed)+offset:]
		styledInput += unwrittenStyle.Render(remainder)

		// join typed input to rest of view
		s += styledInput

		// wrap text in a box
		s = gameBoxStyle.Render(s)

		// render help
		s += "\n" + m.helpView()

	case paused:
		s = lipgloss.PlaceHorizontal(
			gameBoxStyle.GetWidth(),
			lipgloss.Center,
			"\n\nPaused\n\n",
		)
		s = pauseStyle.Render(s)

	case finished:
		s = "Good job! Your final time was: " + m.stopwatch.View()

		s += "\n\nPress Enter to restart"

	case quitting:
		s = "Elapsed: " + m.stopwatch.View() + "\n"

	case erroring:
		s = fmt.Sprint(m.errNote, m.err)
		s = errorStyle.Render(s) + "\n"
	}

	return s
}

// run program
func main() {
	setSeed()

	// create and define textinput
	ti := textinput.New()
	ti.Focus()

	f, _ := os.OpenFile("TypeTester.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	errLogger := log.New(f, "", log.Lshortfile)

	m := Model{
		// create stopwatch
		stopwatch: teatime.NewWithInterval(time.Millisecond),
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
		help:     help.NewModel(),
		state:    playing,
		logger:   *errLogger,
		wordbank: loadWordbank(),
	}

	m.keymap.start.SetEnabled(false)

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
