package games

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var words = []string{
	"GOLPE", "MUSEO", "LIBRE", "PLANO", "VISTA",
	"LENTO", "RADIO", "MUNDO", "PLAYA", "SALSA",
	"LUJOS", "CIELO", "BARCO", "NADAR", "SILLA",
	"LIBRO", "PAPEL", "CALOR", "FRENO", "TRENA",
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
	s := fmt.Sprintf("📝 WORDLE  Intentos: %d/6\n\n", w.attempts)

	for i := 0; i < w.maxAttempts; i++ {
		line := "_____"
		if i < len(w.guesses) {
			guess := w.guesses[i]
			result := w.check(guess)
			for j := 0; j < 5; j++ {
				switch result[j] {
				case 'G':
					line = changeChar(line, j, 'G', guess[j])
				case 'Y':
					line = changeChar(line, j, 'Y', guess[j])
				default:
					line = changeChar(line, j, 'X', guess[j])
				}
			}
		} else if i == len(w.guesses) && len(w.current) > 0 {
			line = w.current + strings.Repeat(" ", 5-len(w.current))
		}
		s += line + "\n"
	}

	s += "\nVerde: correcto  Amarillo:的位置  Gris: no existe\n"
	s += "Enter: enviar  Backspace: borrar  q: menu"

	if w.gameOver {
		if w.win {
			s = "📝 WORDLE  ¡GANASTE!\n\n" + s
		} else {
			s = fmt.Sprintf("📝 WORDLE  PERDISTE! Era: %s\n\n%s\n", w.secret, s)
		}
	}

	return s
}

func changeChar(s string, i int, color byte, r byte) string {
	colors := map[byte]string{
		'G': "🟩",
		'Y': "🟨",
		'X': "⬛",
	}
	result := ""
	for j := 0; j < 5; j++ {
		if j == i {
			result += colors[color] + string(r)
		} else {
			result += " _ "
		}
	}
	return result
}