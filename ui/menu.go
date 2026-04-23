package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"goretro-tui/games"
)

type GameMenu struct {
	cursor   int
	options []string
	submodel tea.Model
}

func NewGameMenu() *GameMenu {
	return &GameMenu{
		options: []string{
			"🐍  Snake",
			"🧱  Tetris",
			"🏓  Pong",
			"🧠  Memory",
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
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.submodel = games.NewSnake()
				return m, m.submodel.Init()
			case 1:
				m.submodel = games.NewTetris()
			case 2:
				m.submodel = games.NewPong()
			case 3:
				m.submodel = games.NewMemory()
			case 4:
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

	s.WriteString(TitleStyle.Render("🎮 GORETRO TUI 🎮"))
	s.WriteString("\n\n")

	menu := strings.Builder{}
	for i, option := range m.options {
		cols := stringWidth(option)
		padding := 12 - cols
		padded := option + strings.Repeat(" ", padding)
		if m.cursor == i {
			menu.WriteString(SelectedStyle.Render("▶ " + padded))
		} else {
			menu.WriteString(NormalStyle.Render("  " + padded))
		}
		menu.WriteString("\n")
	}

	s.WriteString(ContainerStyle.Render(menu.String()))
	s.WriteString(HelpStyle.Render("↑/↓ navegar · Enter: seleccionar · q: salir"))

	return s.String()
}

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