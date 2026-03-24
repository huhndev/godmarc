package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines the keybindings for the application
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Reload key.Binding
	Search key.Binding
	Tab1   key.Binding
	Tab2   key.Binding
	Tab3   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Tab1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "reports"),
		),
		Tab2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "aggregated"),
		),
		Tab3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "failed"),
		),
	}
}

// RenderHelp renders a simplified help view based on view state
func RenderHelp(showReport bool, activeTab int, searching bool, keys KeyMap) string {
	if searching {
		return "type to filter · esc cancel · enter confirm"
	}
	if showReport {
		return "↑/k up · ↓/j down · esc back · q quit"
	}
	return "↑/k up · ↓/j down · enter select · 1/2/3 tabs · / search · r reload · q quit"
}
