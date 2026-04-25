package games

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var words = []string{
	"GOLPE", "MUSEO", "LIBRE", "PLANO", "VISTA",
	"LENTO", "RADIO", "MUNDO", "PLAYA", "SALSA",
	"LUJOS", "CIELO", "BARCO", "NADAR", "SILLA",
}

type Wordle struct {
	secret       string
	guesses      []string
	current      string
	attempts     int
	maxAttempts  int
	gameOver     bool
	win          bool
}

func NewWordle() *Wordle {
	secret := words[rand.Intn(len(words))]
	return &Wordle{
		secret:      secret,
		maxAttempts: 6,
	}
}

func (w *Wordle) Init() tea.Cmd { return nil }

func (w *Wordle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if w.gameOver {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "r" {
				return NewWordle(), nil
			}
		}
		return w, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()

		if s == "enter" && len(w.current) == 5 {
			w.guesses = append(w.guesses, w.current)
			w.attempts++

			if w.current == w.secret {
				w.win = true
				w.gameOver = true
			} else if w.attempts >= w.maxAttempts {
				w.gameOver = true
			}
			w.current = ""
			return w, nil
		}

		if s == "backspace" && len(w.current) > 0 {
			w.current = w.current[:len(w.current)-1]
			return w, nil
		}

		if len(s) == 1 && len(w.current) < 5 {
			c := strings.ToUpper(s)
			if c >= "A" && c <= "Z" {
				w.current += c
			}
		}
	}
	return w, nil
}

func (w *Wordle) check(guess string) string {
	result := make([]byte, 5)
	secret := []byte(w.secret)
	used := make([]bool, 5)

	for i := 0; i < 5; i++ {
		if guess[i] == secret[i] {
			result[i] = 'G'
			used[i] = true
		}
	}

	for i := 0; i < 5; i++ {
		if result[i] == 'G' {
			continue
		}
		found := false
		for j := 0; j < 5; j++ {
			if !used[j] && guess[i] == secret[j] {
				result[i] = 'Y'
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			result[i] = 'X'
		}
	}
	return string(result)
}

func (w *Wordle) View() string {
		var grid strings.Builder

	for i := 0; i < w.maxAttempts; i++ {
		row := ""
		if i < len(w.guesses) {
			guess := w.guesses[i]
			result := w.check(guess)
			for j := 0; j < 5; j++ {
				color := result[j]
				emoji := ""
				switch color {
				case 'G':
					emoji = "🟩"
				case 'Y':
					emoji = "🟨"
				default:
					emoji = "⬛"
				}
				row += emoji + string(guess[j]) + " "
			}
		} else if i == len(w.guesses) && len(w.current) > 0 {
			for _, c := range w.current {
				row += "⬛" + string(c) + " "
			}
			for j := len(w.current); j < 5; j++ {
				row += "⬛_ "
			}
		} else {
			row = "⬛_ ⬛_ ⬛_ ⬛_ ⬛_"
		}
		grid.WriteString(row + "\n")
	}

	gridStr := BoxStyle.Width(60).Align(lipgloss.Center).Render(grid.String())

	headerTitle := TitleStyle.Render("📝 WORDLE")
	headerScore := HeaderScoreStyle.Render(fmt.Sprintf("Intentos: %d/6", w.attempts))

	if w.gameOver {
		resultBox := BoxStyle.Width(25).Align(lipgloss.Center)
		if w.win {
			content := fmt.Sprintf("\n\n  ¡GANASTE!\n\n  Intentos: %d\n\n", w.attempts)
			return headerTitle + "\n" + headerScore + "\n\n" + resultBox.Render(content) + "\n" + HelpStyle.Render("R: reiniciar  Q: menu")
		}
		content := fmt.Sprintf("\n\n  PERDISTE\n\n  Era: %s\n\n", w.secret)
		return headerTitle + "\n" + headerScore + "\n\n" + resultBox.Render(content) + "\n" + HelpStyle.Render("R: reiniciar  Q: menu")
	}

	return headerTitle + "\n" + headerScore + "\n\n" + gridStr + "\n" + HelpStyle.Render("Type: escribir  •  Enter: enviar  •  Backspace: borrar  •  Q: menu")
}