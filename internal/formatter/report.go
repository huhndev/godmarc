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

package formatter

import (
	"fmt"
	"strings"

	"github.com/huhndev/godmarc/internal/model"
)

// FormatReport formats a single DMARC report for display with improved width usage
func FormatReport(report model.DMARCReport) string {
	var sb strings.Builder

	// Report metadata
	sb.WriteString(ReportInfoStyle.Render("Report Metadata") + "\n")
	sb.WriteString(
		fmt.Sprintf("  Organization:    %s\n", report.ReportMetadata.OrgName),
	)
	sb.WriteString(
		fmt.Sprintf("  Email:           %s\n", report.ReportMetadata.Email),
	)
	sb.WriteString(
		fmt.Sprintf("  Report ID:       %s\n", report.ReportMetadata.ReportID),
	)
	sb.WriteString(fmt.Sprintf("  Date Range:      %s to %s\n",
		report.ReportMetadata.DateRange.Begin.Format("2006-01-02"),
		report.ReportMetadata.DateRange.End.Format("2006-01-02")))

	// Policy published
	sb.WriteString("\n" + ReportInfoStyle.Render("Published Policy") + "\n")
	sb.WriteString(
		fmt.Sprintf("  Domain:          %s\n", report.PolicyPublished.Domain),
	)
	sb.WriteString(
		fmt.Sprintf("  DKIM Alignment:  %s\n", report.PolicyPublished.ADKIM),
	)
	sb.WriteString(
		fmt.Sprintf("  SPF Alignment:   %s\n", report.PolicyPublished.ASPF),
	)
	sb.WriteString(
		fmt.Sprintf("  Policy:          %s\n", report.PolicyPublished.P),
	)
	sb.WriteString(
		fmt.Sprintf("  Subdomain Policy: %s\n", report.PolicyPublished.SP),
	)
	sb.WriteString(
		fmt.Sprintf("  Percentage:      %d%%\n", report.PolicyPublished.PCT),
	)

	// Records - using full width for displaying record data
	sb.WriteString(
		"\n" + ReportInfoStyle.Render(
			fmt.Sprintf("Records (%d)", len(report.Records)),
		) + "\n",
	)

	for i, record := range report.Records {
		recordHeader := fmt.Sprintf("Record #%d", i+1)
		sb.WriteString("\n" + SelectedItemStyle.Render(recordHeader) + "\n")

		// Display more details horizontally to use more width
		sb.WriteString(
			fmt.Sprintf("  Source IP: %-40s Count: %-10d Disposition: %-15s\n",
				record.Row.SourceIP,
				record.Row.Count,
				record.Row.PolicyEvaluated.Disposition),
		)

		sb.WriteString(
			fmt.Sprintf(
				"  DKIM Result: %-15s SPF Result: %-15s Header From: %s\n",
				record.Row.PolicyEvaluated.DKIM,
				record.Row.PolicyEvaluated.SPF,
				record.Identifiers.HeaderFrom,
			),
		)

		// DKIM authentication
		sb.WriteString("\n  DKIM Authentication:\n")
		if len(record.AuthResults.DKIM) == 0 {
			sb.WriteString("    None\n")
		} else {
			// Create a more tabular format for DKIM results
			for _, dkim := range record.AuthResults.DKIM {
				sb.WriteString(fmt.Sprintf("    Domain: %-30s Result: %-10s Selector: %s\n",
					dkim.Domain, dkim.Result, dkim.Selector))
			}
		}

		// SPF authentication
		sb.WriteString("\n  SPF Authentication:\n")
		if len(record.AuthResults.SPF) == 0 {
			sb.WriteString("    None\n")
		} else {
			// Create a more tabular format for SPF results
			for _, spf := range record.AuthResults.SPF {
				sb.WriteString(fmt.Sprintf("    Domain: %-30s Result: %-10s Scope: %s\n",
					spf.Domain, spf.Result, spf.Scope))
			}
		}
	}

	return sb.String()
}

// Helper function to truncate strings for display
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper function to get the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to get the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
