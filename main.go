// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/renatonmag/versionctrls-cli/pkg/repository"
// )

// type model struct {
// 	state       int
// 	textInput   textinput.Model
// 	initialized bool
// 	apiKey      string
// 	rootPath    string
// 	url         string
// }

// const (
// 	stateInitMessage = iota
// 	stateAskAPIKey
// 	stateAskPath
// 	stateAskURL
// 	notInRootPath
// 	generalQuitMsg
// 	clearScreen
// 	end
// )

// const secretFilePath = "apikey.txt"

// func initialModel() model {
// 	ti := textinput.New()
// 	ti.Placeholder = ""
// 	ti.Focus()
// 	ti.CharLimit = 256
// 	ti.Width = 30

// 	return model{
// 		state:     stateInitMessage,
// 		textInput: ti,
// 		rootPath:  ".",
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return textinput.Blink
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "enter":
// 			switch m.state {
// 			case stateInitMessage:
// 				m.state = stateAskAPIKey
// 				m.textInput.Placeholder = "Your GitHub API key"
// 			case stateAskAPIKey:
// 				m.apiKey = m.textInput.Value()
// 				err := os.WriteFile(secretFilePath, []byte(m.apiKey), 0600)
// 				if err != nil {
// 					fmt.Println("Error saving API key:", err)
// 					return m, tea.Quit
// 				}
// 				m.textInput.SetValue("")
// 				m.textInput.Placeholder = "Enter your repository URL"
// 				m.state = stateAskURL
// 			// case stateAskPath:
// 			// 	m.rootPath = m.textInput.Value()
// 			// 	m.textInput.SetValue("")
// 			// 	m.textInput.Placeholder = "Enter your repository URL"
// 			// 	m.state = stateAskURL
// 			case stateAskURL:
// 				m.url = m.textInput.Value()
// 				m.initialized = true
// 				return m, tea.Quit
// 			case notInRootPath:
// 				return m, tea.Quit
// 			case clearScreen:
// 				m.state = end
// 				return m, tea.Quit
// 			case end:
// 				m.state = end
// 				return m, tea.Quit
// 			case generalQuitMsg:
// 				m.state = generalQuitMsg
// 				return m, tea.Quit
// 			}

// 		}
// 	}

// 	var cmd tea.Cmd
// 	m.textInput, cmd = m.textInput.Update(msg)

// 	return m, cmd
// }

// func (m model) View() string {
// 	switch m.state {
// 	case stateInitMessage:
// 		return "Initializing Versionctrls repository.\n\nEnter to continue or q to quit."
// 	case stateAskAPIKey:
// 		return fmt.Sprintf(
// 			"Enter your GitHub API Key:\n\n%s\n\n(Enter to confirm)\n",
// 			m.textInput.View(),
// 		)
// 	// case stateAskPath:
// 	// 	return fmt.Sprintf(
// 	// 		"Enter your local repository path:\n\n%s\n\n(Enter to confirm)\n",
// 	// 		m.textInput.View(),
// 	// 	)
// 	case stateAskURL:
// 		return fmt.Sprintf(
// 			"Your Versionctrls integration repository URL:\n\n%s\n\n(Enter to confirm)\n",
// 			m.textInput.View(),
// 		)
// 	case clearScreen:
// 		return ""
// 	case notInRootPath:
// 		return "You are not in the root of a Git repository.\n\n"
// 	case generalQuitMsg:
// 		return "An error occurred. Please try again.\n\n"
// 	default:
// 		return "Unknown state"
// 	}
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: ctrls <command>")
// 		os.Exit(1)
// 	}

// 	cmd := os.Args[1]

// 	if cmd == "init" {
// 		repo := repository.New()
// 		err := repo.PlainOpen(".")
// 		if err != nil {
// 			fmt.Println("You are not in a Git repository.")
// 			return
// 		}

// 		rootPath, err := repo.GetRepoRoot()
// 		if err != nil {
// 			fmt.Println("You are not in the root of the Git repository.")
// 		}

// 		p := tea.NewProgram(initialModel())
// 		m, err := p.Run()
// 		if err != nil {
// 			fmt.Printf("Alas, there's been an error: %v", err)
// 			os.Exit(1)
// 		}
// 		finalModel := m.(model)
// 		if finalModel.initialized {
// 			apiKey, err := os.ReadFile(secretFilePath)
// 			if err != nil {
// 				fmt.Println("Error reading API key:", err)
// 			} else {
// 				fmt.Printf("API Key: %s\n", string(apiKey))
// 			}

// 			fmt.Printf("\nVersionctrls initialized at: %s\n", rootPath)
// 			fmt.Printf("\nJust hit ctrl+s and you're good. Your files are safe forever.")
// 		}
// 	} else {
// 		fmt.Printf("Unknown command: %s\n", cmd)
// 		os.Exit(1)
// 	}
// }

package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renatonmag/versionctrls-cli/pkg/repository"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Initialize ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Initialize"))
)

const (
	stateInit = iota
	end
)

const secretFilePath = "apikey.txt"

type model struct {
	focusIndex  int
	inputs      []textinput.Model
	initialized bool
	state       int
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "GitHub API key"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 1:
			t.Placeholder = "Versionctrls integration repository URL"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.initialized = true
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case stateInit:
		var b strings.Builder

		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		// b.WriteString(helpStyle.Render("cursor mode is "))
		// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		// b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

		return b.String()

	case end:
		return "Goodbye!"

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

	if cmd == "cleanbranch" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}
		repo.CreateEmptyBranchesForChangedFiles()
		// vPath, err := repo.IntegrationSubmodulePath()
		// if err != nil {
		// 	log.Fatalf("Error getting integration submodule path: %v", err)
		// }
		// fmt.Printf("Integration submodule path: %s\n", vPath)
		// vRepo := repository.New()
		// err = vRepo.PlainOpen(vPath)
		// if err != nil {
		// 	log.Fatalf("Error opening integration submodule: %v", err)
		// }
		// err = vRepo.CreateEmptyBranchesForChangedFiles()
		// if err != nil {
		// 	log.Fatalf("Error creating empty branches in integration submodule: %v", err)
		// }

		// err = vRepo.CommitChangedFiles()
		// if err != nil {
		// 	log.Fatalf("Error committing changes in integration submodule: %v", err)
		// }

	} else if cmd == "userinfo" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}
		name, email, err := repo.GetGitUserInfo()
		if err != nil {
			log.Fatalf("Error getting Git user info: %v", err)
		}

		fmt.Printf("Git user name: %s\n", name)
		fmt.Printf("Git user email: %s\n", email)
	} else if cmd == "changes" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}

		files, err := repo.GetChangedFiles()
		if err != nil {
			fmt.Println("Error getting changed files:", err)
			return
		}
		fmt.Println("Files in root:")
		for _, entry := range files {
			fmt.Println(entry)
		}

		vPath, err := repo.IntegrationSubmodulePath()
		if err != nil {
			log.Fatalf("Error getting integration submodule path: %v", err)
		}

		vRepo := repository.New()
		err = vRepo.PlainOpen(vPath)
		if err != nil {
			log.Fatalf("Error opening integration submodule: %v", err)
		}

		files, err = vRepo.GetChangedFiles()
		if err != nil {
			fmt.Println("Error getting changed files:", err)
			return
		}
		fmt.Println("\n\nFiles in integration:")
		for _, entry := range files {
			fmt.Println(entry)
		}

	} else if cmd == "copy" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}

		files, err := repo.GetChangedFiles()
		if err != nil {
			fmt.Println("Error getting changed files:", err)
			return
		}
		for _, entry := range files {
			fmt.Println(entry)
		}

		fmt.Println("\n\nCopying files to submodule...")
		err = repo.CopyChangedFilesToSubmodule()
		if err != nil {
			fmt.Println("Error copying files to submodule:", err)
			return
		}
	} else if cmd == "removeintegration" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}
		err = repo.RemoveSubmodule()
		if err != nil {
			fmt.Println("Error removing submodule:", err)
			return
		}

		fmt.Printf("\nRun th cmds in a clean branch and merge with your main\n\n")
		fmt.Printf("\ngit add .gitmodules versionctrls-integration")
		fmt.Printf("\ngit commit -m 'Remove versionctrls-integration'\n\n")

	} else if cmd == "init" {
		repo := repository.New()
		err := repo.PlainOpen(".")
		if err != nil {
			fmt.Println("You are not in a Git repository.")
			return
		}

		rootPath, err := repo.GetRepoRoot()
		if err != nil {
			fmt.Println("You are not in the root of the Git repository.")
		}

		p := tea.NewProgram(initialModel())
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		finalModel := m.(model)
		if finalModel.initialized {
			apiKey := finalModel.inputs[0].Value()
			err := os.WriteFile(secretFilePath, []byte(apiKey), 0600)
			if err != nil {
				fmt.Println("Error saving API key:", err)
				return
			}

			exists, err := repo.SubmoduleExists("versionctrls-integration")
			if err != nil {
				fmt.Println("Error checking for submodule:", err)
				return
			}
			if !exists {
				submoduleURL := "https://github.com/renatonmag/gitexperimentsintegration.git"
				submodulePath := "versionctrls-integration"
				err := repo.AddSubmodule(submoduleURL, submodulePath)
				if err != nil {
					fmt.Println("Error adding submodule:", err)
					return
				}
				// fmt.Println(output)
			} else {
				fmt.Println("Versionctrls is already initialized.")
			}

			fmt.Printf("\nVersionctrls initialized at: %s\n", rootPath)
			fmt.Printf("\nCommit the changes to .gitmodules and versionctrls-integration folder.\n\n")
			fmt.Printf("\nJust hit ctrl+s and you're good. Your files are safe forever.")
		}
	} else {
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
