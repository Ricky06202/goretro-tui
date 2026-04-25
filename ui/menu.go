package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"goretro-tui/games"
)

func stringWidth(s string) int {
	w := 0
	for _, r := range s {
		if r >= 0x1000 {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

func newSnake() tea.Model  { return games.NewSnake() }
func newTetris() tea.Model { return games.NewTetris() }
func newPong() tea.Model  { return games.NewPong() }
func newMemory() tea.Model { return games.NewMemory() }
func newWordle() tea.Model { return games.NewWordle() }

type GameMenu struct {
	cursor   int
	options  []string
	submodel tea.Model
	quit     bool
}

func NewGameMenu() *GameMenu {
	return &GameMenu{
		options: []string{
			"🐍  Snake",
			"🧱  Tetris",
			"🏓  Pong",
			"🧠  Memory",
			"📝  Wordle",
			"👋  Salir",
		},
		cursor: 0,
	}
}

func (m *GameMenu) Init() tea.Cmd { return nil }

func (m *GameMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.submodel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" || msg.String() == "esc" {
				m.submodel = nil
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.submodel, cmd = m.submodel.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.quit {
				return m, tea.Quit
			}
			m.quit = true
			return m, nil
		case "up", "k", "w":
			m.quit = false
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j", "s":
			m.quit = false
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.quit = false
			switch m.cursor {
			case 0:
				m.submodel = newSnake()
				return m, m.submodel.Init()
			case 1:
				m.submodel = newTetris()
				return m, m.submodel.Init()
			case 2:
				m.submodel = newPong()
				return m, m.submodel.Init()
			case 3:
				m.submodel = newMemory()
			case 4:
				m.submodel = newWordle()
			case 5:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *GameMenu) View() string {
	if m.submodel != nil {
		return m.submodel.View()
	}

	s := strings.Builder{}

	title := TitleStyle.Render("🎮 GORETRO TUI 🎮")
	s.WriteString("\n")
	s.WriteString(title)
	s.WriteString("\n\n")

	menu := strings.Builder{}
	for i, option := range m.options {
		padded := option + strings.Repeat(" ", 20-len(option))
		if m.cursor == i {
			menu.WriteString(SelectedStyle.Render(padded))
		} else {
			menu.WriteString(NormalStyle.Render(padded))
		}
		menu.WriteString("\n")
	}

	menuStr := menu.String()
	menuBox := BoxStyle.Align(lipgloss.Center).Render(menuStr)
	s.WriteString(menuBox)
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/↓ or WASD: navegar  •  Enter: seleccionar  •  Q: salir"))

	if m.quit {
		s.WriteString("\n")
		s.WriteString(SelectedStyle.Render("Presiona Q otra vez para confirmar"))
	}

	return s.String()
}