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
		body:   []Point{{10, 5}, {9, 5}, {8, 5}},
		dir:    Point{1, 0},
		width:  20,
		height: 10,
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

	return fmt.Sprintf("🐍 SNAKE  Score: %d\n\n", s.score) + b.String() + "\n↑/↓/←/→ or WASD: mover · r: reiniciar · q: menu"
}

type Tetris struct {
	board    [][]int
	width   int
	height  int
	piece   []Point
	px, py  int
	score   int
	lines   int
	gameOver bool
	paused  bool
}

var tetrominoes = [][]Point{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{0, 0}, {1, 0}, {2, 0}, {1, 1}},
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
	{{0, 0}, {1, 0}, {2, 0}, {0, 1}},
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
}

func NewTetris() *Tetris {
	board := make([][]int, 20)
	for i := range board {
		board[i] = make([]int, 15)
	}

	t := &Tetris{
		board:  board,
		width:  15,
		height: 20,
	}
	t.spawnPiece()
	return t
}

type tickMsgTetris struct{}

func (t *Tetris) Init() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsgTetris{}
	})
}

func (t *Tetris) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if t.gameOver || t.paused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "r" {
				return NewTetris(), nil
			}
			if msg.String() == "p" {
				t.paused = !t.paused
				if !t.paused {
					return t, t.Init()
				}
			}
		}
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "a":
			if !t.collides(-1, 0) {
				t.px--
			}
		case "right", "l", "d":
			if !t.collides(1, 0) {
				t.px++
			}
		case "down", "j", "s":
			if !t.collides(0, 1) {
				t.py++
			}
		case "up", "k", "w":
			t.rotate()
		case "ctrl+c", "q", "esc":
			return nil, tea.Quit
		case "r":
			return NewTetris(), nil
		case "p":
			t.paused = !t.paused
			if !t.paused {
				return t, t.Init()
			}
		}
	case tickMsgTetris:
		if !t.collides(0, 1) {
			t.py++
		} else {
			t.lockPiece()
			t.clearLines()
			t.spawnPiece()
			if t.collides(0, 0) {
				t.gameOver = true
			}
		}
		return t, t.Init()
	}
	return t, nil
}

func (t *Tetris) collides(dx, dy int) bool {
	for _, p := range t.piece {
		x := t.px + p.X + dx
		y := t.py + p.Y + dy
		if x < 0 || x >= t.width || y >= t.height {
			return true
		}
		if y >= 0 && x >= 0 && x < t.width && t.board[y][x] == 1 {
			return true
		}
	}
	return false
}

func (t *Tetris) rotate() {
	newPiece := make([]Point, len(t.piece))
	for i, p := range t.piece {
		newPiece[i] = Point{X: -p.Y, Y: p.X}
	}
	old := t.piece
	t.piece = newPiece
	if t.collides(0, 0) {
		t.piece = old
	}
}

func (t *Tetris) lockPiece() {
	for _, p := range t.piece {
		x := t.px + p.X
		y := t.py + p.Y
		if y >= 0 && y < t.height && x >= 0 && x < t.width {
			t.board[y][x] = 1
		}
	}
}

func (t *Tetris) clearLines() {
	for y := t.height - 1; y >= 0; y-- {
		full := true
		for x := 0; x < t.width; x++ {
			if t.board[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			copy(t.board[1:], t.board[:y])
			t.board[0] = make([]int, t.width)
			t.lines++
			t.score += 100
			y++
		}
	}
}

func (t *Tetris) spawnPiece() {
	r := rand.Intn(len(tetrominoes))
	t.piece = make([]Point, 4)
	for i := range t.piece {
		t.piece[i] = tetrominoes[r][i]
	}
	t.px = t.width/2 - 2
	t.py = -2
}

func (t *Tetris) View() string {
	if t.gameOver {
		return fmt.Sprintf(`
╔═══════════════╗
║   GAME OVER   ║
║  Lines: %d    ║
║  Score: %d    ║
╚══════════���════╝

r: reiniciar  q: menu
`, t.lines, t.score)
	}

	status := "▶"
	if t.paused {
		status = "⏸"
	}

	s := fmt.Sprintf("🧱 TETRIS  %s  Lines: %d  Score: %d\n\n", status, t.lines, t.score)

	disp := make([][]rune, t.height)
	for i := range disp {
		disp[i] = make([]rune, 0, t.width+2)
		disp[i] = append(disp[i], '│')
		for j := 0; j < t.width; j++ {
			disp[i] = append(disp[i], ' ')
		}
		disp[i] = append(disp[i], '│')
	}

	for y := 0; y < t.height; y++ {
		for x := 0; x < t.width; x++ {
			if t.board[y][x] == 1 {
				disp[y][x+1] = '█'
			}
		}
	}

	for _, p := range t.piece {
		x := t.px + p.X + 1
		y := t.py + p.Y
		if y >= 0 && y < t.height && x > 0 && x <= t.width {
			disp[y][x] = '▒'
		}
	}

	for _, row := range disp {
		s += " " + string(row) + "\n"
	}
	s += " ═════════════════════════════════════════════════════════"

	s += "\n←/→/A/D: mover  ↑/W/K: rotate  ↓/J/S: drop  p: pausa  q: menu"
	return s
}