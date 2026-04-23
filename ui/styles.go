package ui

import "github.com/charmbracelet/lipgloss"

var (
	NeonPink   = lipgloss.Color("#FF006E")
	NeonCyan   = lipgloss.Color("#00F5FF")
	NeonPurple = lipgloss.Color("#9B5DE5")
	NeonYellow = lipgloss.Color("#FEE440")
	LightText  = lipgloss.Color("#FAFAFA")

	TitleStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			MarginTop(3).
			Width(50).
			Align(lipgloss.Center)

	ContainerStyle = lipgloss.NewStyle().
			Width(50).
			Align(lipgloss.Center)

	SelectedStyle = lipgloss.NewStyle().
			Background(NeonCyan).
			Foreground(lipgloss.Color("#0D0D0D")).
			Bold(true).
			PaddingLeft(2).
			PaddingRight(2)

	NormalStyle = lipgloss.NewStyle().
			Foreground(NeonCyan)

	HelpStyle = lipgloss.NewStyle().
			Foreground(NeonCyan).
			MarginTop(2)
)