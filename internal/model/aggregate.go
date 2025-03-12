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

package model

import (
	"strings"
)

// AggregatedReport represents an aggregated view of multiple DMARC reports
type AggregatedReport struct {
	TotalReports  int
	TotalRecords  int
	DateRange     DateRange
	Domains       map[string]int
	Sources       map[string]int
	Dispositions  map[string]int
	DKIMResults   map[string]int
	SPFResults    map[string]int
	FailedRecords []FailedRecord
}

// FailedRecord represents a record that failed DKIM or SPF validation
type FailedRecord struct {
	SourceIP string
	Domain   string
	Count    int
	Reason   string
}

// AggregateReports combines multiple DMARC reports into a single aggregated view
func AggregateReports(reports []DMARCReport) AggregatedReport {
	aggr := AggregatedReport{
		TotalReports: len(reports),
		TotalRecords: 0,
		Domains:      make(map[string]int),
		Sources:      make(map[string]int),
		Dispositions: make(map[string]int),
		DKIMResults:  make(map[string]int),
		SPFResults:   make(map[string]int),
	}

	// Initialize with min/max values for date range
	if len(reports) > 0 {
		aggr.DateRange.Begin = reports[0].ReportMetadata.DateRange.Begin
		aggr.DateRange.End = reports[0].ReportMetadata.DateRange.End
	}

	for _, report := range reports {
		// Update date range
		if report.ReportMetadata.DateRange.Begin.Before(aggr.DateRange.Begin) {
			aggr.DateRange.Begin = report.ReportMetadata.DateRange.Begin
		}
		if report.ReportMetadata.DateRange.End.After(aggr.DateRange.End) {
			aggr.DateRange.End = report.ReportMetadata.DateRange.End
		}

		// Add domain
		aggr.Domains[report.PolicyPublished.Domain]++

		// Process records
		for _, record := range report.Records {
			aggr.TotalRecords++
			aggr.Sources[record.Row.SourceIP]++
			aggr.Dispositions[record.Row.PolicyEvaluated.Disposition]++
			aggr.DKIMResults[record.Row.PolicyEvaluated.DKIM]++
			aggr.SPFResults[record.Row.PolicyEvaluated.SPF]++

			// Track failed authentications
			if record.Row.PolicyEvaluated.DKIM != "pass" ||
				record.Row.PolicyEvaluated.SPF != "pass" {
				reason := ""
				if record.Row.PolicyEvaluated.DKIM != "pass" {
					reason += "DKIM:" + record.Row.PolicyEvaluated.DKIM + " "
				}
				if record.Row.PolicyEvaluated.SPF != "pass" {
					reason += "SPF:" + record.Row.PolicyEvaluated.SPF
				}

				aggr.FailedRecords = append(aggr.FailedRecords, FailedRecord{
					SourceIP: record.Row.SourceIP,
					Domain:   record.Identifiers.HeaderFrom,
					Count:    record.Row.Count,
					Reason:   strings.TrimSpace(reason),
				})
			}
		}
	}

	return aggr
}
