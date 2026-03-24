package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/huhndev/godmarc/model"
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
	// Show pass/fail summary
	dkimPass := 0
	spfPass := 0
	total := len(r.Report.Records)
	for _, rec := range r.Report.Records {
		if rec.Row.PolicyEvaluated.DKIM == "pass" {
			dkimPass++
		}
		if rec.Row.PolicyEvaluated.SPF == "pass" {
			spfPass++
		}
	}

	dkimStatus := PassStyle.Render(fmt.Sprintf("DKIM %d/%d", dkimPass, total))
	spfStatus := PassStyle.Render(fmt.Sprintf("SPF %d/%d", spfPass, total))
	if dkimPass < total {
		dkimStatus = FailStyle.Render(fmt.Sprintf("DKIM %d/%d", dkimPass, total))
	}
	if spfPass < total {
		spfStatus = FailStyle.Render(fmt.Sprintf("SPF %d/%d", spfPass, total))
	}

	return fmt.Sprintf(
		"%s | %d records | %s %s",
		r.Report.ReportMetadata.OrgName,
		total,
		dkimStatus,
		spfStatus,
	)
}

// FilterValue returns the value used for filtering
func (r ReportItem) FilterValue() string {
	return strings.Join([]string{
		r.Report.ReportMetadata.ReportID,
		r.Report.ReportMetadata.OrgName,
		r.Report.PolicyPublished.Domain,
		r.Report.ReportMetadata.DateRange.Begin.Format("2006-01-02"),
	}, " ")
}

// CreateReportListItems creates list items from DMARC reports
func CreateReportListItems(reports []model.DMARCReport) []list.Item {
	items := make([]list.Item, len(reports))
	for i, report := range reports {
		items[i] = ReportItem{Report: report}
	}
	return items
}

// SetupReportList initializes a list with DMARC reports
func SetupReportList(
	reports []model.DMARCReport,
	width, height int,
) list.Model {
	listDelegate := list.NewDefaultDelegate()

	selectedStyle := lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)

	listDelegate.Styles.SelectedTitle = selectedStyle
	listDelegate.Styles.SelectedDesc = selectedStyle

	listDelegate.Styles.NormalTitle = listDelegate.Styles.NormalTitle.MaxWidth(150)
	listDelegate.Styles.NormalDesc = listDelegate.Styles.NormalDesc.MaxWidth(150)
	listDelegate.Styles.SelectedTitle = listDelegate.Styles.SelectedTitle.MaxWidth(150)
	listDelegate.Styles.SelectedDesc = listDelegate.Styles.SelectedDesc.MaxWidth(150)

	listDelegate.SetSpacing(1)

	items := CreateReportListItems(reports)

	l := list.New(items, listDelegate, width, height-6)
	l.Title = "DMARC Reports"
	l.SetShowHelp(false)

	return l
}
