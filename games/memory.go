package games

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var symbols = []rune("🎮🔥🌙⭐🐱🌵🎸🍕")

type Memory struct {
	cards     []rune
	flipped   []int
	matched   []int
	selected int
	width    int
	height   int
	moves    int
	gameOver bool
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
		case "up", "k", "w", "K", "W":
			if m.selected >= m.width {
				m.selected -= m.width
			}
		case "down", "j", "s", "J", "S":
			if m.selected < 12 {
				m.selected += m.width
			}
		case "left", "h", "a", "A":
			if m.selected%m.width > 0 {
				m.selected--
			}
		case "right", "l", "d", "D":
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
	if m.win {
		return fmt.Sprintf(`
╔═══════════════════╗
║   ¡GANASTE!   ║
║  Movidas: %d    ║
╚═══════════════════╝

r: reiniciar  q: menu
`, m.moves)
	}

	status := ""
	if m.waiting {
		status = " (espera...)"
	}

	s := fmt.Sprintf("🧠 MEMORY%s  Movidas: %d\n\n", status, m.moves)

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			idx := y*m.width + x
			c := '?'
			if contains(m.flipped, idx) {
				c = m.cards[idx]
			} else if contains(m.matched, idx) {
				c = m.cards[idx]
				if idx == m.selected {
					c = '◉'
				}
			} else if idx == m.selected && !m.waiting {
				c = '◉'
			}
			s += string(c) + "  "
		}
		s += "\n"
	}

	s += "\n↑/↓/←/→ or WASD: mover  Enter: voltear  q: menu"
	return s
}