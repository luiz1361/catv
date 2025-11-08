package tui

import (
	"fmt"
	"strconv"
	"strings"

	"catv/internal/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AdminModel provides CRUD management of flashcards.
type adminView int

const (
	adminList adminView = iota
	adminCreate
	adminEdit
	adminConfirmDelete
	adminHelp
)

type AdminModel struct {
	flashcards []store.Flashcard
	selected   int
	view       adminView

	// form fields
	questionInput textinput.Model
	answerInput   textinput.Model
	revisitInput  textinput.Model // days until next review

	statusMsg string
	errMsg    string

	storeRef *store.Store
}

func NewAdminModel(storeRef *store.Store, flashcards []store.Flashcard) *AdminModel {
	q := textinput.New()
	q.Placeholder = "Question"
	q.Focus()
	a := textinput.New()
	a.Placeholder = "Answer"
	r := textinput.New()
	r.Placeholder = "Days (e.g. 7)"
	return &AdminModel{
		flashcards:    flashcards,
		selected:      0,
		view:          adminList,
		questionInput: q,
		answerInput:   a,
		revisitInput:  r,
		storeRef:      storeRef,
	}
}

func (m *AdminModel) Init() tea.Cmd { return nil }

func (m *AdminModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.view {
		case adminList:
			switch key {
			case "q":
				return m, tea.Quit
			case "j", "down":
				if m.selected < len(m.flashcards)-1 {
					m.selected++
				}
			case "k", "up":
				if m.selected > 0 {
					m.selected--
				}
			case "c":
				m.view = adminCreate
				m.resetForm()
			case "e":
				if len(m.flashcards) == 0 {
					break
				}
				m.loadSelectedIntoForm()
				m.view = adminEdit
			case "d":
				if len(m.flashcards) == 0 {
					break
				}
				m.view = adminConfirmDelete
			case "r":
				m.reload()
			case "?":
				m.view = adminHelp
			}
		case adminCreate:
			if key == "esc" {
				m.view = adminList
				return m, nil
			}
			if key == "tab" {
				m.cycleFocus()
				return m, nil
			}
			if key == "enter" {
				m.createFlashcard()
				return m, nil
			}
			m.updateInputs(key)
		case adminEdit:
			if key == "esc" {
				m.view = adminList
				return m, nil
			}
			if key == "tab" {
				m.cycleFocus()
				return m, nil
			}
			if key == "enter" {
				m.updateFlashcard()
				return m, nil
			}
			m.updateInputs(key)
		case adminConfirmDelete:
			if key == "y" {
				m.deleteFlashcard()
				return m, nil
			}
			if key == "n" || key == "esc" {
				m.view = adminList
			}
		case adminHelp:
			if key == "esc" {
				m.view = adminList
			}
		}
	}
	return m, nil
}

func (m *AdminModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Admin Flashcards Manager")
	helpBar := lipgloss.NewStyle().Faint(true).Render("[j/k] navigate  [c] create  [e] edit  [d] delete  [r] reload  [?] help  [q] quit")

	if m.errMsg != "" {
		helpBar += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.errMsg)
	} else if m.statusMsg != "" {
		helpBar += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render(m.statusMsg)
	}

	switch m.view {
	case adminList:
		var b strings.Builder
		b.WriteString(header + "\n\n")
		if len(m.flashcards) == 0 {
			b.WriteString("No flashcards. Press c to create.\n")
		} else {
			for i, fc := range m.flashcards {
				prefix := "  "
				if i == m.selected {
					prefix = "> "
				}
				nextDisp := fmt.Sprintf("%d", fc.RevisitIn)
				line := fmt.Sprintf("%s[%d] %s | %s (revisitIn: %s)\n", prefix, fc.ID, truncate(fc.Question, 40), truncate(fc.Answer, 30), nextDisp)
				b.WriteString(line)
			}
		}
		b.WriteString("\n" + helpBar)
		return b.String()
	case adminCreate:
		return fmt.Sprintf("%s\n\nCreate Flashcard\nQuestion: %s\nAnswer: %s\nRevisit (days): %s\n[enter] save  [esc] cancel\n\n%s", header, m.questionInput.View(), m.answerInput.View(), m.revisitInput.View(), helpBar)
	case adminEdit:
		return fmt.Sprintf("%s\n\nEdit Flashcard (ID %d)\nQuestion: %s\nAnswer: %s\nRevisit (days): %s\n[enter] update  [esc] cancel\n\n%s", header, m.flashcards[m.selected].ID, m.questionInput.View(), m.answerInput.View(), m.revisitInput.View(), helpBar)
	case adminConfirmDelete:
		return fmt.Sprintf("%s\n\nDelete Flashcard ID %d? [y/n]\n\n%s", header, m.flashcards[m.selected].ID, helpBar)
	case adminHelp:
		return fmt.Sprintf("%s\n\nHelp\nUse j/k to move, c create, e edit, d delete, r reload.\nesc to return.\n\n%s", header, helpBar)
	}
	return ""
}

// helpers
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (m *AdminModel) resetForm() {
	m.questionInput.SetValue("")
	m.answerInput.SetValue("")
	m.revisitInput.SetValue("")
	m.errMsg = ""
	m.statusMsg = ""
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) loadSelectedIntoForm() {
	fc := m.flashcards[m.selected]
	m.questionInput.SetValue(fc.Question)
	m.answerInput.SetValue(fc.Answer)
	m.revisitInput.SetValue(strconv.Itoa(fc.RevisitIn))
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) updateInputs(key string) {
	// We only simulate typed runes for focused input
	var cmds []tea.Cmd
	if m.questionInput.Focused() {
		m.questionInput, _ = m.questionInput.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	} else if m.answerInput.Focused() {
		m.answerInput, _ = m.answerInput.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	} else if m.revisitInput.Focused() {
		m.revisitInput, _ = m.revisitInput.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	}
	_ = cmds
}

func (m *AdminModel) parseRevisitDays() (int, error) {
	val := strings.TrimSpace(m.revisitInput.Value())
	if val == "" {
		return 0, fmt.Errorf("revisit days required")
	}
	d, err := strconv.Atoi(val)
	if err != nil || d < 0 {
		return 0, fmt.Errorf("invalid days")
	}
	return d, nil
}

func (m *AdminModel) createFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	fc := store.Flashcard{Question: m.questionInput.Value(), Answer: m.answerInput.Value(), File: "manual", RevisitIn: days}
	if err := m.storeRef.InsertFlashcard(fc); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = "Flashcard created"
	m.reload()
	m.view = adminList
}

func (m *AdminModel) updateFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	fc := m.flashcards[m.selected]
	fc.Question = m.questionInput.Value()
	fc.Answer = m.answerInput.Value()
	fc.RevisitIn = days
	if err := m.storeRef.UpdateFlashcardFull(fc); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = "Flashcard updated"
	m.reload()
	m.view = adminList
}

func (m *AdminModel) deleteFlashcard() {
	id := m.flashcards[m.selected].ID
	if err := m.storeRef.DeleteFlashcard(id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = fmt.Sprintf("Deleted flashcard %d", id)
	m.reload()
	if m.selected >= len(m.flashcards) {
		m.selected = len(m.flashcards) - 1
	}
	m.view = adminList
}

func (m *AdminModel) reload() {
	list, err := m.storeRef.GetAllFlashcards()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	m.flashcards = list
}

// cycleFocus switches focus Question -> Answer -> Revisit -> Question
func (m *AdminModel) cycleFocus() {
	if m.questionInput.Focused() {
		m.questionInput.Blur()
		m.answerInput.Focus()
		m.revisitInput.Blur()
		return
	}
	if m.answerInput.Focused() {
		m.answerInput.Blur()
		m.revisitInput.Focus()
		m.questionInput.Blur()
		return
	}
	// default or revisit focused -> go to question
	m.revisitInput.Blur()
	m.questionInput.Focus()
	m.answerInput.Blur()
}
