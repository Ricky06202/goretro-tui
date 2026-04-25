package games

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tetris struct {
	board    [][]int
	boardCol [][]lipgloss.Color
	width   int
	height  int
	piece   []Point
	pieceCol lipgloss.Color
	px, py  int
	score   int
	lines   int
	gameOver bool
	paused  bool
}

var pieceColors = []lipgloss.Color{
	NeonCyan,   // I - cyan
	NeonPurple, // T - purple
	NeonGreen, // L - green
	NeonPink,  // J - pink
	NeonCyan,  // O - cyan
	lipgloss.Color("#FFA500"), // S - orange
	lipgloss.Color("#FF0000"), // Z - red
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
	board := make([][]int, 22)
	boardCol := make([][]lipgloss.Color, 22)
	for i := range board {
		board[i] = make([]int, 25)
		boardCol[i] = make([]lipgloss.Color, 25)
	}

	t := &Tetris{
		board:    board,
		boardCol: boardCol,
		width:   25,
		height:  22,
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
			t.boardCol[y][x] = t.pieceCol
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
	t.pieceCol = pieceColors[r]
	t.px = t.width/2 - 2
	t.py = -2
}

func (t *Tetris) View() string {
	var board strings.Builder

	for y := 0; y < t.height; y++ {
		board.WriteString("│")
		for x := 0; x < t.width; x++ {
			isPiece := false
			pieceCol := t.pieceCol
			for _, p := range t.piece {
				if t.px+p.X == x && t.py+p.Y == y {
					isPiece = true
					break
				}
			}
			if isPiece {
				cell := lipgloss.NewStyle().Foreground(pieceCol).Render("▒")
				board.WriteString(cell)
			} else if t.board[y][x] == 1 {
				cell := lipgloss.NewStyle().Foreground(t.boardCol[y][x]).Render("█")
				board.WriteString(cell)
			} else {
				board.WriteString(" ")
			}
		}
		board.WriteString("│\n")
	}

	boardStr := BoxStyle.Width(60).Align(lipgloss.Center).Render(board.String())

	statusIcon := "▶"
	if t.paused {
		statusIcon = "⏸"
	}
	headerTitle := TitleStyle.Render("🧱 TETRIS")
	headerStatus := HeaderScoreStyle.Render(fmt.Sprintf("Lines: %d  %s", t.lines, statusIcon))

	if t.gameOver {
		gameOverBox := BoxStyle.Width(25).Align(lipgloss.Center)
		content := fmt.Sprintf("\n\n  GAME OVER\n\n  Lines: %d\n  Score: %d\n\n", t.lines, t.score)
		return headerTitle + "\n" + headerStatus + "\n\n" + gameOverBox.Render(content) + "\n" + HelpStyle.Render("R: reiniciar  Q: menu")
	}

	return headerTitle + "\n" + headerStatus + "\n\n" + boardStr + "\n" + HelpStyle.Render("←/→/A/D: mover  •  ↑/W: rotate  •  ↓/S: drop  •  P: pausa  •  Q: menu")
}