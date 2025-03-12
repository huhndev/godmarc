//Copyright (c) 2025, Julian Huhn
//
//Permission to use, copy, modify, and/or distribute this software for any
//purpose with or without fee is hereby granted, provided that the above
//copyright notice and this permission notice appear in all copies.
//
//THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
//WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
//MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
//ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
//WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
//ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
//OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines the keybindings for the application
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Aggr   key.Binding
	Quit   key.Binding
	Reload key.Binding
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
		Aggr: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "aggregated"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload"),
		),
	}
}

// ListKeyMap implements help.KeyMap for list view
type ListKeyMap struct {
	Keys KeyMap
}

// ShortHelp returns the short help message
func (k ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Keys.Up,
		k.Keys.Down,
		k.Keys.Enter,
		k.Keys.Aggr,
		k.Keys.Reload,
		k.Keys.Quit,
	}
}

// FullHelp returns the full help message
func (k ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Keys.Up, k.Keys.Down, k.Keys.Enter},
		{k.Keys.Aggr, k.Keys.Reload, k.Keys.Quit},
	}
}

// ReportKeyMap implements help.KeyMap for report view
type ReportKeyMap struct {
	Keys KeyMap
}

// ShortHelp returns the short help message
func (k ReportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Keys.Up,
		k.Keys.Down,
		k.Keys.Back,
		k.Keys.Quit,
	}
}

// FullHelp returns the full help message
func (k ReportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Keys.Up, k.Keys.Down, k.Keys.Back, k.Keys.Quit},
	}
}

// AggregatedKeyMap implements help.KeyMap for aggregated view
type AggregatedKeyMap struct {
	Keys KeyMap
}

// ShortHelp returns the short help message
func (k AggregatedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Keys.Up,
		k.Keys.Down,
		k.Keys.Back,
		k.Keys.Quit,
	}
}

// FullHelp returns the full help message
func (k AggregatedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Keys.Up, k.Keys.Down, k.Keys.Back, k.Keys.Quit},
	}
}

// GetKeyMapForView returns the appropriate KeyMap for the current view
func GetKeyMapForView(
	showReport, showAggregated bool,
	keys KeyMap,
) help.KeyMap {
	if showReport {
		return ReportKeyMap{Keys: keys}
	} else if showAggregated {
		return AggregatedKeyMap{Keys: keys}
	}
	return ListKeyMap{Keys: keys}
}

// RenderHelp renders a simplified help view
func RenderHelp(showReport, showAggregated bool, keys KeyMap) string {
	if showReport || showAggregated {
		return "↑/k up · ↓/j down · esc back · q quit"
	}
	return "↑/k up · ↓/j down · enter select · a aggregated · r reload · q quit"
}
