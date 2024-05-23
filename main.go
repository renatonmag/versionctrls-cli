package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state       int
	textInput   textinput.Model
	initialized bool
	apiKey      string
	path        string
	url         string
}

const (
	stateInitMessage = iota
	stateAskAPIKey
	stateAskPath
	stateAskURL
)

const secretFilePath = "apikey.txt"

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 30

	return model{
		state:     stateInitMessage,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case stateInitMessage:
				m.state = stateAskAPIKey
			case stateAskAPIKey:
				m.apiKey = m.textInput.Value()
				err := os.WriteFile(secretFilePath, []byte(m.apiKey), 0600)
				if err != nil {
					fmt.Println("Error saving API key:", err)
					return m, tea.Quit
				}
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Your GitHub API key"
				m.state = stateAskPath
			case stateAskPath:
				m.path = m.textInput.Value()
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Enter your repository URL"
				m.state = stateAskURL
			case stateAskURL:
				m.url = m.textInput.Value()
				m.initialized = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateInitMessage:
		return "Initializing Versionctrls repository.\n\nEnter to continue or q to quit."
	case stateAskAPIKey:
		return fmt.Sprintf(
			"Enter your GitHub API Key:\n\n%s\n\n(Enter to confirm)\n",
			m.textInput.View(),
		)
	case stateAskPath:
		return fmt.Sprintf(
			"Enter your local repository path:\n\n%s\n\n(Enter to confirm)\n",
			m.textInput.View(),
		)
	case stateAskURL:
		return fmt.Sprintf(
			"Your Versionctrls integration repository URL:\n\n%s\n\n(Enter to confirm)\n",
			m.textInput.View(),
		)
	default:
		return "Unknown state"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ctrls <command>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "init" {
		p := tea.NewProgram(initialModel())
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		finalModel := m.(model)
		if finalModel.initialized {
			apiKey, err := os.ReadFile(secretFilePath)
			if err != nil {
				fmt.Println("Error reading API key:", err)
			} else {
				fmt.Printf("API Key: %s\n", string(apiKey))
			}
			fmt.Printf("\nVersionctrls initialized at: %s\n", finalModel.path)
			fmt.Printf("\nJust hit ctrl+s to and that's it. Your files are versioned.")
		}
	} else {
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
