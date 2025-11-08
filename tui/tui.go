package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	QuestionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	AnswerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	InfoStyle     = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244"))
	// ErrorStyle removed (unused)
)

// PrintInfo prints an informational message to the console
func PrintInfo(message string) {
	fmt.Println(InfoStyle.Render(message))
}

// PrintError prints an error message to the console
func PrintError(message string, err error) {
	fmt.Println(QuestionStyle.Render(message), err)
}

// PrintSuccess prints a success message to the console
func PrintSuccess(message string) {
	fmt.Println(QuestionStyle.Render(message))
}
