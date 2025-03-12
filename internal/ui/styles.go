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
	"github.com/charmbracelet/lipgloss"
)

var (
	// AppStyle is the main application container style
	AppStyle = lipgloss.NewStyle().
			Padding(1, 0).
			Border(lipgloss.HiddenBorder())

	// TitleStyle is the style for the application title
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Width(100). // Will be adjusted dynamically
			Padding(0, 1).
			Bold(true)

	// StatusBarStyle is the style for the status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#666666")).
			Width(100). // Will be adjusted dynamically
			Padding(0, 1)

	// HelpStyle is used for the help bar
	HelpStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#202020")).
			Foreground(lipgloss.Color("#DDDDDD"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0000")).
			Padding(0, 1).
			Width(100).
			Bold(true)
)

// UpdateStyles updates the styles based on the window width
func UpdateStyles(width int) {
	TitleStyle = TitleStyle.Width(width)
	StatusBarStyle = StatusBarStyle.Width(width)
	HelpStyle = HelpStyle.Width(width)
}
