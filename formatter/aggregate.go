package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/huhndev/godmarc/model"
)

// FormatAggregatedReport formats an aggregated report for display
func FormatAggregatedReport(aggr model.AggregatedReport, width int) string {
	var sb strings.Builder

	if width < 60 {
		width = 60
	}

	// Overview
	sb.WriteString(headerStyle.Render("Aggregated Report") + "\n\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Total Reports:"), valueStyle.Render(fmt.Sprintf("%d", aggr.TotalReports))))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Total Records:"), valueStyle.Render(fmt.Sprintf("%d", aggr.TotalRecords))))
	sb.WriteString(fmt.Sprintf("  %s %s to %s\n",
		labelStyle.Render("Date Range:"),
		valueStyle.Render(aggr.DateRange.Begin.Format("2006-01-02")),
		valueStyle.Render(aggr.DateRange.End.Format("2006-01-02"))))

	// Domains table
	sb.WriteString("\n" + headerStyle.Render("Domains") + "\n\n")
	domainRows := make([][]string, 0, len(aggr.Domains))
	for domain, count := range aggr.Domains {
		domainRows = append(domainRows, []string{domain, fmt.Sprintf("%d", count)})
	}
	sort.Slice(domainRows, func(i, j int) bool { return domainRows[i][0] < domainRows[j][0] })

	dt := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
		Headers("Domain", "Reports").
		Rows(domainRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})

	sb.WriteString(dt.Render() + "\n")

	// Summary statistics - side by side using lipgloss columns
	sb.WriteString("\n" + headerStyle.Render("Summary Statistics") + "\n\n")

	// Build three columns
	dispCol := buildStatColumn("Dispositions", aggr.Dispositions, true)
	dkimCol := buildStatColumn("DKIM Results", aggr.DKIMResults, false)
	spfCol := buildStatColumn("SPF Results", aggr.SPFResults, false)

	colWidth := (width - 10) / 3
	if colWidth < 20 {
		colWidth = 20
	}

	colStyle := lipgloss.NewStyle().Width(colWidth)
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		colStyle.Render(dispCol),
		colStyle.Render(dkimCol),
		colStyle.Render(spfCol),
	) + "\n")

	// Top Sources table
	sb.WriteString("\n" + headerStyle.Render("Top Sources") + "\n\n")

	type sourceCount struct {
		source string
		count  int
	}
	sourceCounts := make([]sourceCount, 0, len(aggr.Sources))
	for source, count := range aggr.Sources {
		sourceCounts = append(sourceCounts, sourceCount{source, count})
	}
	sort.Slice(sourceCounts, func(i, j int) bool {
		return sourceCounts[i].count > sourceCounts[j].count
	})

	topCount := len(sourceCounts)
	if topCount > 20 {
		topCount = 20
	}

	sourceRows := make([][]string, 0, topCount)
	for i := 0; i < topCount; i++ {
		sourceRows = append(sourceRows, []string{
			sourceCounts[i].source,
			fmt.Sprintf("%d", sourceCounts[i].count),
		})
	}

	st := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
		Headers("Source IP", "Count").
		Rows(sourceRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})

	sb.WriteString(st.Render() + "\n")

	if len(sourceCounts) > 20 {
		sb.WriteString(fmt.Sprintf("\n  ... and %d more sources\n", len(sourceCounts)-20))
	}

	return sb.String()
}

// FormatFailedRecords formats only failed records for the "Failed" tab
func FormatFailedRecords(aggr model.AggregatedReport, width int) string {
	var sb strings.Builder

	if width < 60 {
		width = 60
	}

	sb.WriteString(headerStyle.Render(fmt.Sprintf("Failed Records (%d)", len(aggr.FailedRecords))) + "\n\n")

	if len(aggr.FailedRecords) == 0 {
		sb.WriteString(passStyle.Render("  No failed records!") + "\n")
		return sb.String()
	}

	maxRows := len(aggr.FailedRecords)
	if maxRows > 50 {
		maxRows = 50
	}

	rows := make([][]string, 0, maxRows)
	for i := 0; i < maxRows; i++ {
		record := aggr.FailedRecords[i]
		rows = append(rows, []string{
			record.SourceIP,
			record.Domain,
			fmt.Sprintf("%d", record.Count),
			record.Reason,
		})
	}

	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
		Headers("Source IP", "Domain", "Count", "Reason").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
			}
			s := lipgloss.NewStyle().Padding(0, 1)
			if col == 3 && row >= 0 {
				return s.Foreground(lipgloss.Color("#FF4040"))
			}
			return s
		})

	sb.WriteString(t.Render() + "\n")

	if len(aggr.FailedRecords) > 50 {
		sb.WriteString(fmt.Sprintf("\n  ... and %d more failed records\n", len(aggr.FailedRecords)-50))
	}

	return sb.String()
}

func buildStatColumn(title string, data map[string]int, isDisposition bool) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render(title) + "\n\n")

	entries := make([]struct {
		key   string
		count int
	}, 0, len(data))
	for k, v := range data {
		entries = append(entries, struct {
			key   string
			count int
		}{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	for _, e := range entries {
		var colored string
		if isDisposition {
			colored = colorDisposition(e.key)
		} else {
			colored = colorResult(e.key)
		}
		sb.WriteString(fmt.Sprintf("  %-15s %d\n", colored, e.count))
	}

	return sb.String()
}
