package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles

// When a correct character is typed
var hitStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0BF48B"))

// When an incorrect character is typed
var missStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#F12746")).
	ColorWhitespace(true)

// Characters that haven't been written yet
var unwrittenStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#828282")).
	ColorWhitespace(false)
