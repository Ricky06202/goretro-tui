package games

import "github.com/charmbracelet/lipgloss"

var (
	NeonPink   = lipgloss.Color("#FF006E")
	NeonCyan   = lipgloss.Color("#00F5FF")
	NeonPurple = lipgloss.Color("#9B5DE5")
	NeonGreen  = lipgloss.Color("#00FF9F")
	DarkBG     = lipgloss.Color("#0D0D0D")
	LightText  = lipgloss.Color("#FAFAFA")
	GrayText   = lipgloss.Color("#888888")
	WarnRed    = lipgloss.Color("#FF4444")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			Bold(true).
			Width(60).
			Align(lipgloss.Center)

	HeaderScoreStyle = lipgloss.NewStyle().
			Foreground(NeonGreen).
			Bold(true).
			Width(60).
			Align(lipgloss.Center)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(NeonPurple).
			Foreground(LightText).
			Padding(1, 2).
			Width(60)

	SelectedStyle = lipgloss.NewStyle().
			Background(NeonCyan).
			Foreground(DarkBG).
			Bold(true).
			Padding(0, 2)

	NormalStyle = lipgloss.NewStyle().
			Foreground(NeonCyan).
			Padding(0, 2)

	StatusStyle = lipgloss.NewStyle().
			Foreground(NeonGreen).
			Bold(true).
			Width(60).
			Align(lipgloss.Center)

	HelpStyle = lipgloss.NewStyle().
			Foreground(GrayText).
			Width(60).
			Align(lipgloss.Center)

	WarnStyle = lipgloss.NewStyle().
			Foreground(WarnRed).
			Bold(true).
			Width(60).
			Align(lipgloss.Center)
)

func StringWidth(s string) int {
	w := 0
	for _, r := range s {
		if r >= 0x1000 {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}