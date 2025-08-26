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
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/huhndev/godmarc/internal/formatter"
	"github.com/huhndev/godmarc/internal/model"
	"github.com/huhndev/godmarc/internal/storage"
)

// Model represents the state of the application
type Model struct {
	reports        []model.DMARCReport
	aggregated     model.AggregatedReport
	list           list.Model
	viewport       viewport.Model
	help           help.Model
	keys           KeyMap
	selectedReport int
	showReport     bool
	showAggregated bool
	width          int
	height         int
	loader         *storage.ReportLoader
	errorMsg       string
	showError      bool
	errorTimeout   time.Time
}

// showErrorMessage displays an error message for a specified duration
func (m *Model) showErrorMessage(msg string, duration time.Duration) {
	m.errorMsg = msg
	m.showError = true
	m.errorTimeout = time.Now().Add(duration)
}

// clearErrorMessage clears the current error message
func (m *Model) clearErrorMessage() {
	m.errorMsg = ""
	m.showError = false
}

// checkErrorTimeout checks if the error message should be cleared
func (m *Model) checkErrorTimeout() tea.Cmd {
	if m.showError && time.Now().After(m.errorTimeout) {
		m.clearErrorMessage()
	}
	return nil
}

// NewModel creates a new application model
func NewModel() (Model, error) {
	// Setup report loader
	loader, err := storage.NewReportLoader()
	if err != nil {
		return Model{}, fmt.Errorf(
			"failed to initialize report loader: %w",
			err,
		)
	}

	// Load initial reports
	reports, err := loader.LoadReports()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load reports: %w", err)
	}

	// Sort reports by date
	storage.SortReportsByDate(reports)

	// Setup key bindings
	keys := DefaultKeyMap()

	// Create help model
	h := help.New()
	h.ShowAll = true

	// Setup viewport for report display
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		MarginLeft(0).
		MarginRight(0).
		PaddingLeft(0).
		PaddingRight(0).
		MaxWidth(9999) // Effectively no max width

	// Create list (will be properly sized later in WindowSizeMsg)
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)

	// Create model
	m := Model{
		reports:    reports,
		aggregated: model.AggregateReports(reports),
		list:       l,
		viewport:   vp,
		help:       h,
		keys:       keys,
		loader:     loader,
	}

	// Initialize with proper list items
	m.list.SetItems(CreateReportListItems(reports))

	return m, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Check if we need to clear error messages
	if m.showError {
		if time.Now().After(m.errorTimeout) {
			m.clearErrorMessage()
		} else {
			// Set a command to check again after timeout
			cmds = append(cmds, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
				return checkErrorTimeoutMsg{}
			}))
		}
	}

	// Define a custom message type for checking error timeouts
	switch msg := msg.(type) {
	case checkErrorTimeoutMsg:
		return m, m.checkErrorTimeout()

	case errorMsg:
		m.showErrorMessage(msg.error.Error(), 5*time.Second)
		return m, nil

	case tea.WindowSizeMsg:
		m, cmd = m.handleWindowResize(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		// Global key handlers
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

		// View-specific key handlers
		if m.showReport || m.showAggregated {
			m, cmd = m.handleReportViewKeys(msg)
		} else {
			// FIX: Handle explicit up/down navigation for the list to avoid double-processing
			switch {
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
				return m, nil
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
				return m, nil
			default:
				m, cmd = m.handleListViewKeys(msg)
			}
		}
		cmds = append(cmds, cmd)
	}

	// Pass through messages to the active component
	if m.showReport || m.showAggregated {
		m.viewport, cmd = m.viewport.Update(msg)
	} else {
		// Only update the list if we haven't explicitly handled the key above
		if _, ok := msg.(tea.KeyMsg); !ok {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleWindowResize handles window resize events
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	// Update styles
	UpdateStyles(msg.Width)

	// Update list dimensions
	m.list.SetWidth(msg.Width)
	m.list.SetHeight(msg.Height - 6) // Reserve space for title, status, help

	// Update viewport dimensions
	m.viewport.Width = msg.Width
	m.viewport.Height = msg.Height - 6

	// Update help width
	m.help.Width = msg.Width

	// Adjust list styles
	m.list.Styles.Title = m.list.Styles.Title.Width(msg.Width)
	m.list.Styles.FilterPrompt = m.list.Styles.FilterPrompt.Width(msg.Width)

	// Update content if already viewing a report
	if m.showReport && m.selectedReport >= 0 &&
		m.selectedReport < len(m.reports) {
		// Create a fresh viewport to ensure proper rendering
		m.viewport = viewport.New(msg.Width, msg.Height-6)
		m.viewport.Style = lipgloss.NewStyle().
			MarginLeft(0).
			MarginRight(0).
			PaddingLeft(0).
			PaddingRight(0).
			MaxWidth(msg.Width)
		m.viewport.SetContent(
			formatter.FormatReport(m.reports[m.selectedReport]),
		)
		m.viewport.GotoTop()
	} else if m.showAggregated {
		// Create a fresh viewport for aggregated view
		m.viewport = viewport.New(msg.Width, msg.Height-6)
		m.viewport.Style = lipgloss.NewStyle().
			MarginLeft(0).
			MarginRight(0).
			PaddingLeft(0).
			PaddingRight(0).
			MaxWidth(msg.Width)
		m.viewport.SetContent(formatter.FormatAggregatedReport(m.aggregated))
		m.viewport.GotoTop()
	}

	return m, nil
}

// handleListViewKeys handles key events in list view
func (m Model) handleListViewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Reload):
		return m.reloadReports()

	case key.Matches(msg, m.keys.Aggr):
		return m.showAggregatedReport()

	case key.Matches(msg, m.keys.Enter):
		return m.showSelectedReport()
	}

	// We'll handle list navigation directly in the main Update function
	// to avoid double-processing, so we don't pass Up/Down to the list here
	return m, nil
}

// handleReportViewKeys handles key events in report view
func (m Model) handleReportViewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.showReport = false
		m.showAggregated = false
		// Force a resize to refresh the view properly when going back to list
		return m, func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		}
	}

	// Pass the message to the viewport
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// reloadReports reloads reports from disk
func (m Model) reloadReports() (Model, tea.Cmd) {
	reports, err := m.loader.LoadReports()
	if err != nil {
		// Surface the error to the user interface
		return m, func() tea.Msg {
			return errorMsg{err}
		}
	}

	// Sort reports by date
	storage.SortReportsByDate(reports)

	// Update model
	m.reports = reports
	m.aggregated = model.AggregateReports(reports)
	m.list.SetItems(CreateReportListItems(reports))

	// Show success message
	m.showErrorMessage(
		fmt.Sprintf("Successfully loaded %d reports", len(reports)),
		3*time.Second,
	)

	return m, nil
}

// showAggregatedReport switches to the aggregated report view
func (m Model) showAggregatedReport() (Model, tea.Cmd) {
	m.showAggregated = !m.showAggregated
	if m.showAggregated {
		m.viewport.SetContent(formatter.FormatAggregatedReport(m.aggregated))
		m.viewport.GotoTop()
	}
	return m, nil
}

// showSelectedReport switches to the report detail view
func (m Model) showSelectedReport() (Model, tea.Cmd) {
	m.selectedReport = m.list.Index()
	if m.selectedReport >= 0 && m.selectedReport < len(m.reports) {
		m.showReport = true
		// Set viewport to use full width and refresh content
		m.viewport = viewport.New(m.width, m.height-6)
		m.viewport.SetContent(
			formatter.FormatReport(m.reports[m.selectedReport]),
		)
		m.viewport.GotoTop()
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string
	title := TitleStyle.Render("godmarc - DMARC Report Analyzer")
	statusBar := m.renderStatusBar()
	helpView := HelpStyle.Render(
		RenderHelp(m.showReport, m.showAggregated, m.keys),
	)

	// Different rendering approach based on view mode
	if m.showReport || m.showAggregated {
		content = m.viewport.View()
	} else {
		content = m.list.View()
	}

	// Show error message if present
	if m.showError {
		errorView := ErrorStyle.Render(m.errorMsg)
		return AppStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			errorView,
			content,
			helpView,
			statusBar,
		))
	}

	// Normal view without error
	return AppStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		helpView,
		statusBar,
	))
}

// renderStatusBar renders the status bar
func (m Model) renderStatusBar() string {
	var status string

	if m.showReport {
		reportName := m.reports[m.selectedReport].ReportMetadata.ReportID
		status = StatusBarStyle.Render(
			fmt.Sprintf(" Report: %s ", reportName),
		)
	} else if m.showAggregated {
		status = StatusBarStyle.Render(" Aggregated Report ")
	} else {
		totalReports := len(m.reports)
		status = StatusBarStyle.Render(fmt.Sprintf(" %d Report(s) ", totalReports))
	}

	return status
}

// Custom message types for error handling
type checkErrorTimeoutMsg struct{}

type errorMsg struct {
	error error
}
