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
	"sort"
	"strings"

	"github.com/huhnsystems/godmarc/internal/model"
)

// FormatAggregatedReport formats an aggregated report for display with improved width usage
func FormatAggregatedReport(aggr model.AggregatedReport) string {
	var sb strings.Builder

	// Overview
	sb.WriteString(ReportInfoStyle.Render("Aggregated Report") + "\n")
	sb.WriteString(fmt.Sprintf("  Total Reports:    %d\n", aggr.TotalReports))
	sb.WriteString(fmt.Sprintf("  Total Records:    %d\n", aggr.TotalRecords))
	sb.WriteString(fmt.Sprintf("  Date Range:       %s to %s\n",
		aggr.DateRange.Begin.Format("2006-01-02"),
		aggr.DateRange.End.Format("2006-01-02")))

	// Domains - improve layout to use full width
	sb.WriteString("\n" + ReportInfoStyle.Render("Domains") + "\n")
	domainCount := 0
	var domainLine string
	domainEntries := make([]string, 0, len(aggr.Domains))

	// Collect domain entries
	for domain, count := range aggr.Domains {
		domainEntries = append(
			domainEntries,
			fmt.Sprintf("%-35s %d", domain, count),
		)
	}

	// Sort domain entries for consistent display
	sort.Strings(domainEntries)

	// Display domains in a multi-column format (2 columns)
	for i, entry := range domainEntries {
		if i%2 == 0 {
			domainLine = "  " + entry
		} else {
			sb.WriteString(fmt.Sprintf("%-50s %s\n", domainLine, entry))
			domainLine = ""
		}
		domainCount++
	}

	// Handle odd number of domains
	if domainLine != "" {
		sb.WriteString(domainLine + "\n")
	}

	// Show summary data in a more compact layout
	sb.WriteString("\n" + ReportInfoStyle.Render("Summary Statistics") + "\n")

	// Dispositions, DKIM, and SPF in side-by-side columns
	sb.WriteString(
		"  " + ReportInfoStyle.Render("Dispositions") + "          " +
			ReportInfoStyle.Render("DKIM Results") + "          " +
			ReportInfoStyle.Render("SPF Results") + "\n",
	)

	// Get the max length of each category for proper display
	maxDispLen := 0
	maxDkimLen := 0
	maxSpfLen := 0

	for disp := range aggr.Dispositions {
		if len(disp) > maxDispLen {
			maxDispLen = len(disp)
		}
	}

	for dkim := range aggr.DKIMResults {
		if len(dkim) > maxDkimLen {
			maxDkimLen = len(dkim)
		}
	}

	for spf := range aggr.SPFResults {
		if len(spf) > maxSpfLen {
			maxSpfLen = len(spf)
		}
	}

	// Collect and sort entries
	dispEntries := make([]string, 0, len(aggr.Dispositions))
	dkimEntries := make([]string, 0, len(aggr.DKIMResults))
	spfEntries := make([]string, 0, len(aggr.SPFResults))

	for disp, count := range aggr.Dispositions {
		dispEntries = append(
			dispEntries,
			fmt.Sprintf("%-*s %d", maxDispLen+2, disp, count),
		)
	}

	for dkim, count := range aggr.DKIMResults {
		dkimEntries = append(
			dkimEntries,
			fmt.Sprintf("%-*s %d", maxDkimLen+2, dkim, count),
		)
	}

	for spf, count := range aggr.SPFResults {
		spfEntries = append(
			spfEntries,
			fmt.Sprintf("%-*s %d", maxSpfLen+2, spf, count),
		)
	}

	sort.Strings(dispEntries)
	sort.Strings(dkimEntries)
	sort.Strings(spfEntries)

	// Display side by side
	maxRows := Max(len(dispEntries), Max(len(dkimEntries), len(spfEntries)))
	for i := 0; i < maxRows; i++ {
		disp := ""
		dkim := ""
		spf := ""

		if i < len(dispEntries) {
			disp = dispEntries[i]
		}

		if i < len(dkimEntries) {
			dkim = dkimEntries[i]
		}

		if i < len(spfEntries) {
			spf = spfEntries[i]
		}

		// Pad the columns for alignment
		sb.WriteString(
			fmt.Sprintf("  %-20s     %-20s     %-20s\n", disp, dkim, spf),
		)
	}

	// Top Sources - multi-column display
	sb.WriteString("\n" + ReportInfoStyle.Render("Top Sources") + "\n")

	// Convert map to slice for sorting
	type sourceCount struct {
		source string
		count  int
	}

	sourceCounts := make([]sourceCount, 0, len(aggr.Sources))
	for source, count := range aggr.Sources {
		sourceCounts = append(sourceCounts, sourceCount{source, count})
	}

	// Sort by count (descending)
	sort.Slice(sourceCounts, func(i, j int) bool {
		return sourceCounts[i].count > sourceCounts[j].count
	})

	// Display top sources in 2 columns
	topCount := Min(20, len(sourceCounts))
	for i := 0; i < topCount; i += 2 {
		left := fmt.Sprintf(
			"%-20s %d",
			sourceCounts[i].source,
			sourceCounts[i].count,
		)
		right := ""
		if i+1 < topCount {
			right = fmt.Sprintf(
				"%-20s %d",
				sourceCounts[i+1].source,
				sourceCounts[i+1].count,
			)
		}
		if right != "" {
			sb.WriteString(fmt.Sprintf("  %-35s %s\n", left, right))
		} else {
			sb.WriteString("  " + left + "\n")
		}
	}

	// Failed Records - more compact display with columns for key info
	sb.WriteString(
		"\n" + ReportInfoStyle.Render(
			fmt.Sprintf("Failed Records (%d)", len(aggr.FailedRecords)),
		) + "\n",
	)
	sb.WriteString(
		fmt.Sprintf(
			"  %-20s %-30s %-8s %s\n",
			"Source IP",
			"Domain",
			"Count",
			"Reason",
		),
	)
	sb.WriteString("  " + strings.Repeat("-", 78) + "\n")

	for i, record := range aggr.FailedRecords {
		if i >= 20 {
			sb.WriteString(
				fmt.Sprintf(
					"\n  ... and %d more (press the key for details)\n",
					len(aggr.FailedRecords)-20,
				),
			)
			break
		}

		sb.WriteString(fmt.Sprintf("  %-20s %-30s %-8d %s\n",
			TruncateString(record.SourceIP, 20),
			TruncateString(record.Domain, 30),
			record.Count,
			record.Reason))
	}

	return sb.String()
}
