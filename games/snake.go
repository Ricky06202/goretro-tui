package games

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Point struct{ X, Y int }

type Snake struct {
	body     []Point
	food     Point
	dir      Point
	width   int
	height  int
	score   int
	gameOver bool
}

func NewSnake() *Snake {
	s := &Snake{
		body:   []Point{{15, 8}, {14, 8}, {13, 8}},
		dir:    Point{1, 0},
		width:  30,
		height: 15,
	}
	s.spawnFood()
	return s
}

type tickMsgSnake struct{}

func (s *Snake) Init() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsgSnake{}
	})
}

func (s *Snake) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "w":
			if s.dir.Y != 1 {
				s.dir = Point{0, -1}
			}
		case "down", "j", "s":
			if s.dir.Y != -1 {
				s.dir = Point{0, 1}
			}
		case "left", "h", "a":
			if s.dir.X != 1 {
				s.dir = Point{-1, 0}
			}
		case "right", "l", "d":
			if s.dir.X != -1 {
				s.dir = Point{1, 0}
			}
		case "r":
			return NewSnake(), nil
		case "q", "esc":
			return nil, tea.Quit
		}
	case tickMsgSnake:
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

	var board strings.Builder
	for _, row := range grid {
		board.WriteString("│")
		for _, c := range row {
			board.WriteRune(c)
		}
		board.WriteString("│\n")
	}

	boardStr := board.String()
	boardBox := BoxStyle.Width(60).Align(lipgloss.Center).Render(boardStr)

	headerTitle := TitleStyle.Render("🐍 SNAKE")
	headerScore := HeaderScoreStyle.Render(fmt.Sprintf("Score: %d", s.score))

	if s.gameOver {
		return headerTitle + "\n" + headerScore + "\n\n" + boardBox + "\n" + HelpStyle.Render("R: reiniciar  •  Q: menu")
	}

	return headerTitle + "\n" + headerScore + "\n\n" + boardBox + "\n" + HelpStyle.Render("↑/↓/←/→ or WASD: mover  •  R: reiniciar  •  Q: menu")
}