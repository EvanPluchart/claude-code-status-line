package wizard

import (
	"fmt"
	"os"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is the bubbletea model for the wizard.
type Model struct {
	steps      []StepDef
	stepIndex  int
	cursor     int
	scroll     int        // scroll offset for long lists
	maxVisible int        // max visible items
	selections map[StepID]string          // single-select results
	toggled    map[StepID]map[string]bool // multi-select state
	toggleOrder map[StepID][]string       // ordered widget list per step

	confirmed bool
	cancelled bool
	isInit    bool
}

func newModel(existing *config.Config, isInit bool) Model {
	steps := buildSteps()
	m := Model{
		steps:       steps,
		selections:  make(map[StepID]string),
		toggled:     make(map[StepID]map[string]bool),
		toggleOrder: make(map[StepID][]string),
		maxVisible:  12,
		isInit:      isInit,
	}

	// Pre-populate from existing config
	if existing != nil {
		m.selections[StepLocale] = existing.Locale
		m.selections[StepTheme] = existing.Theme
		m.selections[StepCurrency] = existing.Widgets.Cost.Currency

		// Pre-populate widget lines (strip auto-inserted separators/spacers)
		lineSteps := []StepID{StepLine1, StepLine2, StepLine3}
		for i, step := range lineSteps {
			if i < len(existing.Lines) {
				m.toggled[step] = make(map[string]bool)
				var order []string
				for _, w := range existing.Lines[i].Widgets {
					if w == "separator" || w == "spacer" {
						continue
					}
					m.toggled[step][w] = true
					order = append(order, w)
				}
				m.toggleOrder[step] = order
			}
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	step := m.steps[m.stepIndex]

	switch msg.String() {
	case "ctrl+c", "esc":
		m.cancelled = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.scroll {
				m.scroll = m.cursor
			}
		}

	case "down", "j":
		if m.cursor < len(step.Choices)-1 {
			m.cursor++
			if m.cursor >= m.scroll+m.maxVisible {
				m.scroll = m.cursor - m.maxVisible + 1
			}
		}

	case " ":
		if step.IsMulti {
			m.toggleWidget(step)
		}

	case "enter":
		if step.IsMulti {
			// Advance to next step
			return m.nextStep()
		}
		// Single-select: pick and advance
		m.selections[step.ID] = step.Choices[m.cursor].Value
		return m.nextStep()

	case "backspace":
		if m.stepIndex > 0 {
			m.stepIndex--
			m.cursor = 0
			m.scroll = 0
			// Restore cursor position to current selection
			m.restoreCursor()
		}
	}

	return m, nil
}

func (m *Model) toggleWidget(step StepDef) {
	if m.toggled[step.ID] == nil {
		m.toggled[step.ID] = make(map[string]bool)
	}

	val := step.Choices[m.cursor].Value
	if m.toggled[step.ID][val] {
		// Toggle off: remove from order
		m.toggled[step.ID][val] = false
		order := m.toggleOrder[step.ID]
		for i, v := range order {
			if v == val {
				m.toggleOrder[step.ID] = append(order[:i], order[i+1:]...)
				break
			}
		}
	} else {
		// Enforce max widgets per line
		if len(m.toggleOrder[step.ID]) >= maxWidgetsPerLine {
			return
		}
		// Toggle on: append to order
		m.toggled[step.ID][val] = true
		m.toggleOrder[step.ID] = append(m.toggleOrder[step.ID], val)
	}
}

func (m Model) nextStep() (tea.Model, tea.Cmd) {
	if m.stepIndex >= len(m.steps)-1 {
		// Confirm step
		step := m.steps[m.stepIndex]
		m.selections[step.ID] = step.Choices[m.cursor].Value
		m.confirmed = m.selections[StepConfirm] == "yes"
		return m, tea.Quit
	}

	m.stepIndex++
	m.cursor = 0
	m.scroll = 0
	m.restoreCursor()

	return m, nil
}

func (m *Model) restoreCursor() {
	step := m.steps[m.stepIndex]
	if step.IsMulti {
		return // No single selection to restore
	}
	if val, ok := m.selections[step.ID]; ok {
		for i, c := range step.Choices {
			if c.Value == val {
				m.cursor = i
				if m.cursor >= m.maxVisible {
					m.scroll = m.cursor - m.maxVisible + 1
				}
				return
			}
		}
	}
}

func (m Model) View() string {
	var b strings.Builder

	step := m.steps[m.stepIndex]

	// Header
	b.WriteString("\n")
	b.WriteString("  \x1b[1mClaude Code Status Line\x1b[0m")

	// Progress dots
	b.WriteString("                  ")
	for i := range m.steps {
		if i < m.stepIndex {
			b.WriteString("\x1b[32m●\x1b[0m")
		} else if i == m.stepIndex {
			b.WriteString("\x1b[36m●\x1b[0m")
		} else {
			b.WriteString("\x1b[2m○\x1b[0m")
		}
	}
	b.WriteString("\n")

	// Step title
	b.WriteString(fmt.Sprintf("  Step %d/%d: %s\n\n", m.stepIndex+1, len(m.steps), step.Title))

	// Choices
	if step.ID == StepConfirm {
		m.renderConfirmView(&b, step)
	} else if step.IsMulti {
		m.renderMultiSelect(&b, step)
	} else {
		m.renderSingleSelect(&b, step)
	}

	// Preview
	b.WriteString("\n")
	b.WriteString(renderPreview(&m))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	if step.IsMulti {
		b.WriteString("  \x1b[2m↑↓ navigate  Space toggle  Enter confirm  Backspace back  Esc quit\x1b[0m")
	} else {
		b.WriteString("  \x1b[2m↑↓ navigate  Enter select  Backspace back  Esc quit\x1b[0m")
	}
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderSingleSelect(b *strings.Builder, step StepDef) {
	end := len(step.Choices)
	if end > m.scroll+m.maxVisible {
		end = m.scroll + m.maxVisible
	}

	if m.scroll > 0 {
		b.WriteString("    \x1b[2m↑ more\x1b[0m\n")
	}

	for i := m.scroll; i < end; i++ {
		c := step.Choices[i]
		if i == m.cursor {
			b.WriteString(fmt.Sprintf("    \x1b[36m❯\x1b[0m \x1b[1m%-14s\x1b[0m %s\n", c.Value, c.Label))
		} else {
			b.WriteString(fmt.Sprintf("      %-14s \x1b[2m%s\x1b[0m\n", c.Value, c.Label))
		}
	}

	if end < len(step.Choices) {
		b.WriteString("    \x1b[2m↓ more\x1b[0m\n")
	}
}

func (m Model) renderMultiSelect(b *strings.Builder, step StepDef) {
	end := len(step.Choices)
	if end > m.scroll+m.maxVisible {
		end = m.scroll + m.maxVisible
	}

	// Show selected count and order
	order := m.toggleOrder[step.ID]
	if len(order) > 0 {
		b.WriteString(fmt.Sprintf("    \x1b[2mSelected (%d/%d):\x1b[0m ", len(order), maxWidgetsPerLine))
		for i, w := range order {
			if i > 0 {
				b.WriteString(" \x1b[2m→\x1b[0m ")
			}
			b.WriteString(fmt.Sprintf("\x1b[36m%s\x1b[0m", w))
		}
		b.WriteString("\n\n")
	}

	if m.scroll > 0 {
		b.WriteString("    \x1b[2m↑ more\x1b[0m\n")
	}

	for i := m.scroll; i < end; i++ {
		c := step.Choices[i]
		checked := m.toggled[step.ID] != nil && m.toggled[step.ID][c.Value]
		cursor := "  "
		if i == m.cursor {
			cursor = "\x1b[36m❯\x1b[0m "
		}

		box := "☐"
		if checked {
			box = "\x1b[32m☑\x1b[0m"
		}

		if i == m.cursor {
			b.WriteString(fmt.Sprintf("    %s%s \x1b[1m%-16s\x1b[0m %s\n", cursor, box, c.Value, c.Label))
		} else {
			b.WriteString(fmt.Sprintf("    %s%s %-16s \x1b[2m%s\x1b[0m\n", cursor, box, c.Value, c.Label))
		}
	}

	if end < len(step.Choices) {
		b.WriteString("    \x1b[2m↓ more\x1b[0m\n")
	}
}

func (m Model) renderConfirmView(b *strings.Builder, step StepDef) {
	// Show summary of selections
	b.WriteString("    \x1b[2mLocale:\x1b[0m    ")
	b.WriteString(m.selectionLabel(StepLocale))
	b.WriteString("\n")
	b.WriteString("    \x1b[2mTheme:\x1b[0m     ")
	b.WriteString(m.selectionLabel(StepTheme))
	b.WriteString("\n")
	b.WriteString("    \x1b[2mCurrency:\x1b[0m  ")
	b.WriteString(m.selectionLabel(StepCurrency))
	b.WriteString("\n")

	for i, sid := range []StepID{StepLine1, StepLine2, StepLine3} {
		order := m.toggleOrder[sid]
		label := fmt.Sprintf("Line %d:", i+1)
		b.WriteString(fmt.Sprintf("    \x1b[2m%-11s\x1b[0m", label))
		if len(order) == 0 {
			b.WriteString("\x1b[2m(empty)\x1b[0m")
		} else {
			b.WriteString(strings.Join(order, " \x1b[2m│\x1b[0m "))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	for i, c := range step.Choices {
		if i == m.cursor {
			b.WriteString(fmt.Sprintf("    \x1b[36m❯\x1b[0m \x1b[1m%s\x1b[0m\n", c.Label))
		} else {
			b.WriteString(fmt.Sprintf("      \x1b[2m%s\x1b[0m\n", c.Label))
		}
	}
}

func (m Model) selectionLabel(id StepID) string {
	val, ok := m.selections[id]
	if !ok {
		return "\x1b[2m(default)\x1b[0m"
	}
	for _, step := range m.steps {
		if step.ID == id {
			for _, c := range step.Choices {
				if c.Value == val {
					return fmt.Sprintf("%s \x1b[2m(%s)\x1b[0m", c.Value, c.Label)
				}
			}
		}
	}
	return val
}

// Run launches the interactive wizard and returns the resulting config.
// If existing is nil, defaults are used. isInit controls the title wording.
// Returns (config, confirmed, error). confirmed is false if the user cancelled.
func Run(existing *config.Config, isInit bool) (*config.Config, bool, error) {
	if existing == nil {
		existing = config.Default()
	}

	m := newModel(existing, isInit)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))

	final, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	result := final.(Model)
	if result.cancelled {
		return nil, false, nil
	}

	if !result.confirmed {
		return nil, false, nil
	}

	cfg := buildPreviewConfig(&result)
	// Preserve advanced settings from existing config
	cfg.Widgets.TokenBar = existing.Widgets.TokenBar
	cfg.Widgets.Separator = existing.Widgets.Separator
	cfg.Widgets.Model = existing.Widgets.Model
	cfg.Widgets.Timestamp = existing.Widgets.Timestamp
	cfg.Widgets.Cost.Decimals = existing.Widgets.Cost.Decimals
	cfg.Thresholds = existing.Thresholds

	return cfg, true, nil
}
