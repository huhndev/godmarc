package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/huhndev/godmarc/model"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	passStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#25A065")).
			Bold(true)

	failStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4040")).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD")).
			Width(20)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5"))
)

func colorResult(result string) string {
	switch result {
	case "pass":
		return passStyle.Render(result)
	case "fail":
		return failStyle.Render(result)
	case "softfail", "neutral", "temperror", "permerror":
		return warnStyle.Render(result)
	default:
		return result
	}
}

func colorDisposition(disp string) string {
	switch disp {
	case "none":
		return passStyle.Render(disp)
	case "reject":
		return failStyle.Render(disp)
	case "quarantine":
		return warnStyle.Render(disp)
	default:
		return disp
	}
}

// FormatReport formats a single DMARC report for display
func FormatReport(report model.DMARCReport, width int) string {
	var sb strings.Builder

	if width < 60 {
		width = 60
	}

	// Report metadata
	sb.WriteString(headerStyle.Render("Report Metadata") + "\n\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Organization:"), valueStyle.Render(report.ReportMetadata.OrgName)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Email:"), valueStyle.Render(report.ReportMetadata.Email)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Report ID:"), valueStyle.Render(report.ReportMetadata.ReportID)))
	sb.WriteString(fmt.Sprintf("  %s %s to %s\n",
		labelStyle.Render("Date Range:"),
		valueStyle.Render(report.ReportMetadata.DateRange.Begin.Format("2006-01-02")),
		valueStyle.Render(report.ReportMetadata.DateRange.End.Format("2006-01-02"))))

	// Policy published
	sb.WriteString("\n" + headerStyle.Render("Published Policy") + "\n\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Domain:"), valueStyle.Render(report.PolicyPublished.Domain)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("DKIM Alignment:"), valueStyle.Render(report.PolicyPublished.ADKIM)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("SPF Alignment:"), valueStyle.Render(report.PolicyPublished.ASPF)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Policy:"), valueStyle.Render(report.PolicyPublished.P)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Subdomain Policy:"), valueStyle.Render(report.PolicyPublished.SP)))
	sb.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Percentage:"), valueStyle.Render(fmt.Sprintf("%d%%", report.PolicyPublished.PCT))))

	// Records table
	sb.WriteString("\n" + headerStyle.Render(fmt.Sprintf("Records (%d)", len(report.Records))) + "\n\n")

	tableWidth := width - 4
	if tableWidth < 60 {
		tableWidth = 60
	}

	rows := make([][]string, 0, len(report.Records))
	for _, record := range report.Records {
		rows = append(rows, []string{
			record.Row.SourceIP,
			fmt.Sprintf("%d", record.Row.Count),
			colorDisposition(record.Row.PolicyEvaluated.Disposition),
			colorResult(record.Row.PolicyEvaluated.DKIM),
			colorResult(record.Row.PolicyEvaluated.SPF),
			record.Identifiers.HeaderFrom,
		})
	}

	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
		Headers("Source IP", "Count", "Disposition", "DKIM", "SPF", "Header From").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})

	sb.WriteString(t.Render() + "\n")

	// Auth details per record
	for i, record := range report.Records {
		sb.WriteString("\n" + headerStyle.Render(fmt.Sprintf("Record #%d Auth Details", i+1)) + "\n\n")

		// DKIM authentication
		if len(record.AuthResults.DKIM) > 0 {
			dkimRows := make([][]string, 0, len(record.AuthResults.DKIM))
			for _, dkim := range record.AuthResults.DKIM {
				dkimRows = append(dkimRows, []string{
					dkim.Domain,
					colorResult(dkim.Result),
					dkim.Selector,
				})
			}

			dt := ltable.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
				Headers("DKIM Domain", "Result", "Selector").
				Rows(dkimRows...).
				StyleFunc(func(row, col int) lipgloss.Style {
					if row == ltable.HeaderRow {
						return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
					}
					return lipgloss.NewStyle().Padding(0, 1)
				})

			sb.WriteString(dt.Render() + "\n")
		}

		// SPF authentication
		if len(record.AuthResults.SPF) > 0 {
			spfRows := make([][]string, 0, len(record.AuthResults.SPF))
			for _, spf := range record.AuthResults.SPF {
				spfRows = append(spfRows, []string{
					spf.Domain,
					colorResult(spf.Result),
					spf.Scope,
				})
			}

			st := ltable.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))).
				Headers("SPF Domain", "Result", "Scope").
				Rows(spfRows...).
				StyleFunc(func(row, col int) lipgloss.Style {
					if row == ltable.HeaderRow {
						return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87")).Padding(0, 1)
					}
					return lipgloss.NewStyle().Padding(0, 1)
				})

			sb.WriteString(st.Render() + "\n")
		}
	}

	return sb.String()
}

// TruncateString truncates a string to maxLen, adding ellipsis if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
