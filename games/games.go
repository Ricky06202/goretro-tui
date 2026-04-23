package games

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Tetris struct {
	board      [][]int
	width     int
	height    int
	score     int
	gameOver  bool
}

func NewTetris() *Tetris {
	return &Tetris{
		board:  nil,
		width:  10,
		height: 20,
	}
}

func (t *Tetris) Init() tea.Cmd { return nil }

func (t *Tetris) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Tetris) View() string {
	return "🧱 TETRIS  - Coming soon!\n\nq: menu"
}

type Pong struct {
	width, height int
}

func NewPong() *Pong { return &Pong{width: 40, height: 15} }

func (p *Pong) Init() tea.Cmd { return nil }

func (p *Pong) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

func (p *Pong) View() string {
	return "🏓 PONG  - Coming soon!\n\nq: menu"
}

type Memory struct {
	width, height int
}

func NewMemory() *Memory {
	rand.Seed(time.Now().UnixNano())
	return &Memory{}
}

func (m *Memory) Init() tea.Cmd { return nil }

func (m *Memory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Memory) View() string {
	return "🧠 MEMORY  - Coming soon!\n\nq: menu"
}