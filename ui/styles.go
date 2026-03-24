package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
const (
	ColorGreen   = lipgloss.Color("#25A065")
	ColorRed     = lipgloss.Color("#FF4040")
	ColorYellow  = lipgloss.Color("#FFD700")
	ColorPink    = lipgloss.Color("#FF5F87")
	ColorWhite   = lipgloss.Color("#FFFDF5")
	ColorGray    = lipgloss.Color("#666666")
	ColorDimGray = lipgloss.Color("#444444")
	ColorDarkBg  = lipgloss.Color("#202020")
	ColorLightFg = lipgloss.Color("#DDDDDD")
)

var (
	// AppStyle is the main application container style
	AppStyle = lipgloss.NewStyle().
			Padding(1, 0).
			Border(lipgloss.HiddenBorder())

	// TitleStyle is the style for the application title
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorGreen).
			Width(100).
			Padding(0, 1).
			Bold(true)

	// TabStyle is the style for inactive tabs
	TabStyle = lipgloss.NewStyle().
			Foreground(ColorLightFg).
			Background(ColorDimGray).
			Padding(0, 2)

	// ActiveTabStyle is the style for the active tab
	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorGreen).
			Padding(0, 2).
			Bold(true)

	// TabBarStyle is the container for the tab bar
	TabBarStyle = lipgloss.NewStyle().
			Background(ColorDimGray)

	// StatusBarStyle is the style for the status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorGray).
			Width(100).
			Padding(0, 1)

	// HelpStyle is used for the help bar
	HelpStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(ColorDarkBg).
			Foreground(ColorLightFg)

	// ErrorStyle is for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorRed).
			Padding(0, 1).
			Width(100).
			Bold(true)

	// SectionHeaderStyle is used for report section headers
	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorPink).
				Bold(true)

	// SelectedItemStyle is used for highlighting selected items
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true)

	// PassStyle renders pass results in green
	PassStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	// FailStyle renders fail results in red
	FailStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	// WarnStyle renders soft-fail/neutral results in yellow
	WarnStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	// TableHeaderStyle is for table headers
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPink).
				Padding(0, 1)

	// TableCellStyle is for table cells
	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// SearchStyle is for the search input
	SearchStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorDimGray).
			Padding(0, 1)
)

// ColorizeResult returns a styled string for auth results (pass/fail/softfail)
func ColorizeResult(result string) string {
	switch result {
	case "pass":
		return PassStyle.Render(result)
	case "fail":
		return FailStyle.Render(result)
	case "softfail", "neutral", "temperror", "permerror":
		return WarnStyle.Render(result)
	default:
		return result
	}
}

// ColorizeDisposition returns a styled string for dispositions
func ColorizeDisposition(disp string) string {
	switch disp {
	case "none":
		return PassStyle.Render(disp)
	case "reject":
		return FailStyle.Render(disp)
	case "quarantine":
		return WarnStyle.Render(disp)
	default:
		return disp
	}
}

// UpdateStyles updates the styles based on the window width
func UpdateStyles(width int) {
	TitleStyle = TitleStyle.Width(width)
	StatusBarStyle = StatusBarStyle.Width(width)
	HelpStyle = HelpStyle.Width(width)
	TabBarStyle = TabBarStyle.Width(width)
	ErrorStyle = ErrorStyle.Width(width)
}
