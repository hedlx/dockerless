package runtime

import (
	"encoding/json"
	"fmt"

	api "github.com/hedlx/doless/client"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	RMCInitStep    = 0
	RMCNameStep    = iota
	RMCLoadingStep = iota
)

type RuntimeCreateResponse struct {
	Runtime *api.Runtime
	Err     error
}

type RuntimeCreateResponseMsg struct {
	Resp *RuntimeCreateResponse
}

type runtimeCreateStartMsg struct{}

func runtimeCreateStart() tea.Msg {
	return runtimeCreateStartMsg{}
}

type RuntimeCreator interface {
	Create(name string, path string) tea.Cmd
}

type RuntimeCreateModel struct {
	Name    string
	Path    string
	Creator RuntimeCreator

	static string

	resp *RuntimeCreateResponse

	step           int
	nameInput      textinput.Model
	loadingSpinner spinner.Model
}

func InitRuntimeCreateModel(m *RuntimeCreateModel) *RuntimeCreateModel {
	m.nameInput = textinput.New()
	m.nameInput.CharLimit = 156
	m.nameInput.Placeholder = "Runtime name"

	m.loadingSpinner = spinner.New()

	m.loadingSpinner.Spinner = spinner.Dot

	return m
}

func (m RuntimeCreateModel) Init() tea.Cmd {
	return runtimeCreateStart
}

func (m RuntimeCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case runtimeCreateStartMsg:
		return m.incStep()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	switch m.step {
	case RMCNameStep:
		return m.handleRMNameStep(msg)
	case RMCLoadingStep:
		return m.handleRMLoadingStep(msg)
	}
	return m, nil
}

func (m RuntimeCreateModel) View() string {
	static := m.static
	if static != "" {
		static += "\n\n"
	}

	active := ""
	if m.step == RMCNameStep {
		active = m.nameInput.View()
	} else if m.step == RMCLoadingStep {
		active = fmt.Sprintf("%s Creating runtime...", m.loadingSpinner.View())
	}

	return fmt.Sprintf("%s%s", static, active)
}

func (m RuntimeCreateModel) handleRMNameStep(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Name = m.nameInput.Value()
			return m.incStep()
		case tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

func (m RuntimeCreateModel) handleRMLoadingStep(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RuntimeCreateResponseMsg:
		m.resp = msg.Resp
		return m.incStep()
	}

	return m, nil
}

func (m *RuntimeCreateModel) incStep(cmds ...tea.Cmd) (*RuntimeCreateModel, tea.Cmd) {
	if m.step == RMCInitStep {
		m.step++
		m.static = fmt.Sprintf("Path: %s", m.Path)
		return m.incStep(m.nameInput.SetCursorMode(textinput.CursorBlink), m.nameInput.Focus())
	}

	if m.step == RMCNameStep && m.Name != "" {
		m.step++
		m.nameInput.Blur()
		m.static = fmt.Sprintf("%s\nName: %s", m.static, m.Name)

		return m.incStep(m.Creator.Create(m.Name, m.Path))
	}

	if m.step == RMCLoadingStep && m.resp != nil {
		m.step++
		if m.resp.Err != nil {
			m.static = fmt.Sprintf("%s\n\nFailed to create runtime: %s", m.static, m.resp.Err)
		} else {
			j, _ := json.MarshalIndent(m.resp.Runtime, "", "  ")
			m.static = fmt.Sprintf("%s\n\n%s", m.static, j)
		}

		return m.incStep(tea.Quit)
	}

	return m, tea.Batch(cmds...)
}