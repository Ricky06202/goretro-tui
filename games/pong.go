package games

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Pong struct {
	ballX, ballY      int
	ballVX, ballVY    int
	paddle1Y, paddle2Y int
	score1, score2    int
	width             int
	height            int
	gameOver          bool
	inverted          bool
}

func NewPong() *Pong {
	return &Pong{
		ballX:     30,
		ballY:     10,
		ballVX:    1,
		ballVY:    1,
		paddle1Y: 8,
		paddle2Y: 8,
		width:    60,
		height:   20,
	}
}

type tickMsgPong struct{}

func (p *Pong) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsgPong{}
	})
}

func (p *Pong) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.gameOver {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "r" {
				return NewPong(), nil
			}
		}
		return p, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "a", "W", "A":
			if p.paddle1Y > 0 {
				p.paddle1Y--
			}
		case "s", "S", "d", "D":
			if p.paddle1Y < p.height-4 {
				p.paddle1Y++
			}
		case "up", "k", "K":
			if p.paddle2Y > 0 {
				p.paddle2Y--
			}
		case "down", "j", "J":
			if p.paddle2Y < p.height-4 {
				p.paddle2Y++
			}
		case "r":
			return NewPong(), nil
		case "p":
			p.inverted = !p.inverted
		}
	case tickMsgPong:
		if !p.inverted {
			p.move()
		}
		return p, p.Init()
	}
	return p, nil
}

func (p *Pong) move() {
	p.ballX += p.ballVX
	p.ballY += p.ballVY

	if p.ballY <= 0 || p.ballY >= p.height-1 {
		p.ballVY = -p.ballVY
		if p.ballY < 0 {
			p.ballY = 0
		}
		if p.ballY >= p.height-1 {
			p.ballY = p.height - 1
		}
	}

	if p.ballX == 2 && p.ballY >= p.paddle1Y && p.ballY <= p.paddle1Y+3 {
		p.ballVX = -p.ballVX
		p.ballX = 3
		dy := p.ballY - (p.paddle1Y + 1)
		if dy == 0 {
			p.ballVY = -1
		} else if dy == 1 {
			p.ballVY = 0
		} else {
			p.ballVY = 1
		}
	}
	if p.ballX == p.width-3 && p.ballY >= p.paddle2Y && p.ballY <= p.paddle2Y+3 {
		p.ballVX = -p.ballVX
		p.ballX = p.width - 4
		dy := p.ballY - (p.paddle2Y + 1)
		if dy == 0 {
			p.ballVY = -1
		} else if dy == 1 {
			p.ballVY = 0
		} else {
			p.ballVY = 1
		}
	}

	if p.ballX <= 0 {
		p.score2++
		p.resetBall()
	}
	if p.ballX >= p.width-1 {
		p.score1++
		p.resetBall()
	}

	if p.score1 >= 5 || p.score2 >= 5 {
		p.gameOver = true
	}
}

func (p *Pong) resetBall() {
	p.ballX = p.width / 2
	p.ballY = p.height / 2
	p.ballVX = 1
	p.ballVY = 1
}

func (p *Pong) View() string {
	if p.gameOver {
		winner := "Jugador 1"
		if p.score2 >= 5 {
			winner = "Jugador 2"
		}
		return fmt.Sprintf(`
╔═══════════════════════╗
║   ¡GANADOR!          ║
║   %s         ║
║   %d - %d               ║
╚═══════════════════════╝

r: reiniciar  q: menu
`, winner, p.score1, p.score2)
	}

	status := "▶ JUGANDO"
	if p.inverted {
		status = "⏸ PAUSA"
	}

	s := fmt.Sprintf("🏓 PONG  %s  %d | %d\n\n", status, p.score1, p.score2)

	canvas := make([][]rune, p.height)
	for i := range canvas {
		canvas[i] = make([]rune, p.width)
		for j := range canvas[i] {
			if j == 0 || j == p.width-1 {
				canvas[i][j] = '│'
			} else {
				canvas[i][j] = ' '
			}
		}
	}

	canvas[p.ballY][p.ballX] = '●'

	for y := p.paddle1Y; y < p.paddle1Y+4; y++ {
		if y < p.height {
			canvas[y][1] = '█'
		}
	}
	for y := p.paddle2Y; y < p.paddle2Y+4; y++ {
		if y < p.height {
			canvas[y][p.width-2] = '█'
		}
	}

	for _, row := range canvas {
		s += " " + string(row) + "\n"
	}

	s += "\nW/A/S/D: Jug1  ↑/↓/K/J: Jug2  p: pausa  q: menu"
	return s
}