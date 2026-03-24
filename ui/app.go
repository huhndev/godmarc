package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/huhndev/godmarc/formatter"
	"github.com/huhndev/godmarc/model"
	"github.com/huhndev/godmarc/storage"
)

const (
	tabReports    = 0
	tabAggregated = 1
	tabFailed     = 2
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
	activeTab      int
	width          int
	height         int
	loader         *storage.ReportLoader
	errorMsg       string
	showError      bool
	errorTimeout   time.Time
	searchInput    textinput.Model
	searching      bool
	searchFilter   string
	filteredItems  []list.Item
	allItems       []list.Item
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

// NewModel creates a new application model
func NewModel() (Model, error) {
	loader, err := storage.NewReportLoader()
	if err != nil {
		return Model{}, fmt.Errorf(
			"failed to initialize report loader: %w",
			err,
		)
	}

	reports, err := loader.LoadReports()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load reports: %w", err)
	}

	storage.SortReportsByDate(reports)

	keys := DefaultKeyMap()

	h := help.New()
	h.ShowAll = true

	vp := viewport.New(0, 0)

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "Search reports..."
	ti.CharLimit = 100

	items := CreateReportListItems(reports)

	m := Model{
		reports:     reports,
		aggregated:  model.AggregateReports(reports),
		list:        l,
		viewport:    vp,
		help:        h,
		keys:        keys,
		loader:      loader,
		searchInput: ti,
		allItems:    items,
	}

	m.list.SetItems(items)

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

	// Check error timeout
	if m.showError {
		if time.Now().After(m.errorTimeout) {
			m.clearErrorMessage()
		} else {
			cmds = append(cmds, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
				return checkErrorTimeoutMsg{}
			}))
		}
	}

	switch msg := msg.(type) {
	case checkErrorTimeoutMsg:
		if m.showError && time.Now().After(m.errorTimeout) {
			m.clearErrorMessage()
		}
		return m, nil

	case errorMsg:
		m.showErrorMessage(msg.error.Error(), 5*time.Second)
		return m, nil

	case tea.WindowSizeMsg:
		m, cmd = m.handleWindowResize(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		// Handle search input mode
		if m.searching {
			return m.handleSearchKeys(msg)
		}

		// Global key handlers
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

		// View-specific key handlers
		if m.showReport {
			m, cmd = m.handleReportViewKeys(msg)
		} else {
			// Tab switching (available in all tab views)
			switch {
			case key.Matches(msg, m.keys.Tab1):
				m.activeTab = tabReports
				m.refreshTabContent()
				return m, nil
			case key.Matches(msg, m.keys.Tab2):
				m.activeTab = tabAggregated
				m.refreshTabContent()
				return m, nil
			case key.Matches(msg, m.keys.Tab3):
				m.activeTab = tabFailed
				m.refreshTabContent()
				return m, nil
			case key.Matches(msg, m.keys.Search):
				m.searching = true
				m.searchInput.Focus()
				return m, textinput.Blink
			}

			if m.activeTab == tabReports {
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
			} else {
				// Viewport-based tabs (aggregated, failed)
				m.viewport, cmd = m.viewport.Update(msg)
			}
		}
		cmds = append(cmds, cmd)
	}

	// Pass through messages to the active component
	if m.showReport || m.activeTab != tabReports {
		m.viewport, cmd = m.viewport.Update(msg)
	} else {
		if _, ok := msg.(tea.KeyMsg); !ok {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleSearchKeys handles key events during search mode
func (m Model) handleSearchKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.searching = false
		m.searchFilter = ""
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.list.SetItems(m.allItems)
		return m, nil

	case tea.KeyEnter:
		m.searching = false
		m.searchFilter = m.searchInput.Value()
		m.searchInput.Blur()
		m.applySearchFilter()
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	// Live filter as user types
	m.searchFilter = m.searchInput.Value()
	m.applySearchFilter()

	return m, cmd
}

// applySearchFilter filters the report list based on the search term
func (m *Model) applySearchFilter() {
	if m.searchFilter == "" {
		m.list.SetItems(m.allItems)
		return
	}

	filter := strings.ToLower(m.searchFilter)
	filtered := make([]list.Item, 0)
	for _, item := range m.allItems {
		if ri, ok := item.(ReportItem); ok {
			searchable := strings.ToLower(ri.FilterValue())
			if strings.Contains(searchable, filter) {
				filtered = append(filtered, item)
			}
		}
	}
	m.list.SetItems(filtered)
}

// handleWindowResize handles window resize events
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	UpdateStyles(msg.Width)

	contentHeight := msg.Height - 8 // title + tabs + help + status

	m.list.SetWidth(msg.Width)
	m.list.SetHeight(contentHeight)

	m.viewport.Width = msg.Width
	m.viewport.Height = contentHeight

	m.help.Width = msg.Width

	m.list.Styles.Title = m.list.Styles.Title.Width(msg.Width)

	m.refreshTabContent()

	return m, nil
}

// refreshTabContent updates the viewport content for the current tab
func (m *Model) refreshTabContent() {
	if m.width == 0 {
		return
	}

	contentHeight := m.height - 8

	switch {
	case m.showReport && m.selectedReport >= 0 && m.selectedReport < len(m.reports):
		m.viewport = viewport.New(m.width, contentHeight)
		m.viewport.SetContent(formatter.FormatReport(m.reports[m.selectedReport], m.width))
		m.viewport.GotoTop()
	case m.activeTab == tabAggregated:
		m.viewport = viewport.New(m.width, contentHeight)
		m.viewport.SetContent(formatter.FormatAggregatedReport(m.aggregated, m.width))
		m.viewport.GotoTop()
	case m.activeTab == tabFailed:
		m.viewport = viewport.New(m.width, contentHeight)
		m.viewport.SetContent(formatter.FormatFailedRecords(m.aggregated, m.width))
		m.viewport.GotoTop()
	}
}

// handleListViewKeys handles key events in list view
func (m Model) handleListViewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Reload):
		return m.reloadReports()
	case key.Matches(msg, m.keys.Enter):
		return m.showSelectedReport()
	}
	return m, nil
}

// handleReportViewKeys handles key events in report view
func (m Model) handleReportViewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.showReport = false
		return m, func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// reloadReports reloads reports from disk
func (m Model) reloadReports() (Model, tea.Cmd) {
	reports, err := m.loader.LoadReports()
	if err != nil {
		return m, func() tea.Msg {
			return errorMsg{err}
		}
	}

	storage.SortReportsByDate(reports)

	m.reports = reports
	m.aggregated = model.AggregateReports(reports)
	items := CreateReportListItems(reports)
	m.allItems = items
	m.list.SetItems(items)

	// Re-apply search filter if active
	if m.searchFilter != "" {
		m.applySearchFilter()
	}

	m.showErrorMessage(
		fmt.Sprintf("Loaded %d reports", len(reports)),
		3*time.Second,
	)

	m.refreshTabContent()

	return m, nil
}

// showSelectedReport switches to the report detail view
func (m Model) showSelectedReport() (Model, tea.Cmd) {
	m.selectedReport = m.list.Index()
	if m.selectedReport >= 0 && m.selectedReport < len(m.reports) {
		m.showReport = true
		m.refreshTabContent()
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	title := TitleStyle.Render("godmarc - DMARC Report Analyzer")
	tabBar := m.renderTabBar()
	statusBar := m.renderStatusBar()
	helpView := HelpStyle.Render(
		RenderHelp(m.showReport, m.activeTab, m.searching, m.keys),
	)

	var content string
	if m.showReport {
		content = m.viewport.View()
	} else if m.activeTab == tabReports {
		content = m.list.View()
	} else {
		content = m.viewport.View()
	}

	// Search bar
	var searchBar string
	if m.searching {
		searchBar = SearchStyle.Width(m.width).Render("/ " + m.searchInput.View())
	}

	parts := []string{title, tabBar}
	if m.showError {
		parts = append(parts, ErrorStyle.Render(m.errorMsg))
	}
	if searchBar != "" {
		parts = append(parts, searchBar)
	}
	parts = append(parts, content, helpView, statusBar)

	return AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

// renderTabBar renders the tab navigation bar
func (m Model) renderTabBar() string {
	if m.showReport {
		return ActiveTabStyle.Width(m.width).Render("Report Detail")
	}

	tabs := []struct {
		label  string
		active bool
	}{
		{"Reports", m.activeTab == tabReports},
		{"Aggregated", m.activeTab == tabAggregated},
		{"Failed", m.activeTab == tabFailed},
	}

	rendered := make([]string, len(tabs))
	for i, tab := range tabs {
		label := fmt.Sprintf(" %d %s ", i+1, tab.label)
		if tab.active {
			rendered[i] = ActiveTabStyle.Render(label)
		} else {
			rendered[i] = TabStyle.Render(label)
		}
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)

	// Fill remaining width with background
	tabWidth := lipgloss.Width(tabRow)
	if tabWidth < m.width {
		filler := TabBarStyle.Width(m.width - tabWidth).Render("")
		tabRow = lipgloss.JoinHorizontal(lipgloss.Top, tabRow, filler)
	}

	return tabRow
}

// renderStatusBar renders the status bar with useful info
func (m Model) renderStatusBar() string {
	var left string
	var right string

	if m.showReport {
		reportName := m.reports[m.selectedReport].ReportMetadata.ReportID
		left = fmt.Sprintf(" Report: %s", reportName)
		right = fmt.Sprintf("scroll: %.0f%% ", m.viewport.ScrollPercent()*100)
	} else {
		totalReports := len(m.reports)
		totalRecords := m.aggregated.TotalRecords
		failedCount := len(m.aggregated.FailedRecords)

		left = fmt.Sprintf(" %d reports | %d records", totalReports, totalRecords)

		if failedCount > 0 {
			right = FailStyle.Render(fmt.Sprintf("%d failed ", failedCount))
		} else {
			right = PassStyle.Render("all pass ")
		}

		if m.searchFilter != "" {
			left += fmt.Sprintf(" | filter: %q", m.searchFilter)
		}

		if m.aggregated.TotalReports > 0 {
			dateRange := fmt.Sprintf("%s to %s ",
				m.aggregated.DateRange.Begin.Format("2006-01-02"),
				m.aggregated.DateRange.End.Format("2006-01-02"))
			right = dateRange + right
		}
	}

	// Build status bar with left and right alignment
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := m.width - leftWidth - rightWidth
	if gap < 0 {
		gap = 0
	}

	return StatusBarStyle.Render(left + strings.Repeat(" ", gap) + right)
}

// Custom message types
type checkErrorTimeoutMsg struct{}

type errorMsg struct {
	error error
}
