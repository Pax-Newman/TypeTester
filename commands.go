package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type logMsg string
type startMsg struct{}
type stopMsg struct{}
type errMsg struct {
	err  error
	note string
}
type wordbankMsg []string
type seedMsg int64

// Commands
func stopGame() tea.Msg {
	return stopMsg{}
}

func startGame() tea.Msg {
	return startMsg{}
}

func logText(s string) tea.Cmd {
	return func() tea.Msg { return logMsg(s) }
}

func reportError(e error, msg string) tea.Cmd {
	return func() tea.Msg { return errMsg{e, msg} }
}

// Sets a pseudorandom seed for the random generator
func setSeed() tea.Msg {
	var seedbytes [8]byte
	_, err := crypto_rand.Read(seedbytes[:])
	if err != nil {
		return errMsg{err, "Fatal error occured while generating random seed: "}
	}
	return seedMsg(int64(binary.LittleEndian.Uint64(seedbytes[:])))
}

// Loads a list of words from a given file
//
// Unused until https://github.com/charmbracelet/bubbletea/commit/989d49f3e69f2e67951a6b803a2f6973de8443d0
// is implemented on main branch
//
// https://github.com/charmbracelet/bubbletea/issues/413
/* func loadWordbank() tea.Msg {
	f, err := os.Open("wordbank.txt")
	if err != nil {
		return errMsg{err, "Fatal error while loading wordbank file: "}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	words := []string{}
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return errMsg{err, "Fatal error occured while reading wordbank file: "}
	}
	return wordbankMsg(words)
} */
