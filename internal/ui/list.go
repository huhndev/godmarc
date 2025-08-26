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

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/huhndev/godmarc/internal/model"
)

// ReportItem is a list item for the list model
type ReportItem struct {
	Report model.DMARCReport
}

// Title returns the title for the item
func (r ReportItem) Title() string {
	return fmt.Sprintf(
		"%s - %s",
		r.Report.ReportMetadata.ReportID,
		r.Report.ReportMetadata.DateRange.Begin.Format("2006-01-02"),
	)
}

// Description returns the description for the item
func (r ReportItem) Description() string {
	return fmt.Sprintf(
		"From: %s, Records: %d",
		r.Report.ReportMetadata.OrgName,
		len(r.Report.Records),
	)
}

// FilterValue returns the value used for filtering
func (r ReportItem) FilterValue() string {
	return r.Title()
}

// CreateReportListItems creates list items from DMARC reports
func CreateReportListItems(reports []model.DMARCReport) []list.Item {
	items := make([]list.Item, len(reports))
	for i, report := range reports {
		items[i] = ReportItem{
			Report: report,
		}
	}
	return items
}

// SetupReportList initializes a list with DMARC reports
func SetupReportList(
	reports []model.DMARCReport,
	width, height int,
) list.Model {
	// Create the delegate with custom styling
	listDelegate := list.NewDefaultDelegate()

	// Style the list items
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#25A065")).
		Bold(true)

	listDelegate.Styles.SelectedTitle = selectedStyle
	listDelegate.Styles.SelectedDesc = selectedStyle

	// Adjust the item styles to use more width
	listDelegate.Styles.NormalTitle = listDelegate.Styles.NormalTitle.MaxWidth(
		150,
	)
	listDelegate.Styles.NormalDesc = listDelegate.Styles.NormalDesc.MaxWidth(
		150,
	)
	listDelegate.Styles.SelectedTitle = listDelegate.Styles.SelectedTitle.MaxWidth(
		150,
	)
	listDelegate.Styles.SelectedDesc = listDelegate.Styles.SelectedDesc.MaxWidth(
		150,
	)

	// Fix: Set a more appropriate spacing to avoid skipping items
	listDelegate.SetSpacing(1) // Reduced from default of 2

	// Create list items from reports
	items := CreateReportListItems(reports)

	// Create the list
	l := list.New(items, listDelegate, width, height-6)
	l.Title = "DMARC Reports"
	l.SetShowHelp(false) // We'll handle the help display ourselves

	// Set additional help keys
	l.AdditionalShortHelpKeys = func() []key.Binding {
		keys := DefaultKeyMap()
		return []key.Binding{
			keys.Aggr,
			keys.Reload,
			keys.Quit,
		}
	}

	return l
}
