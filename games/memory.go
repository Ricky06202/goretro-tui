package games

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var symbols = []rune("🎮🔥🌙⭐🐱🌵🎸🍕")

type Memory struct {
	cards     []rune
	flipped   []int
	matched   []int
	selected  int
	width    int
	height   int
	moves    int
	win      bool
	waiting  bool
}

func NewMemory() *Memory {
	pairs := make([]rune, 16)
	for i := 0; i < 8; i++ {
		pairs[i*2] = symbols[i]
		pairs[i*2+1] = symbols[i]
	}
	rand.Shuffle(16, func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})

	return &Memory{
		cards:    pairs,
		width:   4,
		height:  4,
	}
}

type tickMsgMemory struct{}

func (m *Memory) Init() tea.Cmd {
	if m.waiting {
		return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return tickMsgMemory{}
		})
	}
	return nil
}

func (m *Memory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.waiting {
		switch msg.(type) {
		case tickMsgMemory:
			m.waiting = false
			if len(m.flipped) == 2 {
				idx1, idx2 := m.flipped[0], m.flipped[1]
				if m.cards[idx1] == m.cards[idx2] {
					m.matched = append(m.matched, idx1, idx2)
					if len(m.matched) == 16 {
						m.win = true
					}
				}
				m.flipped = nil
			}
			return m, nil
		}
		return m, nil
	}

	if m.win {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return nil, tea.Quit
		case "r":
			return NewMemory(), nil
		case "up", "k", "w":
			if m.selected >= m.width {
				m.selected -= m.width
			}
		case "down", "j", "s":
			if m.selected < 12 {
				m.selected += m.width
			}
		case "left", "h", "a":
			if m.selected%m.width > 0 {
				m.selected--
			}
		case "right", "l", "d":
			if m.selected%m.width < m.width-1 {
				m.selected++
			}
		case "enter", " ":
			if !contains(m.flipped, m.selected) && !contains(m.matched, m.selected) {
				cmd := m.flipCard(m.selected)
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m *Memory) flipCard(n int) tea.Cmd {
	m.flipped = append(m.flipped, n)

	if len(m.flipped) == 2 {
		m.moves++
		m.waiting = true
		return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return tickMsgMemory{}
		})
	}
	return nil
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func (m *Memory) View() string {
	cellWidth := 13

	var grid strings.Builder
	for y := 0; y < m.height; y++ {
		grid.WriteString("│")
		for x := 0; x < m.width; x++ {
			idx := y*m.width + x
			c := '?'
			if contains(m.flipped, idx) || contains(m.matched, idx) {
				c = m.cards[idx]
			} else if idx == m.selected && !m.waiting {
				c = '◉'
			}
			cell := fmt.Sprintf("%-*c", cellWidth-1, c)
			grid.WriteString(cell)
		}
		grid.WriteString("│\n")
	}

	gridStr := BoxStyle.Width(60).Align(lipgloss.Center).Render(grid.String())

	statusMsg := ""
	if m.waiting {
		statusMsg = "(espera...)"
	}
	headerTitle := TitleStyle.Render("🧠 MEMORY")
	headerScore := HeaderScoreStyle.Render(fmt.Sprintf("Movidas: %d  %s", m.moves, statusMsg))

	if m.win {
		winBox := BoxStyle.Width(25).Align(lipgloss.Center)
		content := fmt.Sprintf("\n\n  ¡GANASTE!\n\n  Movidas: %d\n\n", m.moves)
		return headerTitle + "\n" + headerScore + "\n\n" + winBox.Render(content) + "\n" + HelpStyle.Render("R: reiniciar  Q: menu")
	}

	return headerTitle + "\n" + headerScore + "\n\n" + gridStr + "\n" + HelpStyle.Render("↑/↓/←/→ or WASD: mover  •  Enter: voltear  •  R: reiniciar  •  Q: menu")
}