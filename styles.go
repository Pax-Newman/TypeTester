package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles

// Default
var defaultStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF"))

// Game Box
var gameBoxStyle = lipgloss.NewStyle().
	Width(50)
	// Border(lipgloss.RoundedBorder())

// When a correct character is typed
var hitStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Foreground(lipgloss.Color("#0BF48B"))

// When an incorrect character is typed
var missStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Background(lipgloss.Color("#F12746")).
	ColorWhitespace(true)

// Characters that haven't been written yet
var unwrittenStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Foreground(lipgloss.Color("#828282")).
	ColorWhitespace(false)

// Error Message
var errorStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#F12746")).
	Padding(0, 1).
	AlignHorizontal(lipgloss.Center)

var cursorStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Background(lipgloss.Color("#828282"))

var pauseStyle = lipgloss.NewStyle().
	Inherit(defaultStyle).
	Width(50).
	Border(lipgloss.RoundedBorder())
