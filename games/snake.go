package games

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Point struct{ X, Y int }

type Snake struct {
	body      []Point
	food     Point
	dir      Point
	width    int
	height   int
	score    int
	gameOver bool
}

func NewSnake() *Snake {
	rand.Seed(time.Now().UnixNano())
	s := &Snake{
		body:   []Point{{10, 5}, {9, 5}, {8, 5}},
		dir:    Point{1, 0},
		width:  20,
		height: 10,
	}
	s.spawnFood()
	return s
}

func (s *Snake) Init() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

func (s *Snake) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.dir.Y != 1 {
				s.dir = Point{0, -1}
			}
		case "down", "j":
			if s.dir.Y != -1 {
				s.dir = Point{0, 1}
			}
		case "left", "h":
			if s.dir.X != 1 {
				s.dir = Point{-1, 0}
			}
		case "right", "l":
			if s.dir.X != -1 {
				s.dir = Point{1, 0}
			}
		case "r":
			return NewSnake(), nil
		case "q", "esc":
			return nil, tea.Quit
		}
	case tickMsg:
		s.move()
		return s, s.Init()
	}
	return s, nil
}

func (s *Snake) move() {
	head := s.body[0]
	newHead := Point{head.X + s.dir.X, head.Y + s.dir.Y}

	if newHead.X < 0 {
		newHead.X = s.width - 1
	} else if newHead.X >= s.width {
		newHead.X = 0
	}
	if newHead.Y < 0 {
		newHead.Y = s.height - 1
	} else if newHead.Y >= s.height {
		newHead.Y = 0
	}

	for _, p := range s.body[1:] {
		if newHead == p {
			s.gameOver = true
			return
		}
	}

	s.body = append([]Point{newHead}, s.body...)

	if newHead == s.food {
		s.score++
		s.spawnFood()
	} else {
		s.body = s.body[:len(s.body)-1]
	}
}

func (s *Snake) spawnFood() {
	for {
		food := Point{rand.Intn(s.width), rand.Intn(s.height)}
		valid := true
		for _, p := range s.body {
			if food == p {
				valid = false
				break
			}
		}
		if valid {
			s.food = food
			return
		}
	}
}

func (s *Snake) View() string {
	if s.gameOver {
		return fmt.Sprintf(`
╔════════════════════╗
║   GAME OVER!        ║
║   Score: %d          ║
╚════════════════════╝

r: reiniciar  q: menu
`, s.score)
	}

	grid := make([][]rune, s.height)
	for i := range grid {
		grid[i] = make([]rune, s.width)
		for j := range grid[i] {
			grid[i][j] = '░'
		}
	}

	for _, p := range s.body {
		grid[p.Y][p.X] = '█'
	}
	grid[s.food.Y][s.food.X] = '●'

	var b strings.Builder
	for _, row := range grid {
		b.WriteString(" ")
		for _, c := range row {
			b.WriteRune(c)
		}
		b.WriteByte('\n')
	}

	return fmt.Sprintf("🐍 SNAKE  Score: %d\n\n", s.score) + b.String() + "\n↑/↓/←/→ mover · r: reiniciar · q: menu"
}